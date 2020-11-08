package main

import (
  "fmt"
)

func main() {
  var array = [10]int{2, 100, 5, 7, 11, 13, 101, 99, 124, 2}
  var i int = 1
  var max int = array[0]

  for i < len(array) {
    if array[i] > max {
      max = array[i]
      i = i + 1
    }
  }

  fmt.Println("Result", max)
}
