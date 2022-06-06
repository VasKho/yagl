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
func (context *Context) GetNodes(exchange *Exchange[Arcs]) {
	if exchange != nil {
		if exchange.Wg != nil { defer exchange.Wg.Done() }
	}
	defer context.RUnlock()
	context.RLock()
	exchange.Result <-context.nodes
}

// Get number of nodes in context
func (context *Context) GetSize() int {
	defer context.RUnlock()
	context.RLock()
	return context.size
}

// Get node by identifier
func (context *Context) GetNode(identifier string, el_type int, exchange *Exchange[*Node]) {
	if exchange != nil {
		if exchange.Wg != nil { defer exchange.Wg.Done() }
	}
	defer context.RUnlock()
	hash := genHash(identifier, el_type)
	context.RLock()
	if node, ok := context.nodes[hash]; ok {
		exchange.Result <-node
		return
	}
	exchange.Result <-&Node{}	
}

// Get node address by identifier
func (context *Context) GetNodeAddr(identifier string, el_type int, exchange *Exchange[string]) {
	if exchange != nil {
		if exchange.Wg != nil { defer exchange.Wg.Done() }
	}
	defer context.RUnlock()
	hash := genHash(identifier, el_type)
	context.RLock()
	if _, ok := context.nodes[hash]; ok {
		exchange.Result <-hash
		return
	}
	exchange.Result <-""
}

// Get arc by its begin and end nodes
func (context *Context) GetArc(begin, end *Node, el_type int, exchange *Exchange[*Node]) {
	if exchange != nil {
		if exchange.Wg != nil { defer exchange.Wg.Done() }
	}
	defer begin.RUnlock()
	defer end.RUnlock()
	defer context.RUnlock()
	begin.RLock()
	end.RLock()
	if begin.IsEmpty() {
		exchange.Result <-&Node{}
		return
	}
	if end.IsEmpty() {
		exchange.Result <-&Node{}
		return
	}
	context.RLock()
	for index_1, arc_1 := range begin.child {
		for index_2 := range end.parent {
			if index_1 == index_2 && arc_1.el_type == el_type {
				exchange.Result <-arc_1
				return
			}
		}
	}
	exchange.Result <-&Node{}
}

// Get arc address by its begin and end nodes
func (context *Context) GetArcAddr(begin, end *Node, el_type int, exchange *Exchange[string]) {
	if exchange != nil {
		if exchange.Wg != nil { defer exchange.Wg.Done() }
	}
	defer begin.RUnlock()
	defer end.RUnlock()
	defer context.RUnlock()
	begin.RLock()
	end.RLock()
	if begin.IsEmpty() {
		exchange.Result <-""
		return
	}
	if end.IsEmpty() {
		exchange.Result <-""
		return
	}
	context.RLock()
	for index_1, arc_1 := range begin.child {
		for index_2 := range begin.parent {
			if index_1 == index_2 && arc_1.el_type == el_type {
				exchange.Result <-index_1
				return
			}
		}
	}
	exchange.Result <-""
}

// Add new node to context
func (context *Context) AddNode(identifier string, el_type int, exchange *Exchange[error]) {
	if exchange != nil {
		if exchange.Wg != nil { defer exchange.Wg.Done() }
	}
	exch := Exchange[*Node]{make(chan *Node), nil}
	go context.GetNode(identifier, el_type, &exch)
	if !(<-exch.Result).IsEmpty() {
		exchange.Result <-fmt.Errorf("AddNode: Key %s already exists", identifier)
		return
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
}

// Create new arc of given type between two nodes
func (context *Context) AddArc(begin, end *Node, el_type int, exchange *Exchange[error]) {
	if exchange != nil {
		if exchange.Wg != nil { defer exchange.Wg.Done() }
	}
	if el_type&Node_t != 0 {
		exchange.Result <-fmt.Errorf("AddArc: Unable to create arc of Node_t type")
		return
	}
	begin.RLock()
	if begin.IsEmpty() {
		defer begin.RUnlock()
		exchange.Result <-fmt.Errorf("AddArc: Begin node is out of context")
		return
	}
	end.RLock()
	if end.IsEmpty() {
		defer end.RUnlock()
		exchange.Result <-fmt.Errorf("AddArc: End node is out of context")
		return
	}
	begin.RUnlock()
	end.RUnlock()
	get_res := Exchange[*Node]{make(chan *Node), nil}
	go context.GetArc(begin, end, el_type, &get_res)
	if !(<-get_res.Result).IsEmpty() {
		return
	}
	begin.RLock()
	end.RLock()
	context.Lock()
	context.unnamed_num++
	new_node, hash := NewNode("unnamed_" + fmt.Sprintf("%d", context.unnamed_num), el_type)
	context.Unlock()
	add_node_res := Exchange[error]{make(chan error, 2), &sync.WaitGroup{}}
	add_node_res.Wg.Add(1)
	go context.AddNode("unnamed_" + fmt.Sprintf("%d", context.unnamed_num), el_type, &add_node_res)
	add_node_res.Wg.Wait()
	begin.child[hash] = new_node
	end.parent[hash] = new_node
	begin.RUnlock()
	end.RUnlock()
}

// Remove node by its identifier
func (context *Context) RemoveNode(identifier string, el_type int, exchange *Exchange[error]) {
	if exchange != nil {
		if exchange.Wg != nil { defer exchange.Wg.Done() }
	}
	node_res := Exchange[*Node]{make(chan *Node), nil}
	addr_res := Exchange[string]{make(chan string), nil}
	go context.GetNode(identifier, el_type, &node_res)
	go context.GetNodeAddr(identifier, el_type, &addr_res)
	node_1, addr := <-node_res.Result, <-addr_res.Result
	if node_1.IsEmpty() {
		exchange.Result <-fmt.Errorf("RemoveNode: Unknown key %s", identifier)
		return
	}
	node_1.Lock()
	for _, node := range node_1.parent {
		node.Lock()
		for _, nodes_parent := range node.parent {
			nodes_parent.Lock()
			if nodes_parent.el_type&Node_t != 0 {
				go context.RemoveArc(nodes_parent, node_1, node.el_type, nil)
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
				go context.RemoveArc(node_1, nodes_child, node.el_type, nil)
			}
			nodes_child.Unlock()
		}
		node.Unlock()
	}
	node_1.Unlock()
	context.Lock()
	delete(context.nodes, addr)
	context.Unlock()
}

// Remove arc by its begin, end and type
func (context *Context) RemoveArc(begin, end *Node, el_type int, exchange *Exchange[error]) {
	if exchange != nil {
		if exchange.Wg != nil { defer exchange.Wg.Done() }
	}
	begin.RLock()
	if begin.IsEmpty() {
		defer begin.RUnlock()
		exchange.Result <-fmt.Errorf("RemoveArc: Begin node is out of context")
		return
	}
	end.RLock()
	if end.IsEmpty() {
		defer end.RUnlock()
		exchange.Result <-fmt.Errorf("RemoveArc: End node is out of context")
		return
	}
	arc_res := Exchange[*Node]{make(chan *Node), nil}
	arc_addr_res := Exchange[string]{make(chan string), nil}
	go context.GetArc(begin, end, el_type, &arc_res)
	go context.GetArcAddr(begin, end, el_type, &arc_addr_res)
	arc, arc_addr := <-arc_res.Result, <-arc_addr_res.Result
	if arc_addr != "" {
		arc.Lock()
		for hash, node := range arc.parent {
			node.Lock()
			if node.el_type&Node_t != 0 {
				delete(node.child, arc_addr)
				delete(arc.parent, hash)
			} else {
				for _, parent := range node.parent {
					go context.RemoveArc(parent, arc, node.el_type, nil)
				}
				for _, child := range node.child {
					go context.RemoveArc(arc, child, node.el_type, nil)
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
					go context.RemoveArc(parent, arc, node.el_type, nil)
				}
				for _, child := range node.child {
					go context.RemoveArc(arc, child, node.el_type, nil)
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
}
