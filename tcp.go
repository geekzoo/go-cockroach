package main
import (
  "fmt"
  "net"
  "time"
)

func Tcc(metric string) {
    
        if conn == nil {
        conn, err := net.DialTimeout("tcp", carbon_host + ":" + carbon_port , 60*time.Second) //Spawns a new connection each call //Need to find a clean way of doing this
        }
        
go func() {
  if err != nil {
 //       fmt.Printf("\tCarbon: %v\n",err)
        blow_out = false
 //   time.Sleep(60*time.Second)
    return
  }
  defer conn.Close()
  for metric != "" {
    fmt.Fprintf(conn, "%v\n", metric)
    metric = ""
    break
  }
  //conn.Close()
}()  
}
