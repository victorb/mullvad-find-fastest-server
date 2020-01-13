package main

import (
	"fmt"
	"sort"
	"time"

	ping "github.com/sparrc/go-ping"
	servers "github.com/victorb/mullvad-find-fastest-server/servers"
)

var pingCount = 3

var mullvadAddr = ".mullvad.net"

var results = map[string]int64{}

func main() {
	for _, server := range servers.GetServers() {
		pinger, err := ping.NewPinger(server)
		if err != nil {
			panic(err)
		}
		pinger.Count = pingCount
		pinger.OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}
		pinger.Run()                 // blocks until finished
		stats := pinger.Statistics() // get send/receive/rtt stats
		fmt.Printf("%s = %s\n", server, stats.AvgRtt.String())
		results[server] = stats.AvgRtt.Nanoseconds()
	}
	type kv struct {
		Key   string
		Value int64
	}

	var ss []kv
	for k, v := range results {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	fmt.Println("## Final Results (least latency first):")
	for _, kv := range ss {
		durr, _ := time.ParseDuration(fmt.Sprintf("%dns", kv.Value))
		fmt.Printf("%s = %s\n", kv.Key, durr.String())
	}
}
