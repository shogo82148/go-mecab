package mecab

// #include <mecab.h>
// #include <stdlib.h>
import "C"

// Node is a node in a lattice.
type Node struct {
	node *C.mecab_node_t

	// actual data of node is stored in mecab or lattice.
	// they are here to avoid garbage collection.
	mecab   *mecab
	lattice *lattice
}

// NodeStat is status of a node.
type NodeStat int

const (
	// NormalNode is status for normal node.
	NormalNode NodeStat = 0

	// UnknownNode is status for unknown node.
	UnknownNode NodeStat = 1

	// BOSNode is status for BOS(Begin Of Sentence) node.
	BOSNode NodeStat = 2

	// EOSNode is status for EOS(End Of Sentence) node.
	EOSNode NodeStat = 3

	// EONNode is status for EON(End Of Node) node.
	EONNode NodeStat = 4
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

// String returns Surface and Feature
func (node Node) String() string {
	return node.Surface() + "\t" + node.Feature()
}

// Surface returns the surface string.
func (node Node) Surface() string {
	return C.GoStringN(node.node.surface, C.int(node.node.length))
}

// Feature returns the feature.
func (node Node) Feature() string {
	return C.GoString(node.node.feature)
}

// Length returns the length of the surface string.
func (node Node) Length() int {
	return int(node.node.length)
}

// RLength returns the length of the surface string including white space before the morph.
func (node Node) RLength() int {
	return int(node.node.rlength)
}

// PosID returns the part-of-speech id.
func (node Node) PosID() int {
	return int(node.node.posid)
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
	return Node{
		node:    (*C.mecab_node_t)(node.node.enext),
		mecab:   node.mecab,
		lattice: node.lattice,
	}
}

// BNext resturns a node which begins same position
func (node Node) BNext() Node {
	return Node{
		node:    (*C.mecab_node_t)(node.node.bnext),
		mecab:   node.mecab,
		lattice: node.lattice,
	}
}

// Stat returns the type of Node.
func (node Node) Stat() NodeStat {
	return NodeStat(node.node.stat)
}

// ID returns the id of Node.
func (node Node) ID() int {
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

// IsZero returns whether the node is zero.
func (node Node) IsZero() bool {
	return node.node == nil
}
