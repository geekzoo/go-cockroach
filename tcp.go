package main
import (
  "fmt"
  "net"
  "time"
)
func Tcc(bar string) {
  conn, err := net.DialTimeout("tcp", carbon_host + ":" + carbon_port , 1*time.Second)
  if err != nil {
    fmt.Printf("\tCarbon: %v\n",err)
    return
  }
  defer conn.Close()
  for bar != "" {
    fmt.Fprintf(conn, "%v\n", bar)
    bar = ""
    break
  }
  
}