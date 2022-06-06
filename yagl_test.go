package yagl_test

import (
	"sync"
	"testing"

	"github.com/VasKho/yagl"
)

func TestNewContext(t *testing.T) {
	context := yagl.NewContext()
	if context == nil {
		t.Fail()
	}
}

func TestAddNode(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 2), &sync.WaitGroup{}}
	add_res.Wg.Add(1)
	go context.AddNode("test", yagl.Node_t, &add_res)
	exch := yagl.Exchange[yagl.Arcs]{make(chan yagl.Arcs), nil}
	add_res.Wg.Wait()
	go context.GetNodes(&exch)
	nodes := <-exch.Result
	get_addr_res := yagl.Exchange[string]{make(chan string), nil}
	go context.GetNodeAddr("test", yagl.Node_t, &get_addr_res)
	if context.GetSize() != 1 {
		t.Fail()
	}
	if _, ok := nodes[<-get_addr_res.Result]; !ok {
		t.Fail()
	}
}

func TestAddSameNode(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 2), &sync.WaitGroup{}}
	add_res.Wg.Add(1)
	go context.AddNode("test", yagl.Node_t, &add_res)
	add_res.Wg.Wait()
	add_res.Wg.Add(1)
	go context.AddNode("test", yagl.Node_t, &add_res)
	err := <-add_res.Result
	if context.GetSize() != 1 {
		t.Fail()
	}
	if err == nil {
		t.Fail()
	}
}

func TestAddSameNodeDiffType(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 2), &sync.WaitGroup{}}
	add_res.Wg.Add(1)
	go context.AddNode("test", yagl.Node_t, &add_res)
	add_res.Wg.Wait()
	add_res.Wg.Add(1)
	go context.AddNode("test", yagl.Node_t|yagl.Const_t, &add_res)
	add_res.Wg.Wait()
	if len(add_res.Result) != 0 {
		t.Fail()
	}
	get_addr_res := yagl.Exchange[string]{make(chan string, 2), &sync.WaitGroup{}}
	get_addr_res.Wg.Add(2)
	go context.GetNodeAddr("test", yagl.Node_t, &get_addr_res)
	go context.GetNodeAddr("test", yagl.Node_t|yagl.Const_t, &get_addr_res)
	get_addr_res.Wg.Wait()
	if <-get_addr_res.Result == <-get_addr_res.Result {
		t.Fail()
	}
}

func TestAddNoIdtfNode(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 2), &sync.WaitGroup{}}
	add_res.Wg.Add(1)
	go context.AddNode("", yagl.Node_t, &add_res)
	add_res.Wg.Wait()
	if len(add_res.Result) != 0 {
		t.Fail()
	}
}

func TestGetNode(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 3), &sync.WaitGroup{}}
	add_res.Wg.Add(3)
	go context.AddNode("test", yagl.Node_t, &add_res)
	go context.AddNode("test1", yagl.Node_t, &add_res)
	go context.AddNode("test2", yagl.Node_t, &add_res)
	add_res.Wg.Wait()
	get_res := yagl.Exchange[yagl.Arcs]{make(chan yagl.Arcs, 2), &sync.WaitGroup{}}
	get_res.Wg.Add(1)
	go context.GetNodes(&get_res)
	get_res.Wg.Wait()
	node_addr_res := yagl.Exchange[string]{make(chan string), nil}
	get_node_res := yagl.Exchange[*yagl.Node]{make(chan *yagl.Node), nil}
	go context.GetNodeAddr("test1", yagl.Node_t, &node_addr_res)
	hash := <-node_addr_res.Result
	go context.GetNode("test1", yagl.Node_t, &get_node_res)
	test_node := <-get_node_res.Result
	if !test_node.IsEqual(*(<-get_res.Result)[hash]) {
		t.Fail()
	}
}

func TestGetNullNode(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 3), &sync.WaitGroup{}}
	add_res.Wg.Add(3)
	go context.AddNode("test", yagl.Node_t, &add_res)
	go context.AddNode("test1", yagl.Node_t, &add_res)
	go context.AddNode("test2", yagl.Node_t, &add_res)
	add_res.Wg.Wait()
	get_node_res := yagl.Exchange[*yagl.Node]{make(chan *yagl.Node), nil}
	go context.GetNode("test3", yagl.Node_t, &get_node_res)
	test_node := <-get_node_res.Result
	if !test_node.IsEmpty() {
		t.Fail()
	}
}

func TestAddArc(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 2), &sync.WaitGroup{}}
	add_res.Wg.Add(2)
	go context.AddNode("a", yagl.Node_t, &add_res)
	go context.AddNode("b", yagl.Node_t, &add_res)
	add_res.Wg.Wait()
	get_res := yagl.Exchange[*yagl.Node]{make(chan *yagl.Node), nil}
	go context.GetNode("a", yagl.Node_t, &get_res)
	a := <-get_res.Result
	go context.GetNode("b", yagl.Node_t, &get_res)
	b := <-get_res.Result
	add_res.Wg.Add(1)
	go context.AddArc(a, b, yagl.Arc_t, &add_res)
	add_res.Wg.Wait()
	if len(add_res.Result) != 0 || context.GetSize() != 3 {
		t.Fail()
	}
}

func TestGetArc(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 2), &sync.WaitGroup{}}
	add_res.Wg.Add(2)
	go context.AddNode("a", yagl.Node_t, &add_res)
	go context.AddNode("b", yagl.Node_t, &add_res)
	add_res.Wg.Wait()
	get_res := yagl.Exchange[*yagl.Node]{make(chan *yagl.Node), nil}
	go context.GetNode("a", yagl.Node_t, &get_res)
	a := <-get_res.Result
	go context.GetNode("b", yagl.Node_t, &get_res)
	b := <-get_res.Result
	add_res.Wg.Add(1)
	go context.AddArc(a, b, yagl.Arc_t, &add_res)
	add_res.Wg.Wait()
	go context.GetArc(a, b, yagl.Arc_t, &get_res)
	if (<-get_res.Result).IsEmpty() {
		t.Fail()
	}
}

func TestAddSameArc(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 2), &sync.WaitGroup{}}
	add_res.Wg.Add(2)
	go context.AddNode("a", yagl.Node_t, &add_res)
	go context.AddNode("b", yagl.Node_t, &add_res)
	add_res.Wg.Wait()
	get_res := yagl.Exchange[*yagl.Node]{make(chan *yagl.Node), nil}
	go context.GetNode("a", yagl.Node_t, &get_res)
	a := <-get_res.Result
	go context.GetNode("b", yagl.Node_t, &get_res)
	b := <-get_res.Result
	add_res.Wg.Add(1)
	go context.AddArc(a, b, yagl.Arc_t, &add_res)
	add_res.Wg.Wait()
	add_res.Wg.Add(1)
	go context.AddArc(a, b, yagl.Arc_t, &add_res)
	add_res.Wg.Wait()
	if len(add_res.Result) != 0 || context.GetSize() != 3 {
		t.Fail()
	}
}

func TestAddArcTypeNode(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 2), &sync.WaitGroup{}}
	add_res.Wg.Add(2)
	go context.AddNode("a", yagl.Node_t, &add_res)
	go context.AddNode("b", yagl.Node_t, &add_res)
	add_res.Wg.Wait()
	get_res := yagl.Exchange[*yagl.Node]{make(chan *yagl.Node), nil}
	go context.GetNode("a", yagl.Node_t, &get_res)
	a := <-get_res.Result
	go context.GetNode("b", yagl.Node_t, &get_res)
	b := <-get_res.Result
	add_res.Wg.Add(1)
	go context.AddArc(a, b, yagl.Node_t|yagl.Const_t, &add_res)
	add_res.Wg.Wait()
	if len(add_res.Result) == 0 || context.GetSize() != 2 {
		t.Fail()
	}
}

func TestRemoveArc(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 3), &sync.WaitGroup{}}
	add_res.Wg.Add(3)
	go context.AddNode("a", yagl.Node_t, &add_res)
	go context.AddNode("b", yagl.Node_t, &add_res)
	go context.AddNode("c", yagl.Node_t|yagl.Role_t, &add_res)
	add_res.Wg.Wait()
	get_res := yagl.Exchange[*yagl.Node]{make(chan *yagl.Node), nil}
	go context.GetNode("a", yagl.Node_t, &get_res)
	a := <-get_res.Result
	go context.GetNode("b", yagl.Node_t, &get_res)
	b := <-get_res.Result
	go context.GetNode("c", yagl.Node_t|yagl.Role_t, &get_res)
	c := <-get_res.Result
	add_res.Wg.Add(1)
	go context.AddArc(a, b, yagl.Arc_t, &add_res)
	add_res.Wg.Wait()
	get_res = yagl.Exchange[*yagl.Node]{make(chan *yagl.Node), &sync.WaitGroup{}}
	get_res.Wg.Add(1)
	go context.GetArc(a, b, yagl.Arc_t, &get_res)
	arc := <-get_res.Result
	add_res.Wg.Add(1)
	go context.AddArc(c, arc, yagl.Arc_t, &add_res)
	add_res.Wg.Wait()
	add_res.Wg.Add(1)
	go context.RemoveArc(a, b, yagl.Arc_t, &add_res)
	add_res.Wg.Wait()
	if len(add_res.Result) != 0 {
		t.Fail()
	}
}

func TestRemoveNode(t *testing.T) {
	context := yagl.NewContext()
	add_res := yagl.Exchange[error]{make(chan error, 3), &sync.WaitGroup{}}
	add_res.Wg.Add(3)
	go context.AddNode("a", yagl.Node_t, &add_res)
	go context.AddNode("b", yagl.Node_t, &add_res)
	go context.AddNode("c", yagl.Node_t|yagl.Role_t, &add_res)
	add_res.Wg.Wait()
	get_res := yagl.Exchange[*yagl.Node]{make(chan *yagl.Node), nil}
	go context.GetNode("a", yagl.Node_t, &get_res)
	a := <-get_res.Result
	go context.GetNode("b", yagl.Node_t, &get_res)
	b := <-get_res.Result
	go context.GetNode("c", yagl.Node_t|yagl.Role_t, &get_res)
	c := <-get_res.Result
	add_res.Wg.Add(1)
	go context.AddArc(a, b, yagl.Arc_t, &add_res)
	add_res.Wg.Wait()
	get_res = yagl.Exchange[*yagl.Node]{make(chan *yagl.Node), &sync.WaitGroup{}}
	get_res.Wg.Add(1)
	go context.GetArc(a, b, yagl.Arc_t, &get_res)
	arc := <-get_res.Result
	add_res.Wg.Add(1)
	go context.AddArc(c, arc, yagl.Arc_t, &add_res)
	add_res.Wg.Wait()
	add_res.Wg.Add(1)
	go context.RemoveNode("b", yagl.Node_t, &add_res)
	add_res.Wg.Wait()
	if len(add_res.Result) != 0 {
		t.Fail()
	}
}
