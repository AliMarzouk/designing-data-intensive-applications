package main

import (
	"slices"
	"testing"
)

func TestSplitNode(t *testing.T) {
	tosplit := NewNode([]int{1, 2, 3, 4, 5}, nil)

	middleElement, left, right := splitNode(tosplit)

	if middleElement != 3 {
		t.Errorf("Expected %v Got %v", 3, middleElement)
	}
	if !slices.Equal(*left.keys, []int{1, 2}) {
		t.Errorf("Expected %v Got %v", []int{1, 2}, *left.keys)
	}
	if !slices.Equal(*right.keys, []int{4, 5}) {
		t.Errorf("Expected %v Got %v", []int{4, 5}, *right.keys)
	}

	tosplit = NewNode([]int{1, 2, 3, 4, 5, 6, 7}, nil)

	middleElement, left, right = splitNode(tosplit)

	if middleElement != 4 {
		t.Errorf("Expected %v Got %v", 4, middleElement)
	}
	if !slices.Equal(*left.keys, []int{1, 2, 3}) {
		t.Errorf("Expected %v Got %v", []int{1, 2, 3}, *left.keys)
	}
	if !slices.Equal(*right.keys, []int{5, 6, 7}) {
		t.Errorf("Expected %v Got %v", []int{5, 6, 7}, *right.keys)
	}
}

func TestInsertWithOrder(t *testing.T) {
	keys := []int{55, 78, 98}
	children := []Node{}

	newKeys, newChildren := insertWithOrder(keys, children, 100, nil)

	if !slices.Equal(newKeys, []int{55, 78, 98, 100}) {
		t.Errorf("Expected %v Got %v", []int{55, 78, 98, 100}, newKeys)
	}

	if !slices.Equal(newChildren, []Node{}) {
		t.Errorf("Expected %v Got %v", []Node{}, newChildren)
	}
}
