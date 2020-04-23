package main

import (
  "fmt"
)

func main()  {
  array := [10]int{2, 100, 5, 7, 11, 13, 101, 99, 124, 2}
  max := array[0]

  for i := 1; i < len(array); i++ {
    if array[i] > max {
      max = array[i]
    }
  }

  fmt.Println("Result", max)
}
