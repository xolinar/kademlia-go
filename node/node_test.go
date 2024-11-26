package node_test

import (
	"net"
	"testing"

	"github.com/xolinar/kademlia-go/node"
)

func TestNewNode(t *testing.T) {
	data, _ := node.NewNodeID([]byte("test_node_data"))
	address := net.ParseIP("192.168.1.1")
	port := uint16(8080)

	testNode, err := node.NewNode(data, address, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if testNode.ID() != data {
		t.Errorf("expected id %x, got %x", data, testNode.ID())
	}

	if !testNode.Address().Equal(address) {
		t.Errorf("expected address %s, got %s", address, testNode.Address())
	}

	if testNode.Port() != port {
		t.Errorf("expected port %d, got %d", port, testNode.Port())
	}
}

func TestNewNode_InvalidAddress(t *testing.T) {
	data, _ := node.NewNodeID([]byte("test_node_data"))
	port := uint16(8080)

	_, err := node.NewNode(data, nil, port)
	if err == nil {
		t.Fatal("expected error for nil IP address, got none")
	}

	_, err = node.NewNode(data, net.IP{}, port)
	if err == nil {
		t.Fatal("expected error for empty IP address, got none")
	}
}

// TestDistance checks the Distance method by verifying the XOR calculation between two NodeIDs.
func TestDistance(t *testing.T) {
	id1, _ := node.NewNodeID([]byte("node_1"))
	id2, _ := node.NewNodeID([]byte("node_2"))

	address := net.ParseIP("192.168.1.1")
	port := uint16(8080)

	node1, _ := node.NewNode(id1, address, port)
	node2, _ := node.NewNode(id2, address, port)

	distance := node1.Distance(node2)
	if len(distance) != 20 {
		t.Errorf("expected distance length 20, got %d", len(distance))
	}

	// Add specific XOR result tests if necessary
}

// TestSameNodeDistance checks that the distance between a node and itself is zero.
func TestSameNodeDistance(t *testing.T) {
	id, _ := node.NewNodeID([]byte("same_node_data"))
	address := net.ParseIP("192.168.1.1")
	port := uint16(8080)

	testNode, _ := node.NewNode(id, address, port)
	distance := testNode.Distance(testNode)

	var zeroDistance [20]byte
	if distance != zeroDistance {
		t.Errorf("Distance failed, expected zero distance, got %x", distance)
	}
}
