/*
 * Just junk atm Idea is to count qps and report to stdout and carbon/influx and elasticsearch
 * not sure how I want this to work so this is a junk function atm.
 */
package main
import (
  "fmt"
  "time"
)

func qps(time time.Duration, now int64) (time.Duration, int64) {
    /*
     * Collector <- count every seconds
     * time.
     */
    fmt.Printf("%d %v\n", time, now)
  return 10, 100
}
