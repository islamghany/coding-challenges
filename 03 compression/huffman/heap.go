package huffman

import "fmt"

type HeapNode struct {
	Char  string
	Count int
	Left  *HeapNode
	Right *HeapNode
}

type HuffmanTree struct {
	Root *HeapNode
}

func NewNode(char string, count int) *HeapNode {
	return &HeapNode{
		Char:  char,
		Count: count,
	}
}

func NewHuffmanTree(root *HeapNode) *HuffmanTree {
	return &HuffmanTree{
		Root: root,
	}
}

type Heap struct {
	Nodes []*HeapNode
}

func NewHeapNode(char string, count int) *HeapNode {
	return &HeapNode{
		Char:  char,
		Count: count,
	}
}

func NewHeap() *Heap {
	return &Heap{
		Nodes: []*HeapNode{},
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
func (h *Heap) Insert(node *HeapNode) {
	h.Nodes = append(h.Nodes, node)
	// heapify the node up
	h.HeapifyUp(len(h.Nodes) - 1)
}

func (h *Heap) ExtractMin() *HeapNode {
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
