package cli

import "testing"

func TestSplitArgs(t *testing.T) {
	SplitArgs([]string{"wmexec", "--user", "1200", "--", "wmspawner"})
}
