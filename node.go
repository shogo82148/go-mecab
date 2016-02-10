package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"

type Node struct {
	node *C.mecab_node_t
}

// surface string
func (node *Node) Surface() string {
	return C.GoStringN(node.node.surface, C.int(node.node.length))
}

func (node *Node) Next() *Node {
	next := node.node.next
	if next == nil {
		return nil
	}
	return &Node{node: next}
}

func (node *Node) Length() int {
	return int(node.node.length)
}
