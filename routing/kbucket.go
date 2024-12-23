package routing

import (
	"sync"

	"github.com/xolinar/kademlia-go/node"
)

// KSize defines the maximum number of nodes a KBucket can hold.
type KSize uint8

// IKBucket defines the interface for managing a K-bucket in the Kademlia DHT routing table.
//
// A K-bucket is a structure used in the Kademlia protocol to store a subset of nodes based on
// their XOR distance from the local node. Each K-bucket contains nodes within a specific distance range,
// enabling efficient node management by keeping track of nodes that are close in the keyspace.
// Each K-bucket has a fixed capacity (KSize) and follows Kademlia's rules for adding, removing, and evicting nodes.
type IKBucket interface {
	// Nodes returns a slice of nodes currently stored in the K-bucket, ordered by proximity
	// to the bucket's range in the keyspace. The most recently active nodes are placed at the end.
	// This allows the K-bucket to prioritize recently active nodes for future interactions.
	Nodes() []node.INode

	// KSize returns the maximum number of nodes that the K-bucket can hold.
	// This value is typically fixed across all K-buckets, defining their capacity limit.
	KSize() KSize

	// Add inserts a new node into the K-bucket. If the K-bucket is full,
	// the least recently seen node may be evicted to make room for the new node.
	// This follows the Kademlia protocol's rule of prioritizing recently active nodes.
	Add(newNode node.INode)

	// Remove deletes a node from the K-bucket using its NodeID.
	// This is typically used to remove nodes that are unreachable or outdated.
	Remove(id node.NodeID)

	// Update refreshes the position of a node in the K-bucket by its NodeID,
	// indicating that it was recently active. This moves the node to the end
	// of the list, reflecting recent interaction.
	Update(id node.NodeID)

	// Size returns the current number of nodes in the K-bucket.
	// This helps monitor how many nodes are actively being tracked and managed within the bucket.
	Size() KSize

	// IsFull checks whether the K-bucket has reached its maximum capacity.
	IsFull() bool
}

// KBucket represents a container for nodes within the Kademlia network.
//
// A KBucket is a segment of Kademlia's routing table that maintains a list of nodes
// within a specific distance range from the current node. Each KBucket has a limited
// capacity and follows a Least Recently Seen (LRS) eviction policy, which ensures
// that the oldest nodes are removed to make space for newly added or recently active nodes.
//
// References:
//   - [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"] [Section 2.2, "Node State"]
//
// [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"]: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
type KBucket struct {
	// nodes is a slice of nodes stored in this KBucket. These nodes represent peers at a specific
	// distance range from the current node. The slice maintains nodes in order of
	// activity, with the most recently seen node positioned at the end.
	nodes []node.INode

	// ksize is the maximum number of nodes that the KBucket can contain. If this limit is reached
	// when adding a new node, the oldest node is evicted to make room for the new node.
	ksize KSize

	// mu is a mutex used to synchronize access to the nodes slice, ensuring that all operations
	// on the KBucket are thread-safe in concurrent environments.
	mu sync.Mutex
}

// NewKBucket creates and returns a new KBucket instance with a specified capacity for storing nodes.
func NewKBucket(ksize KSize) *KBucket {
	return &KBucket{
		nodes: make([]node.INode, 0, ksize),
		ksize: ksize,
		mu:    sync.Mutex{},
	}
}

// Nodes returns a slice of nodes stored in the KBucket.
//
// This method provides access to the nodes contained within the KBucket, representing peers
// at a specific distance from the local node. The nodes are ordered by their last-seen time.
func (kb *KBucket) Nodes() []node.INode {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	return kb.nodes
}

// KSize returns the maximum number of nodes that the KBucket can hold.
//
// This method provides the maximum capacity of the KBucket, which is fixed and determined
// during initialization. The capacity is used to manage node evictions when the KBucket is full.
func (kb *KBucket) KSize() KSize {
	return kb.ksize
}

// Add inserts a new node into the KBucket.
//
// If the node already exists, it is removed from its current position and re-added to the end
// of the list to reflect its recent activity. If the KBucket is full and does not contain the new node,
// the oldest node (at the beginning) is removed to make space.
func (kb *KBucket) Add(newNode node.INode) {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	// Check if the node already exists and update its position if necessary.
	for i, n := range kb.nodes {
		if n.ID().Equals(newNode.ID()) {
			// Move the node to the end of the list (most recently seen).
			kb.nodes = append(kb.nodes[:i], kb.nodes[i+1:]...)
			kb.nodes = append(kb.nodes, newNode)
			return
		}
	}

	// If the bucket is not full, just add the new node.
	if len(kb.nodes) < int(kb.ksize) {
		kb.nodes = append(kb.nodes, newNode)
		return
	}

	// If the bucket is full, evict the oldest node and add the new one.
	kb.nodes = append(kb.nodes[1:], newNode)
}

// Remove deletes a node from the KBucket using its NodeID.
//
// This is typically used to remove nodes that are unreachable or outdated.
func (kb *KBucket) Remove(id node.NodeID) {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	for i, n := range kb.nodes {
		if n.ID().Equals(id) {
			// Remove the node from the slice.
			kb.nodes = append(kb.nodes[:i], kb.nodes[i+1:]...)
			return
		}
	}
}

// Update refreshes the position of a node in the KBucket by its NodeID,
// indicating that it was recently active. This moves the node to the end
// of the list, reflecting recent interaction.
func (kb *KBucket) Update(id node.NodeID) {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	for i, n := range kb.nodes {
		if n.ID().Equals(id) {
			// Move the node to the end of the list (most recently seen).
			nodeToMove := kb.nodes[i]
			kb.nodes = append(kb.nodes[:i], kb.nodes[i+1:]...)
			kb.nodes = append(kb.nodes, nodeToMove)
			return
		}
	}
}

// Size returns the current number of nodes in the KBucket.
//
// This method calculates and returns the number of nodes currently stored in the KBucket.
// It is useful for monitoring the number of active nodes being tracked.
func (kb *KBucket) Size() KSize {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	return KSize(len(kb.nodes))
}

// IsFull checks whether the KBucket has reached its maximum capacity.
//
// This method returns true if the current number of nodes is equal to or greater than the KBucket's capacity.
// If there is still space available, it returns false.
func (kb *KBucket) IsFull() bool {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	return len(kb.nodes) >= int(kb.ksize)
}
