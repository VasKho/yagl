package yagl

// Type for Node representation
type Node struct {
	parent     Arcs
	child      Arcs
	el_type    int
	identifier string
}

type Arcs map[string]*Node

// Create node
func NewNode(identifier string, el_type int) (*Node, string) {
	hash := genHash(identifier)
	parent, child := make(Arcs), make(Arcs)
	new_node := &Node{parent, child, el_type, identifier}
	return new_node, hash
}

// Compare two nodes
func (node_1 Node) IsEqual(node_2 Node) bool {
	return node_1.identifier == node_2.identifier
}

// Check node on emptyness
func (node_1 Node) IsEmpty() bool {
	return node_1.IsEqual(Node{})
}
