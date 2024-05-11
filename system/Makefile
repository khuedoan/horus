.POSIX:
.PHONY: *

default: apply

apply:
	timoni bundle apply --file gitops.cue

diff:
	timoni bundle apply --file gitops.cue --dry-run --diff
