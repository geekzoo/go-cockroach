package main

import (
	"log"
//	"os"
	"fmt"
	"io/ioutil"
)

func s_sys() {
file, err := ioutil.ReadFile("/proc/loadavg") // For read access.
if err != nil {
	log.Fatal(err)
}
fmt.Print(string(file))
}