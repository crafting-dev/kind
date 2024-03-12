package crafting

import (
	"bytes"
	"fmt"
	"os"

	"sigs.k8s.io/kind/pkg/cluster/internal/providers"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"
	"sigs.k8s.io/kind/pkg/errors"
)

func InSandboxWorkspace() bool {
	// TODO. Better way to check.
	_, ok := os.LookupEnv("SANDBOX_NAME")
	return ok
}

func IsAvailable() bool {
	return InSandboxWorkspace()
}

func PatchKubeProxy(p providers.Provider, name string) error {
	// find a control plane node to patch kube-proxy
	n, err := p.ListNodes(name)
	if err != nil {
		return err
	}
	nodes, err := nodeutils.ControlPlaneNodes(n)
	if err != nil {
		return err
	}
	if len(nodes) < 1 {
		return errors.Errorf("could not locate any control plane nodes for cluster named '%s'. "+
			"Use the --name option to select a different cluster", name)
	}
	var buff bytes.Buffer
	node := nodes[0]
	if err = node.Command(
		"kubectl", "-n", "kube-system", "patch", "ds", "kube-proxy", "--type", "json",
		"-p", "[{'op': 'replace', 'path': '/spec/template/spec/containers/0/securityContext', 'value':{'capabilities':{'add':['NET_ADMIN']}}}]").
		SetStderr(&buff).Run(); err != nil {
		return fmt.Errorf("patch kube-system capabilities error: %w, %s", err, buff.String())
	}

	return nil
}
