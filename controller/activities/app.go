package activities

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"go.temporal.io/sdk/activity"
	"gopkg.in/yaml.v3"
)

func PushRenderedApp(ctx context.Context, appsPath, namespace, app, cluster, registry string) (*PushResult, error) {
	logger := activity.GetLogger(ctx)

	tmpDir, err := os.MkdirTemp("", fmt.Sprintf("%s-%s-", app, cluster))
	if err != nil {
		logger.Error("failed to create temp dir", "error", err)
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.CommandContext(
		ctx,
		"helm", "template",
		"--namespace", namespace,
		app,
		"oci://ghcr.io/bjw-s-labs/helm/app-template:4.1.1",
		"--values", path.Join(namespace, app, cluster+".yaml"),
	)
	cmd.Dir = appsPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	logger.Info("running helm template", "cmd", cmd.String())

	if err := cmd.Run(); err != nil {
		logger.Error("helm template failed", "error", err, "stderr", stderr.String())
		return nil, err
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "rendered.yaml"), stdout.Bytes(), 0644); err != nil {
		logger.Error("failed to write rendered output to file", "error", err)
		return nil, err
	}

	outputPath, err := filepath.Abs(tmpDir)
	if err != nil {
		logger.Error("failed to get absolute path to rendered manifests", "error", err)
		return nil, err
	}

	imageRef := fmt.Sprintf("%s/%s/%s:%s", registry, namespace, app, cluster)
	result, err := PushManifests(ctx, outputPath, imageRef)
	if err != nil {
		logger.Error("failed to push manifests", "error", err)
		return nil, err
	}

	return result, nil
}

func DiscoverApps(ctx context.Context, appsDir string, cluster string) ([]string, error) {
	// TODO logs
	_ = activity.GetLogger(ctx)
	var matched []string
	err := filepath.Walk(appsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), cluster+".yaml") {
			matched = append(matched, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matched, nil
}

type Image struct {
	Repository string
	Tag        string
}

func updateImageTags(node *yaml.Node, newImages []Image) error {
	var walk func(n *yaml.Node)
	walk = func(n *yaml.Node) {
		if n.Kind != yaml.MappingNode {
			for _, child := range n.Content {
				walk(child)
			}
			return
		}
		for i := 0; i < len(n.Content)-1; i += 2 {
			key := n.Content[i]
			val := n.Content[i+1]
			if key.Value == "image" && val.Kind == yaml.MappingNode {
				var repoNode, tagNode *yaml.Node
				for j := 0; j < len(val.Content)-1; j += 2 {
					k := val.Content[j]
					v := val.Content[j+1]
					switch k.Value {
					case "repository":
						repoNode = v
					case "tag":
						tagNode = v
					}
				}
				if repoNode != nil && tagNode != nil {
					for _, img := range newImages {
						if repoNode.Value == img.Repository {
							tagNode.Value = img.Tag
						}
					}
				}
			} else {
				walk(val)
			}
		}
	}
	walk(node)
	return nil
}

func UpdateAppVersion(ctx context.Context, appsDir, namespace, app, cluster string, newImages []Image) error {
	path := filepath.Join(appsDir, namespace, app, fmt.Sprintf("%s.yaml", cluster))

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	if err := updateImageTags(&node, newImages); err != nil {
		return fmt.Errorf("failed to update image tags: %w", err)
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(&node); err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}
	encoder.Close()

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}
