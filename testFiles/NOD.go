package main

/*
  COMMENT
*/
//asdasd

import (
  "fmt"
)

func main()  {
    var bigger int = 6538
    var smaller int = 1547
    var tmp int = 0
    var result int = 0

    for result == 0 {
        if bigger < smaller {
            tmp = smaller
            smaller = bigger
            bigger = tmp
        }

        if bigger % smaller != 0 {
            bigger = bigger % smaller
        } else {
            result = smaller
        }
    }

    fmt.Println("Result", result)
}
