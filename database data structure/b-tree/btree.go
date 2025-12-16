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
		// Recurse on child i if not leaf
		if !node.isLeaf && len(node.children) > i {
			node.children[i].PrintInOrder()
		}
		// Print current key
		fmt.Printf("%d ", node.keys[i])
	}
	// Finally, recurse on last child if not leaf
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
	// Print current node's keys with indentation
	fmt.Printf("%s", strings.Repeat("  ", level))
	fmt.Printf("%v\n", node.keys)

	// Recurse on children if not leaf
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
	parent.keys = slices.Insert(parent.keys, childIndex, middleElement)

	newKeys := childToSplit.keys[BRANCHING/2:]
	childToSplit.keys = childToSplit.keys[:BRANCHING/2-1]
	childToSplit.length = BRANCHING/2 - 1

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
	newChild := NewNode(newKeys, newChildren, len(newKeys), newIsLeaf) // determine is leaf value

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
			}
		}
		if nextChildPos == -1 {
			nextChildPos = BRANCHING - 1
		}

		if node.children[nextChildPos].length == MAX_KEYS {
			node.split(nextChildPos)
			if node.keys[nextChildPos] < newKey {
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
	// leafLeft := NewNode([]int{30, 40}, [BRANCHING]*Node{}, 2, true)
	// leafMiddle := NewNode([]int{65, 67}, [BRANCHING]*Node{}, 2, true)
	// leafMiddleRight := NewNode([]int{80, 81}, [BRANCHING]*Node{}, 2, true)
	// interLeft := NewNode([]int{50, 60, 70}, [BRANCHING]*Node{leafLeft, nil, leafMiddle, leafMiddleRight}, 3, false)
	// leafRight := NewNode([]int{100, 200, 300}, [BRANCHING]*Node{}, 3, true)
	// root := NewNode([]int{2, 85}, [BRANCHING]*Node{nil, interLeft, leafRight}, 1, false)
	// root.PrintInOrder()
	// fmt.Println()
	// root.Print(0)
	// fmt.Println("============")
	// root.split(1)
	// root.PrintInOrder()
	// fmt.Println()
	// root.Print(0)
	// fmt.Println("============")

	test := NewNode([]int{30, 40, 50}, [BRANCHING]*Node{}, 3, true)
	test = test.addRoot(3)
	test = test.addRoot(4)
	test = test.addRoot(5)
	test.Print(0)
	test.PrintInOrder()
}
