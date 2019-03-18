package main
import (
  "fmt"
  "time"
)

func qps(time time.Duration, now int64) (time.Duration, int64) {
    fmt.Printf("%d %v\n", time, now)
  return 10, 100
}