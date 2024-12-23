package node

import (
	"fmt"
	"net"
)

// INode defines the interface for interacting with nodes in a Kademlia DHT network.
//
// The INode interface provides a standard set of methods for interacting with nodes,
// enabling the abstraction of node properties such as identifiers, network addresses,
// and distance calculations. Implementations of this interface can include various
// node types, which may extend the functionality with additional attributes or behaviors,
// such as custom metadata or specialized communication capabilities.
//
// By implementing the INode interface, various node structures can be created to handle
// different use cases or extended behaviors, while maintaining compatibility with the
// Kademlia routing logic. This allows for flexibility and extensibility in how nodes
// are represented and interacted with in the distributed network.
type INode interface {
	// ID returns the unique identifier of the node, used for routing and distance calculations
	// within the Kademlia network.
	ID() NodeID

	// Address returns the IP address of the node, either IPv4 or IPv6, which is used for
	// establishing network communication.
	Address() net.IP

	// Port returns the port number the node listens on, used in conjunction with the IP address
	// to facilitate network connections.
	Port() uint16

	// Distance calculates and returns the XOR-based distance between the current node and another
	// node, which is essential for determining proximity and routing decisions in the Kademlia protocol.
	Distance(node INode) [20]byte
}

// Node represents a node in the Kademlia DHT network.
//
// Each Node is identified by a unique identifier (ID), and it is associated
// with an IP address and port for establishing network connections. The Node
// struct is fundamental in the Kademlia protocol, storing the necessary information
// for routing, communication, and maintaining a decentralized distributed hash table (DHT).
//
// References:
//   - [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"] [Section 2.2, "Node State"]
//
// [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"]: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
type Node struct {
	id      NodeID // Unique identifier for the node.
	address net.IP // IP address for network communication (IPv4 or IPv6).
	port    uint16 // Port for listening to incoming connections (range 0-65535).
}

// NewNode creates and returns a new Node instance with a unique NodeID.
// Added validations for address and port correctness.
func NewNode(id NodeID, address net.IP, port uint16) (*Node, error) {
	if len(address) == 0 || address == nil {
		return nil, fmt.Errorf("invalid IP address: cannot be nil or empty")
	}
	return &Node{
		id:      id,
		address: address,
		port:    port,
	}, nil
}

// ID returns the NodeID of the current node.
//
// The NodeID is a unique identifier generated from the node's relevant data,
// such as its IP address and other parameters, used for sorting and determining
// proximity to other nodes in the Kademlia network.
func (n *Node) ID() NodeID {
	return n.id
}

// Address returns the IP address of the current node.
//
// The IP address is used for network communication and can be either IPv4 or IPv6.
// This address is necessary for establishing connections with other nodes in the network.
func (n *Node) Address() net.IP {
	return n.address
}

// Port returns the port number the current node is listening on.
//
// The port is used for establishing network connections, either over TCP or UDP,
// and must be within the valid range (0-65535) to ensure proper communication.
func (n *Node) Port() uint16 {
	return n.port
}

// Distance calculates the distance between the current node and another node in the Kademlia DHT.
//
// The distance is determined using the XOR metric, which is applied between the NodeIDs of
// the current node and the other node. The result is a 160-bit value that represents the
// proximity or distance between the nodes in the Kademlia keyspace.
//
// The smaller the result, the closer the nodes are in the network.
//
// References:
//   - [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"] [Section 2.1, "XOR Metric"]
//
// [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"]: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
func (n *Node) Distance(other INode) [20]byte {
	return n.id.XOR(other.ID())
}
