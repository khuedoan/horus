package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"go.temporal.io/sdk/activity"
)

type PushResult struct {
	Reference       string            `json:"reference"`
	MediaType       string            `json:"mediaType"`
	Digest          string            `json:"digest"`
	Size            int               `json:"size"`
	Annotations     map[string]string `json:"annotations"`
	ArtifactType    string            `json:"artifactType"`
	ReferenceAsTags []string          `json:"referenceAsTags"`
}

func PushManifests(ctx context.Context, path string, image string) (*PushResult, error) {
	logger := activity.GetLogger(ctx)
	cmd := exec.CommandContext(ctx, "nix", "develop", "--command", "oras", "push", "--format=json", "--plain-http", image, ".")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logger.Error("oras push failed", "error", err, "stderr", stderr.String())
		return nil, err
	}

	var result PushResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		logger.Error("failed to parse oras output", "error", err, "output", stdout.String())
		return nil, err
	}

	return &result, nil
}

func PushRenderedHelm(ctx context.Context, appsPath, namespace, app, cluster, registry string) (*PushResult, error) {
	logger := activity.GetLogger(ctx)

	tmpDir, err := os.MkdirTemp("", fmt.Sprintf("%s-%s-", app, cluster))
	if err != nil {
		logger.Error("failed to create temp dir", "error", err)
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.CommandContext(
		ctx,
		"nix", "develop", "--command",
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
