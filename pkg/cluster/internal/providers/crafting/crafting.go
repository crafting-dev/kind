package crafting

import (
	"os"
)

func InSandboxWorkspace() bool {
	// TODO. Better way to check.
	_, ok := os.LookupEnv("SANDBOX_NAME")
	return ok
}

func IsAvailable() bool {
	return InSandboxWorkspace()
}
