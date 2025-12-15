package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const BRANCHING = 4 // 2*t => t = 2

const MAX_KEYS = BRANCHING - 1
const MIN_KEYS = BRANCHING/2 - 1

var ErrAlreadyExists = errors.New("key Already exists")

type Node struct {
	keys     []int
	children []*Node
}

func (node *Node) length() int {
	return len(node.keys)
}

func (node *Node) isLeaf() bool {
	return len(node.children) == 0
}

func NewNode(keys []int, children []*Node) *Node {
	return &Node{keys, children}
}

func (node *Node) PrintInOrder() {
	n := len(node.keys)
	for i := 0; i < n; i++ {
		// Recurse on child i if not leaf
		if !node.isLeaf() && len(node.children) > i {
			node.children[i].PrintInOrder()
		}
		// Print current key
		fmt.Printf("%d ", node.keys[i])
	}
	// Finally, recurse on last child if not leaf
	if !node.isLeaf() && len(node.children) > n {
		node.children[n].PrintInOrder()
	}
}

func (node *Node) Print(level int) {
	// Print current node's keys with indentation
	fmt.Printf("%s", strings.Repeat("  ", level))
	fmt.Printf("%v\n", node.keys)

	// Recurse on children if not leaf
	if !node.isLeaf() {
		for _, child := range node.children {
			child.Print(level + 1)
		}
	}
}

func (parent *Node) split(childIndex int) {
	childToSplit := parent.children[childIndex]
	middleElement := childToSplit.keys[BRANCHING/2-1]
	parent.keys = slices.Insert(parent.keys, childIndex, middleElement)

	rightKeys := childToSplit.keys[BRANCHING/2:]
	childToSplit.keys = childToSplit.keys[:BRANCHING/2-1]

	newSibling := NewNode(rightKeys, []*Node{})
	if !childToSplit.isLeaf() {
		fmt.Println("childToSplit.children", childToSplit)
		newSibling.children = childToSplit.children[BRANCHING/2:]
		childToSplit.children = childToSplit.children[:BRANCHING/2]
	}
	parent.children = slices.Insert(parent.children, childIndex+1, newSibling)
}

func (root *Node) add(newItem int) {
	if root.isLeaf() {

	}

}

func (tree *Node) addRoot(newItem int) *Node {
	root := tree
	if tree.length() == MAX_KEYS {
		newRoot := NewNode([]int{}, []*Node{tree})
		newRoot.split(0)
		root = newRoot
	}

	// continue implementation

	return root
}

func (node Node) String() string {
	return fmt.Sprintf("{keys: (%v) - children: (%v)}", node.keys, node.children)
}

func main() {
	leafLeft := NewNode([]int{1, 2}, []*Node{})
	interLeft := NewNode([]int{50, 60, 70}, []*Node{leafLeft})
	leafRight := NewNode([]int{100, 200, 300}, []*Node{})
	root := NewNode([]int{85}, []*Node{interLeft, leafRight})
	root.PrintInOrder()
	fmt.Println(root)
	fmt.Println("============")
	root.split(0)
	root.PrintInOrder()
	fmt.Println(root)

	// fmt.Println("before adding", tree)
	// newRoot, _ := tree.add(40, nil)
	// fmt.Println("after adding", newRoot)
	// newRoot, _ = tree.add(68, nil)
	// fmt.Println("after adding", newRoot)
	// newRoot, _ = tree.add(666, nil)
	// fmt.Println("after adding", newRoot)
	// newRoot, _ = tree.add(69, nil)
	// fmt.Println("after adding", newRoot)
	// result, _ := tree.search(32)
	// fmt.Println(result)
	// result, _ = tree.search(46)
	// fmt.Println(result)
}
