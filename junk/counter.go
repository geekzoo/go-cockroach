package main

import (
    "fmt"
 //   "sync"
    "time"
//    "container/list"
)


var f_id int
func proc(a int) {
ticker := time.NewTicker(5 * time.Second)
quit := make(chan struct{})
m := make(map[int]int)

go func() {
    for {
       select {
        case <- ticker.C:
  f_id++
  m[f_id] = a
  fmt.Println("QPS:", len(m), m)

  break

        case <- quit:
            ticker.Stop()
            return
        }
    }
 }()
}
