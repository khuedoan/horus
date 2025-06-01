package activities

import (
	"strings"
	"testing"
)

func TestPruneGraph(t *testing.T) {
	input := `
digraph {
	"azure-policies" ;
	"azure-vm-disk-backups/backup-policies/daily-30-day-retention" ;
	"azure-vm-disk-backups/backup-policies/daily-30-day-retention" -> "cloud";
	"azure-vm-disk-backups/backup-policies/daily-30-day-retention" -> "azure-vm-disk-backups/backup-vault";
	"azure-vm-disk-backups/backup-vault" ;
	"azure-vm-disk-backups/backup-vault" -> "cloud";
	"cloud" ;
	"ecomnet-vng" ;
	"generated-secrets" ;
	"generated-secrets" -> "topology";
	"il4/generated-secrets" ;
	"il4/generated-secrets" -> "topology";
	"legacy-bridge" ;
	"legacy-bridge" -> "topology";
	"legacy-bridge" -> "network";
	"local-distribution" ;
	"network" ;
	"secrets" ;
	"secrets" -> "topology";
	"topology" ;
}`

	expectedContains := []string{
		`"azure-vm-disk-backups/backup-policies/daily-30-day-retention"`,
		`"azure-vm-disk-backups/backup-policies/daily-30-day-retention" -> "azure-vm-disk-backups/backup-vault"`,
		`"azure-vm-disk-backups/backup-vault"`,
		`"topology"`,
		`"generated-secrets" -> "topology"`,
		`"il4/generated-secrets" -> "topology"`,
		`"legacy-bridge" -> "topology"`,
		`"secrets" -> "topology"`,
		`"local-distribution"`,
	}

	changed := []string{
		`"azure-vm-disk-backups/backup-vault"`,
		`"topology"`,
		`"local-distribution"`,
	}

	pruned, err := pruneGraph(input, changed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, expected := range expectedContains {
		if !strings.Contains(pruned, expected) {
			t.Errorf("expected pruned output to contain: %s", expected)
		}
	}
}
