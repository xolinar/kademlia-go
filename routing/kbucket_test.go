package routing_test

import (
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xolinar/kademlia-go/node"
	"github.com/xolinar/kademlia-go/routing"
)

// newTestNodeID generates a NodeID for testing purposes.
func newTestNodeID(id string) node.NodeID {
	val, _ := node.NewNodeID([]byte(id))
	return val // Assumes NodeID is a string; adjust if it's another type.
}

// newTestNode creates a new node with the given ID, IP address, and port for testing.
func newTestNode(id string, address string, port uint16) *node.Node {
	ip := net.ParseIP(address)
	n, _ := node.NewNode(newTestNodeID(id), ip, port)
	return n
}

// TestAddAndRetrieveNodes ensures that nodes can be added to the KBucket
// and retrieved correctly. It validates that all added nodes are stored,
// and the KBucket respects its maximum capacity.
func TestAddAndRetrieveNodes(t *testing.T) {
	kbucket := routing.NewKBucket(3)

	node1 := newTestNode("node1", "127.0.0.1", 8001)
	node2 := newTestNode("node2", "127.0.0.1", 8002)
	node3 := newTestNode("node3", "127.0.0.1", 8003)

	// Add nodes to the KBucket
	kbucket.Add(node1)
	kbucket.Add(node2)
	kbucket.Add(node3)

	// Verify that all nodes were added
	nodes := kbucket.Nodes()
	assert.Len(t, nodes, 3)
	assert.Contains(t, nodes, node1)
	assert.Contains(t, nodes, node2)
	assert.Contains(t, nodes, node3)
}

// TestEvictionPolicy verifies that the KBucket evicts the least recently seen
// node when it reaches its maximum capacity, following the LRS eviction policy.
func TestEvictionPolicy(t *testing.T) {
	kbucket := routing.NewKBucket(2)

	node1 := newTestNode("node1", "127.0.0.1", 8001)
	node2 := newTestNode("node2", "127.0.0.1", 8002)
	node3 := newTestNode("node3", "127.0.0.1", 8003)

	// Add 3 nodes to a KBucket with a size of 2
	kbucket.Add(node1)
	kbucket.Add(node2)
	kbucket.Add(node3)

	// Verify that the oldest node (node1) is evicted
	nodes := kbucket.Nodes()
	assert.Len(t, nodes, 2)
	assert.NotContains(t, nodes, node1)
	assert.Contains(t, nodes, node2)
	assert.Contains(t, nodes, node3)
}

// TestRemoveNode checks that a specific node can be removed
// from the KBucket based on its NodeID.
func TestRemoveNode(t *testing.T) {
	kbucket := routing.NewKBucket(3)

	node1 := newTestNode("node1", "127.0.0.1", 8001)
	node2 := newTestNode("node2", "127.0.0.1", 8002)

	// Add nodes
	kbucket.Add(node1)
	kbucket.Add(node2)

	// Remove node1
	kbucket.Remove(node1.ID())

	// Verify that node1 was removed
	nodes := kbucket.Nodes()
	assert.Len(t, nodes, 1)
	assert.Contains(t, nodes, node2)
	assert.NotContains(t, nodes, node1)
}

// TestUpdateNode validates that a node's position is updated in the KBucket
// when it is accessed, ensuring that it reflects the most recently seen order.
func TestUpdateNode(t *testing.T) {
	kbucket := routing.NewKBucket(3)

	node1 := newTestNode("node1", "127.0.0.1", 8001)
	node2 := newTestNode("node2", "127.0.0.1", 8002)

	// Add nodes
	kbucket.Add(node1)
	kbucket.Add(node2)

	// Update node1 to reflect recent activity
	kbucket.Update(node1.ID())

	nodes := kbucket.Nodes()
	// Verify that node1 is moved to the end of the list
	assert.Equal(t, nodes[len(nodes)-1], node1)
}

// TestRaceCondition tests for race conditions when multiple goroutines
// concurrently modify or access a KBucket instance under heavy load.
func TestRaceCondition(t *testing.T) {
	kbucket := routing.NewKBucket(5)

	// Simulate heavy load with multiple goroutines
	var wg sync.WaitGroup
	numGoroutines := 1000
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			node := newTestNode(fmt.Sprintf("node%d", i), "127.0.0.1", uint16(8000+i))
			kbucket.Add(node)
		}(i)
	}

	wg.Wait()

	// Verify that the KBucket still holds a reasonable number of nodes
	// Check if the number of nodes doesn't exceed the KBucket capacity
	nodes := kbucket.Nodes()
	assert.Len(t, nodes, 5) // KBucket should hold only 5 nodes, with eviction applied
}
