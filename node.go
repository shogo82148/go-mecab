package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"

type Node struct {
	node *C.mecab_node_t
}

type NodeStat int

const (
	NormalNode  NodeStat = 0
	UnknownNode          = 1
	BOSNode              = 2
	EOSNode              = 3
	EONNode              = 4
)

func (stat NodeStat) String() string {
	switch stat {
	case NormalNode:
		return "Normal"
	case UnknownNode:
		return "Unknown"
	case BOSNode:
		return "BOS"
	case EOSNode:
		return "EOS"
	case EONNode:
		return "EON"
	}
	return ""
}

// Surface returns the surface string.
func (node Node) Surface() string {
	return C.GoStringN(node.node.surface, C.int(node.node.length))
}

// Length returns the length of the surface string.
func (node Node) Length() int {
	return int(node.node.length)
}

// RLength returns the length of the surface string including white space before the morph.
func (node Node) RLength() int {
	return int(node.node.rlength)
}

// Prev returns the previous Node.
func (node Node) Prev() Node {
	return Node{node: (*C.mecab_node_t)(node.node.prev)}
}

// Next returns the next Node.
func (node Node) Next() Node {
	return Node{node: (*C.mecab_node_t)(node.node.next)}
}

// ENext resturns a node which ends same position
func (node Node) ENext() Node {
	return Node{node: (*C.mecab_node_t)(node.node.enext)}
}

// ENext resturns a node which begins same position
func (node Node) BNext() Node {
	return Node{node: (*C.mecab_node_t)(node.node.bnext)}
}

// Stat returns the type of Node.
func (node Node) Stat() NodeStat {
	return NodeStat(node.node.stat)
}

// Id returns the id of Node.
func (node Node) Id() int {
	return int(node.node.id)
}

// RCAttr returns the right context attribute.
func (node Node) RCAttr() int {
	return int(node.node.rcAttr)
}

// LCAttr returns the right context attribute.
func (node Node) LCAttr() int {
	return int(node.node.lcAttr)
}

// CharType returns the character type.
func (node Node) CharType() int {
	return int(node.node.char_type)
}

// IsBest returns that if the Node is the best solution.
func (node Node) IsBest() bool {
	return node.node.isbest != 0
}

// Alpha returns the forward accumulative log summation.
func (node Node) Alpha() float32 {
	return float32(node.node.alpha)
}

// Beta returns the backward accumulative log summation.
func (node Node) Beta() float32 {
	return float32(node.node.beta)
}

// Prob returns the marginal probability.
func (node Node) Prob() float32 {
	return float32(node.node.prob)
}

// WCost returns word cost.
func (node Node) WCost() int {
	return int(node.node.wcost)
}

// Cost returns the best accumulative cost from bos node to this node.
func (node Node) Cost() int {
	return int(node.node.cost)
}
