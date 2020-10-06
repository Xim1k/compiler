package main

import (
  "fmt"
  // "time"
  // "os"
  // "strconv"
)

type AstTree struct {
  key string
  level int
	typ  itemType
	data  string
  text string
	parent  *AstTree
	childs []*AstTree
}

func indexOfInSlice(a []*AstTree, b AstTree) int {
    for i, n := range a {
        if b.key == n.key {
          return i
        }
    }

    return -1
}

func removeFromSlice(slice []*AstTree, s int) []*AstTree {
    return append(slice[:s], slice[s+1:]...)
}

func (tree *AstTree) addChild(child *AstTree) *AstTree {
  if (child.parent != nil) {
    if indexOfInSlice(child.parent.childs, *child) != -1 {
      child.parent.childs = removeFromSlice(child.parent.childs, indexOfInSlice(child.parent.childs, *child))
    }
  }

  if indexOfInSlice(tree.childs, *child) != -1 {
    tree.childs = removeFromSlice(tree.childs, indexOfInSlice(tree.childs, *child))
  }

  tree.childs = append(tree.childs, child)
  child.parent = tree

  return child
}

func (tree *AstTree) removeChild(child AstTree) {
  if indexOfInSlice(tree.childs, child) != -1 {
    tree.childs = removeFromSlice(tree.childs, indexOfInSlice(tree.childs, child))

    if child.parent == tree {
      child.parent = nil
    }
  }
}

func printTree(tree *AstTree) {
  if tree.typ == -1 {
    fmt.Println("[ Abstract syntax tree ]")
  }

  for index := 0; index < tree.level; index++ {
    fmt.Print("-")
  }

  if tree.typ != -1 {
    fmt.Print(">")

    if tree.typ == itemNode {
      fmt.Println("[ Type:", tree.text, "]")
    } else {
        fmt.Println("[ Type:", tree.text, ", value: '", tree.data, "' ] ")
    }
  }

  for _, child := range tree.childs {
      printTree(child)
  }
}
