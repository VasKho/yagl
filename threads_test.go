package yagl_test

import (
	"fmt"
	"testing"

	"github.com/VasKho/yagl"
)

func TestAddGorot(t *testing.T) {
	context := yagl.NewContext()
	go context.AddNode("test1", yagl.Node_t)
	go context.AddNode("test1", yagl.Node_t)
	go context.AddNode("test2", yagl.Node_t)
	go context.AddNode("test3", yagl.Node_t)
	go context.AddNode("test4", yagl.Node_t)
	go context.AddNode("test5", yagl.Node_t)
	fmt.Println(context.GetNodes())
}
