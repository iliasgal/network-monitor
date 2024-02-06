package main

import (
	"fmt"
	"time"

	"github.com/iliasgal/network-monitor/pkg/metrics"
)

func main() {

	host := "google.com"
	count := 4
	for {
		err := metrics.PingHost(host, count)
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(5 * time.Second)
	}
}
