package routing

import "github.com/xolinar/kademlia-go/node"

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

type IKBucketManager interface {
}
