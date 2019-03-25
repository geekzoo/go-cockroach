package main
import (
  "fmt"
  "time"
)
func thread(lat time.Duration, average float64) (time.Duration, float64) {
  fmt.Println(lat)
  return lat, 0.0
}