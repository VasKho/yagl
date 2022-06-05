package yagl

import (
	"fmt"
	"sync"
)

// Type for context representation
type Context struct {
	nodes       Arcs
	unnamed_num int
	size        int
	*sync.RWMutex
}

// Function to create new context
func NewContext() *Context {
	context := Context{make(Arcs), 0, 0, &sync.RWMutex{}}
	return &context
}

// Get all nodes from context
func (context Context) GetNodes() Arcs {
	return context.nodes
}

// Get number of nodes in context
func (context Context) GetSize() int {
	return context.size
}

// Get node by identifier
func (context Context) GetNode(identifier string, el_type int) *Node {
	defer context.RUnlock()
	hash := genHash(identifier, el_type)
	context.RLock()
	if node, ok := context.nodes[hash]; ok {
		return node
	}
	return &Node{}
}

// Get node address by identifier
func (context Context) GetNodeAddr(identifier string, el_type int) string {
	defer context.RUnlock()
	hash := genHash(identifier, el_type)
	context.RLock()
	if _, ok := context.nodes[hash]; ok {
		return hash
	}
	return ""
}

// Get arc by its begin and end nodes
func (context Context) GetArc(begin, end *Node, el_type int) *Node {
	defer begin.RUnlock()
	defer end.RUnlock()
	defer context.RUnlock()
	begin.RLock()
	end.RLock()
	if begin.IsEmpty() {
		return &Node{}
	}
	if end.IsEmpty() {
		return &Node{}
	}
	context.RLock()
	for index_1, arc_1 := range begin.child {
		for index_2 := range end.parent {
			if index_1 == index_2 && arc_1.el_type == el_type {
				return arc_1
			}
		}
	}
	return &Node{}
}

// Get arc address by its begin and end nodes
func (context Context) GetArcAddr(begin, end *Node, el_type int) string {
	defer begin.RUnlock()
	defer end.RUnlock()
	defer context.RUnlock()
	begin.RLock()
	end.RLock()
	if begin.IsEmpty() {
		return ""
	}
	if end.IsEmpty() {
		return ""
	}
	context.RLock()
	for index_1, arc_1 := range begin.child {
		for index_2 := range begin.parent {
			if index_1 == index_2 && arc_1.el_type == el_type {
				return index_1
			}
		}
	}
	return ""
}

// Add new node to context
func (context *Context) AddNode(identifier string, el_type int) error {
	if !context.GetNode(identifier, el_type).IsEmpty() {
		return fmt.Errorf("Error: Key %s already exists", identifier)
	}
	if identifier == "" {
		identifier = "unnamed_" + fmt.Sprintf("%d", context.unnamed_num)
		context.unnamed_num++
	}
	new_node, hash := NewNode(identifier, el_type)
	context.Lock()
	context.nodes[hash] = new_node
	context.size++
	context.Unlock()
	return nil
}

// Create new arc of given type between two nodes
func (context *Context) AddArc(begin, end *Node, el_type int) error {
	if el_type&Node_t != 0 {
		return fmt.Errorf("Error: Unable to create arc of Node_t type")
	}
	begin.RLock()
	if begin.IsEmpty() {
		defer begin.RUnlock()
		return fmt.Errorf("Error: Unknown key %s", begin.identifier)
	}
	end.RLock()
	if end.IsEmpty() {
		defer end.RUnlock()
		return fmt.Errorf("Error: Unknown key %s", end.identifier)
	}
	begin.RUnlock()
	end.RUnlock()
	if !context.GetArc(begin, end, el_type).IsEmpty() {
		return nil
	}
	begin.RLock()
	end.RLock()
	context.Lock()
	context.unnamed_num++
	new_node, hash := NewNode("unnamed_"+fmt.Sprintf("%d", context.unnamed_num), el_type)
	context.Unlock()
	context.AddNode("unnamed_" + fmt.Sprintf("%d", context.unnamed_num), el_type)
	begin.child[hash] = new_node
	end.parent[hash] = new_node
	begin.RUnlock()
	end.RUnlock()
	return nil
}

// Remove node by its identifier
func (context *Context) RemoveNode(identifier string, el_type int) error {
	node_1, addr := context.GetNode(identifier, el_type), context.GetNodeAddr(identifier, el_type)
	if node_1.IsEmpty() {
		return fmt.Errorf("Error: Unknown key %s", identifier)
	}
	node_1.Lock()
	for _, node := range node_1.parent {
		node.Lock()
		for _, nodes_parent := range node.parent {
			nodes_parent.Lock()
			if nodes_parent.el_type&Node_t != 0 {
				context.RemoveArc(nodes_parent, node_1, node.el_type)
			}
			nodes_parent.Unlock()
		}
		node.Unlock()
	}
	for _, node := range node_1.child {
		node.Lock()
		for _, nodes_child := range node.child {
			nodes_child.Lock()
			if nodes_child.el_type&Node_t != 0 {
				context.RemoveArc(node_1, nodes_child, node.el_type)
			}
			nodes_child.Unlock()
		}
		node.Unlock()
	}
	node_1.Unlock()
	context.Lock()
	delete(context.nodes, addr)
	context.Unlock()
	return nil
}

// Remove arc by its begin, end and type
func (context *Context) RemoveArc(begin, end *Node, el_type int) error {
	begin.RLock()
	if begin.IsEmpty() {
		defer begin.RUnlock()
		return fmt.Errorf("Error: Unknown begin node")
	}
	end.RLock()
	if end.IsEmpty() {
		defer end.RUnlock()
		return fmt.Errorf("Error: Unknown end node")
	}
	arc, arc_addr := context.GetArc(begin, end, el_type), context.GetArcAddr(begin, end, el_type)
	if arc_addr != "" {
		arc.Lock()
		for hash, node := range arc.parent {
			node.Lock()
			if node.el_type&Node_t != 0 {
				delete(node.child, arc_addr)
				delete(arc.parent, hash)
			} else {
				for _, parent := range node.parent {
					context.RemoveArc(parent, arc, node.el_type)
				}
				for _, child := range node.child {
					context.RemoveArc(arc, child, node.el_type)
				}
			}
			node.Unlock()
		}
		for hash, node := range arc.child {
			node.Lock()
			if node.el_type&Node_t != 0 {
				delete(node.parent, arc_addr)
				delete(arc.child, hash)
			} else {
				for _, parent := range node.parent {
					context.RemoveArc(parent, arc, node.el_type)
				}
				for _, child := range node.child {
					context.RemoveArc(arc, child, node.el_type)
				}
			}
			node.Unlock()
		}
		arc.Unlock()
		context.Lock()
		delete(context.nodes, arc_addr)
		context.unnamed_num--
		context.Unlock()
		begin.RUnlock()
		end.RUnlock()
	}
	return nil
}
