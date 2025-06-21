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
		"helm", "template", "--namespace", namespace, app, "oci://ghcr.io/bjw-s-labs/helm/app-template:4.1.1", "--values", path.Join(namespace, app, cluster+".yaml"), "--output-dir", tmpDir,
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
