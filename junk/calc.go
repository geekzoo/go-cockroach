/*
 * Just junk atm Idea is to count qps and report to stdout and carbon/influx and elasticsearch
 * not sure how I want this to work so this is a junk function atm.
 */
package main
import (
  "fmt"
  "time"
)

func qps(now ...int) (int) {

    if now != nil {

	fmt.Printf("qps: %v\n", len(now))
    }
  return len(now)
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func helloworld(t time.Time) {
	qps()
	fmt.Printf("%v: Hello, World!\n", t)
}