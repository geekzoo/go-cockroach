package main

import (
    "fmt"
    "sync"
    "time"
    "container/list"
)

type Counter struct {
	count int
}
var wg = sync.WaitGroup{}


var total int

func (self Counter) currentValue() int {
	return self.count
}
func (self *Counter) increment() {
	self.count++
}

func proc(a int) {
ticker := time.NewTicker(1 * time.Second)
//counts := make(chan int)
quit := make(chan struct{})
go func() {
    for {
       select {
        case <- ticker.C:
            l := list.New()
            l.InsertAfter(a, e1)
            
            fmt.Printf("\t\t\033[37mCOUNTER: %v\033[0m\n", l)
            total = 0
        case <- quit:
            ticker.Stop()
            return
        }
    }
 }()
}
/*func main() {
	counter := Counter{1}
	counter.increment()
	counter.increment()

	fmt.Printf("current value %d", counter.currentValue())
}*/
