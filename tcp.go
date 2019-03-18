package main
import (
  "fmt"
  "net"
  "time"
)
func Tcc(metric string) {
  conn, err := net.DialTimeout("tcp", carbon_host + ":" + carbon_port , 1*time.Second)
  if err != nil {
    fmt.Printf("\tCarbon: %v\n",err)
    return
  }
  defer conn.Close()
  for metric != "" {
    fmt.Fprintf(conn, "%v\n", metric)
    metric = ""
    break
  }
}