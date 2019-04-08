package main
import (
  "fmt"
  "net"
  "time"
)

func Tcc(metric string) {
    
        t_conn, err := net.DialTimeout("tcp", carbon_host + ":" + carbon_port , 60*time.Second) //Spawns a new connection each call //Need to find a clean way of doing this
        
go func() {
  if err != nil {
        blow_out = false
    return
  }
  defer t_conn.Close()
  for metric != "" {
    fmt.Fprintf(t_conn, "%v\n", metric)
    metric = ""
    break
  }
  t_conn.Close()
}()  
}
