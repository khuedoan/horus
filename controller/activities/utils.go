package activities

import (
	"context"

	"go.temporal.io/sdk/activity"
)

// safeHeartbeat sends a heartbeat only if we're in an activity context
func safeHeartbeat(ctx context.Context, details string) {
	defer func() {
		if r := recover(); r != nil {
			// Ignore panic - we're not in an activity context
		}
	}()
	activity.RecordHeartbeat(ctx, details)
}
