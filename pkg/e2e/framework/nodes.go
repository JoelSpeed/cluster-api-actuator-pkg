package framework

import (
	"context"
	"fmt"

	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AddNodeCondition adds a condition in the given Node's status.
func AddNodeCondition(c client.Client, node *corev1.Node, cond corev1.NodeCondition) error {
	nodeCopy := node.DeepCopy()
	nodeCopy.Status.Conditions = append(nodeCopy.Status.Conditions, cond)

	return c.Status().Patch(context.Background(), nodeCopy, client.MergeFrom(node))
}

// GetNodeForMachine retrieves the node backing the given Machine.
func GetNodeForMachine(c client.Client, m *mapiv1beta1.Machine) (*corev1.Node, error) {
	if m.Status.NodeRef == nil {
		return nil, fmt.Errorf("%s: machine has no NodeRef", m.Name)
	}

	node := &corev1.Node{}
	nodeName := client.ObjectKey{Name: m.Status.NodeRef.Name}

	if err := c.Get(context.Background(), nodeName, node); err != nil {
		return nil, err
	}

	return node, nil
}

// GetWorkerNodes returns all nodes with the nodeWorkerRoleLabel label
func GetWorkerNodes(c client.Client) ([]corev1.Node, error) {
	workerNodes := &corev1.NodeList{}
	err := c.List(context.TODO(), workerNodes,
		client.InNamespace(MachineAPINamespace),
		client.MatchingLabels(map[string]string{WorkerNodeRoleLabel: ""}),
	)

	if err != nil {
		return nil, err
	}

	return workerNodes.Items, nil
}

// FilterReadyNodes fileter the list of nodes and returns the list with ready nodes
func FilterReadyNodes(nodes []corev1.Node) []corev1.Node {
	var readyNodes []corev1.Node
	for _, n := range nodes {
		if IsNodeReady(&n) {
			readyNodes = append(readyNodes, n)
		}
	}
	return readyNodes
}
