package crafting

import (
	"bytes"
	"fmt"
	"os"

	"sigs.k8s.io/kind/pkg/cluster/internal/create/actions"
	"sigs.k8s.io/kind/pkg/cluster/internal/providers"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/internal/apis/config"
)

type kubeProxyPatchAction struct {
	p   providers.Provider
	cfg *config.Cluster
}

func NewKubeProxyPatchAction(p providers.Provider, cfg *config.Cluster) actions.Action {
	return &kubeProxyPatchAction{p: p, cfg: cfg}
}

func (a *kubeProxyPatchAction) Execute(ctx *actions.ActionContext) error {
	name := a.cfg.Name
	// find a control plane node to patch kube-proxy
	n, err := a.p.ListNodes(name)
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

func InSandboxWorkspace() bool {
	// TODO. Better way to check.
	_, err := os.Stat("/run/sandbox/svc/workspace.sock")
	return err == nil
}

func IsAvailable() bool {
	return InSandboxWorkspace()
}
