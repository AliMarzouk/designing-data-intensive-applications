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
	children [BRANCHING]*Node
	length   int
	isLeaf   bool
}

func NewNode(keys []int, children [BRANCHING]*Node, length int, isLeaf bool) *Node {
	return &Node{keys, children, length, isLeaf}
}

func (node *Node) PrintInOrder() {
	if node == nil {
		return
	}
	n := len(node.keys)
	for i := 0; i < n; i++ {
		if !node.isLeaf && len(node.children) > i {
			node.children[i].PrintInOrder()
		}
		fmt.Printf("%d ", node.keys[i])
	}
	if !node.isLeaf && len(node.children) > n {
		node.children[n].PrintInOrder()
	}
}

func (node *Node) Print(level int) {
	if node == nil {
		fmt.Printf("%s", strings.Repeat("  ", level))
		fmt.Println("<nil>")
		return
	}
	fmt.Printf("%s", strings.Repeat("  ", level))
	fmt.Printf("%v [length=%v; isLeaf=%v]\n", node.keys, node.length, node.isLeaf)

	if !node.isLeaf {
		for _, child := range node.children {
			child.Print(level + 1)
		}
	}
}

func (parent *Node) split(childIndex int) {
	childToSplit := parent.children[childIndex]
	if childToSplit == nil {
		panic("cannot split nil child")
	}
	middleElement := childToSplit.keys[BRANCHING/2-1]
	if childIndex == MAX_KEYS {
		parent.keys = append(parent.keys, middleElement)
	} else {
		parent.keys = slices.Insert(parent.keys, childIndex, middleElement)
	}
	parent.length++

	newKeys := make([]int, MIN_KEYS)
	copy(newKeys, childToSplit.keys[BRANCHING/2:])
	childToSplit.keys = childToSplit.keys[:BRANCHING/2-1]
	childToSplit.length = len(newKeys)

	var newChildren = [BRANCHING]*Node{}
	newIsLeaf := true
	if !childToSplit.isLeaf {
		for i := BRANCHING / 2; i < BRANCHING; i++ {
			newChildren[i-BRANCHING/2] = childToSplit.children[i] // take children form position childIndex to end in the childToSplit and put them in the newChild.children
			childToSplit.children[i] = nil                        // remove children from childToSplit because they have been moved to the newChild
			if newChildren[i-BRANCHING/2] != nil {
				newIsLeaf = false
			}
		}
	}
	newChild := NewNode(newKeys, newChildren, len(newKeys), newIsLeaf)

	insertChild(&parent.children, childIndex+1, newChild) // insert newChild at the parent.children correct position in parent
}

func insertChild(children *[BRANCHING]*Node, index int, newNode *Node) {
	var currentNode *Node
	for i := index; i < len(*children); i++ {
		currentNode = (*children)[i]
		(*children)[i] = newNode
		newNode = currentNode
	}
}

func (node *Node) add(newKey int) {
	if node.isLeaf {
		for i := 0; i < len(node.keys); i++ {
			if node.keys[i] == newKey {
				panic(ErrAlreadyExists)
			}

			if node.keys[i] > newKey {
				node.keys = slices.Insert(node.keys, i, newKey)
				node.length += 1
				return
			}
		}
		node.keys = append(node.keys, newKey)
		node.length = node.length + 1
	} else {
		var nextChildPos int = -1
		for i := 0; i < len(node.keys); i++ {
			if node.keys[i] == newKey {
				panic(ErrAlreadyExists)
			}

			if node.keys[i] > newKey {
				nextChildPos = i
				break
			}
		}
		if nextChildPos == -1 {
			nextChildPos = node.length
		}

		if node.children[nextChildPos] == nil {
			node.children[nextChildPos] = NewNode([]int{}, [BRANCHING]*Node{}, 0, true)
			node.isLeaf = false
		}

		if node.children[nextChildPos].length == MAX_KEYS {
			node.split(nextChildPos)
			if nextChildPos < len(node.keys) && node.keys[nextChildPos] < newKey {
				nextChildPos += 1
			}
		}

		node.children[nextChildPos].add(newKey)
	}
}

func (tree *Node) addRoot(newItem int) *Node {
	root := tree
	if tree.length == MAX_KEYS {
		newRoot := NewNode([]int{}, [BRANCHING]*Node{tree}, 0, false)
		newRoot.split(0)
		root = newRoot
	}

	root.add(newItem)

	return root
}

func (node Node) String() string {
	return fmt.Sprintf("{keys: (%v) - children: (%v)}", node.keys, node.children)
}

func main() {

	root := NewNode([]int{}, [BRANCHING]*Node{}, 0, true)
	root = root.addRoot(2)
	root = root.addRoot(60)
	root = root.addRoot(85)
	root = root.addRoot(50)
	root = root.addRoot(70)
	root = root.addRoot(30)
	root = root.addRoot(40)
	root = root.addRoot(100)
	root = root.addRoot(200)
	root = root.addRoot(65)
	root = root.addRoot(75)
	root = root.addRoot(64)

	root = root.addRoot(73)
	root = root.addRoot(92)
	root = root.addRoot(13)
	root = root.addRoot(18)
	root = root.addRoot(24)
	root = root.addRoot(44)

	root.Print(0)
	root.PrintInOrder()
}
