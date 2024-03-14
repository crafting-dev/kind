package crafting

import (
	"os"
)

func InSandboxWorkspace() bool {
	_, err := os.Stat("/run/sandbox/svc/workspace.sock")
	return err == nil
}

func IsAvailable() bool {
	return InSandboxWorkspace()
}
