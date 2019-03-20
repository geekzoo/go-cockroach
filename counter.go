package main

//import "fmt"

type Counter struct {
	count int
}

func (self Counter) currentValue() int {
	return self.count
}
func (self *Counter) increment() {
	self.count++
}

/*func main() {
	counter := Counter{1}
	counter.increment()
	counter.increment()

	fmt.Printf("current value %d", counter.currentValue())
}*/