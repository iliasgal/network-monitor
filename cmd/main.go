package main

import (
	"fmt"
	"time"

	"github.com/iliasgal/network-monitor/pkg/ping"
)

func main() {

	host := "google.com"
	count := 4
	for {
		err := ping.PingHost(host, count)
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(5 * time.Second)
	}
}
