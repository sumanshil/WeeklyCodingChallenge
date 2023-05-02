package main

import (
	"compression_tool/priority_queue"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type HuffBaseNode interface {
	IsLeaf() bool
	Weight() int
}

type HuffLeafNode struct {
	element string
	weight  int
}

func NewLeafNode(element string, weight int) *HuffLeafNode {
	return &HuffLeafNode{element: element, weight: weight}
}

func (node *HuffLeafNode) IsLeaf() bool {
	return true
}

func (node *HuffLeafNode) Weight() int {
	return node.weight
}

func (node *HuffLeafNode) Value() int {
	return node.weight
}

type HuffInternalNode struct {
	weight int
	left   HuffBaseNode
	right  HuffBaseNode
}

func NewInternalNode(weight int, left HuffBaseNode, right HuffBaseNode) *HuffInternalNode {
	return &HuffInternalNode{weight: weight, left: left, right: right}
}

func (node *HuffInternalNode) IsLeaf() bool {
	return false
}

func (node *HuffInternalNode) Weight() int {
	return node.weight
}

func (node *HuffInternalNode) Value() int {
	return node.weight
}

type HuffTree struct {
	root HuffBaseNode
}

func (node *HuffTree) Weight() int {
	return node.root.Weight()
}

func (node *HuffTree) Value() int {
	return node.root.Weight()
}

func NewHuffTreeWithLeafNode(el string, weight int) *HuffTree {
	return &HuffTree{
		root: &HuffLeafNode{
			element: el,
			weight:  weight,
		},
	}
}

func NewHuffTreeWithInternalNode(left HuffBaseNode,
	right HuffBaseNode, weight int) *HuffTree {
	return &HuffTree{
		root: &HuffInternalNode{
			left:   left,
			right:  right,
			weight: weight,
		},
	}
}

func assignPrefixes(baseNode HuffBaseNode, prefixMap map[string]string, prefixes []string) {
	v, ok := baseNode.(*HuffLeafNode)
	if ok {
		prefixMap[v.element] = strings.Join(prefixes, "")
		return
	}
	internal, _ := baseNode.(*HuffInternalNode)
	prefixes = append(prefixes, "0")
	assignPrefixes(internal.left, prefixMap, prefixes)
	prefixes = prefixes[:len(prefixes)-1]
	prefixes = append(prefixes, "1")
	assignPrefixes(internal.right, prefixMap, prefixes)
	prefixes = prefixes[:len(prefixes)-1]
}

func main() {
	fmt.Println("Hello world")
	data, _ := os.ReadFile("./data/data.txt")
	//	fmt.Print(string(data))
	charFrequency := make(map[string]int)
	for _, char := range data {
		val, ok := charFrequency[string(char)]
		if ok {
			charFrequency[string(char)] = val + 1
		} else {
			charFrequency[string(char)] = 1
		}
	}
	//treeLeafNodes := make([]*HuffTree, 0)
	priorityQueue := priority_queue.NewPriorityQueue()
	for str, frequency := range charFrequency {
		//treeLeafNodes = append(treeLeafNodes, NewHuffTreeWithLeafNode(str, frequency))
		priorityQueue.Insert(NewHuffTreeWithLeafNode(str, frequency))
		println(str + " ===> " + strconv.Itoa(frequency))
	}
	var tree3 *HuffTree
	for priorityQueue.Size() > 1 {
		node1 := priorityQueue.Delete()
		tree1 := node1.(*HuffTree)
		node2 := priorityQueue.Delete()
		tree2 := node2.(*HuffTree)
		tree3 = NewHuffTreeWithInternalNode(tree1.root, tree2.root,
			tree1.Weight()+tree2.Weight())
		priorityQueue.Insert(tree3)
	}
	println(tree3.Weight())
	prefixMap := make(map[string]string)
	prefixes := make([]string, 0)
	assignPrefixes(tree3.root, prefixMap, prefixes)
	for str, prefix := range prefixMap {
		println(str)
		println(prefix)
	}
}
