package huffman

import "fmt"

// =======================================================================
// huffman tree definition

// define the node the huffman tree
type Node struct {
	Char  rune // rune is an alias for int32 and is equivalent to int32 in all ways. It is used, by convention, to distinguish character values from integer values.
	Count int
	Left  *Node
	Right *Node
}

// isLeaf() returns true if the node is a leaf node, i.e., it has no children.
func (n *Node) isLeaf() bool {
	return n.Left == nil && n.Right == nil
}

type HuffmanTree struct {
	Root *Node
}

func NewNode(char rune, count int) *Node {
	return &Node{
		Char:  char,
		Count: count,
	}
}

func NewHuffmanTree(root *Node) *HuffmanTree {
	return &HuffmanTree{
		Root: root,
	}
}

// =======================================================================
// define the heap data structure
type Heap struct {
	Nodes []*Node
}

func NewHeap() *Heap {
	return &Heap{
		Nodes: []*Node{},
	}
}

func (h *Heap) leftChild(idx int) int {
	return 2*idx + 1
}

func (h *Heap) rightChild(idx int) int {
	return 2*idx + 2
}

func (h *Heap) parent(idx int) int {
	return (idx - 1) / 2
}

func (h *Heap) Size() int {
	return len(h.Nodes)
}

func (h *Heap) HeapifyUp(idx int) {
	// swap the node with it's parent until the parent is less than the node or it's the root
	for h.parent(idx) >= 0 && h.Nodes[h.parent(idx)].Count > h.Nodes[idx].Count {
		h.Nodes[h.parent(idx)], h.Nodes[idx] = h.Nodes[idx], h.Nodes[h.parent(idx)]
		idx = h.parent(idx)
	}
}

func (h *Heap) HeapifyDown(idx int) {
	min := idx
	l := h.leftChild(idx)
	r := h.rightChild(idx)
	// swap the node the smallest child of it, and repeat the process until
	// we the node is samller that it's children or it is a leaf
	if l < h.Size() && h.Nodes[l].Count < h.Nodes[min].Count {
		min = l
	}
	if r < h.Size() && h.Nodes[r].Count < h.Nodes[min].Count {
		min = r
	}
	if min != idx {
		h.Nodes[min], h.Nodes[idx] = h.Nodes[idx], h.Nodes[min]
		h.HeapifyDown(min)
	}
}
func (h *Heap) Insert(node *Node) {
	h.Nodes = append(h.Nodes, node)
	// heapify the node up
	h.HeapifyUp(len(h.Nodes) - 1)
}

func (h *Heap) ExtractMin() *Node {
	// if there is no node return nil
	if h.Size() == 0 {
		return nil
	}
	// swap the root node with the last index node and delete the last node
	min := h.Nodes[0]
	h.Nodes[0] = h.Nodes[h.Size()-1]
	h.Nodes = h.Nodes[:h.Size()-1]
	// heapify the node down
	h.HeapifyDown(0)
	return min

}

func (h *Heap) Print() {
	for _, node := range h.Nodes {
		fmt.Printf("%s: %d\n", node.Char, node.Count)
	}
}
