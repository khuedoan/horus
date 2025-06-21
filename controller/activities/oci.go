package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"

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
	cmd := exec.CommandContext(ctx, "oras", "push", "--format=json", "--plain-http", image, ".")
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
