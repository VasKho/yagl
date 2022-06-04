package yagl_test

import (
	"github.com/yagl"
	"testing"
)

func TestNewContext(t *testing.T) {
	context := yagl.NewContext()
	if context == nil {
		t.Fail()
	}
}

func TestAddNode(t *testing.T) {
	context := yagl.NewContext()
	context.AddNode("test", yagl.Node_t)
	nodes := context.GetNodes()
	hash := context.GetNodeAddr("test", yagl.Node_t)
	if context.GetSize() != 1 {
		t.Fail()
	}
	if _, ok := nodes[hash]; !ok {
		t.Fail()
	}
}

func TestAddSameNode(t *testing.T) {
	context := yagl.NewContext()
	context.AddNode("test", yagl.Node_t)
	err := context.AddNode("test", yagl.Node_t)
	if context.GetSize() != 1 {
		t.Fail()
	}
	if err == nil {
		t.Fail()
	}
}

func TestAddSameNodeDiffType(t *testing.T) {
	context := yagl.NewContext()
	context.AddNode("test", yagl.Node_t)
	err := context.AddNode("test", yagl.Node_t|yagl.Const_t)
	if err != nil {
		t.Fail()
	}
}

func TestAddNoIdtfNode(t *testing.T) {
	context := yagl.NewContext()
	err := context.AddNode("", yagl.Node_t)
	if err != nil {
		t.Fail()
	}
}

func TestGetNode(t *testing.T) {
	context := yagl.NewContext()
	context.AddNode("test", yagl.Node_t)
	context.AddNode("test1", yagl.Node_t)
	context.AddNode("test2", yagl.Node_t)
	nodes := context.GetNodes()
	hash := context.GetNodeAddr("test1", yagl.Node_t)
	test_node := context.GetNode("test1", yagl.Node_t)
	if !test_node.IsEqual(*nodes[hash]) {
		t.Fail()
	}
}

func TestGetNullNode(t *testing.T) {
	context := yagl.NewContext()
	context.AddNode("test", yagl.Node_t)
	context.AddNode("test1", yagl.Node_t)
	context.AddNode("test2", yagl.Node_t)
	test_node := context.GetNode("test3", yagl.Node_t)
	if !test_node.IsEmpty() {
		t.Fail()
	}
}

func TestAddArc(t *testing.T) {
	context := yagl.NewContext()
	context.AddNode("a", yagl.Node_t)
	context.AddNode("b", yagl.Node_t)
	a := context.GetNode("a", yagl.Node_t)
	b := context.GetNode("b", yagl.Node_t)
	err := context.AddArc(a, b, yagl.Arc_t)
	if err != nil {
		t.Fail()
	}
}

func TestAddSameArc(t *testing.T) {
	context := yagl.NewContext()
	context.AddNode("a", yagl.Node_t)
	context.AddNode("b", yagl.Node_t)
	a := context.GetNode("a", yagl.Node_t)
	b := context.GetNode("b", yagl.Node_t)
	context.AddArc(a, b, yagl.Arc_t)
	err := context.AddArc(a, b, yagl.Arc_t)
	if err != nil {
		t.Fail()
	}
}

func TestAddNodeTypeArc(t *testing.T) {
	context := yagl.NewContext()
	context.AddNode("a", yagl.Node_t)
	context.AddNode("b", yagl.Node_t)
	a := context.GetNode("a", yagl.Node_t)
	b := context.GetNode("b", yagl.Node_t)
	err := context.AddArc(a, b, yagl.Node_t|yagl.Const_t)
	if err == nil {
		t.Fail()
	}
}

func TestRemoveArc(t *testing.T) {
	context := yagl.NewContext()
	context.AddNode("a", yagl.Node_t)
	context.AddNode("b", yagl.Node_t)
	context.AddNode("c", yagl.Node_t|yagl.Role_t)
	a := context.GetNode("a", yagl.Node_t)
	b := context.GetNode("b", yagl.Node_t)
	c := context.GetNode("c", yagl.Node_t|yagl.Role_t)
	context.AddArc(a, b, yagl.Arc_t)
	arc := context.GetArc(a, b, yagl.Arc_t)
	context.AddArc(c, arc, yagl.Arc_t)
	err := context.RemoveArc(a, b, yagl.Arc_t)
	if err != nil {
		t.Fail()
	}
}

func TestRemoveNode(t *testing.T) {
	context := yagl.NewContext()
	context.AddNode("a", yagl.Node_t)
	context.AddNode("b", yagl.Node_t)
	context.AddNode("c", yagl.Node_t|yagl.Role_t)
	a := context.GetNode("a", yagl.Node_t)
	b := context.GetNode("b", yagl.Node_t)
	c := context.GetNode("c", yagl.Node_t|yagl.Role_t)
	context.AddArc(a, b, yagl.Arc_t)
	arc := context.GetArc(a, b, yagl.Arc_t)
	context.AddArc(c, arc, yagl.Arc_t)
	err := context.RemoveNode("b", yagl.Node_t)
	if err != nil {
		t.Fail()
	}
}
