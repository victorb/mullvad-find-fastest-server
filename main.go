package main

import (
	"fmt"
	"sort"
	"time"

	ping "github.com/sparrc/go-ping"
)

// Get the list by going to https://mullvad.net/en/servers/
// and run the following snippet in the JS console:
// `var arr = []; document.querySelectorAll('div.server-section:nth-child(7) > div:nth-child(3) > table:nth-child(1) > tbody:nth-child(1) > tr > td:nth-child(1)').forEach(a => arr.push(a.innerText)); copy(arr)`
// You'll know have a JS array in the clipboard that you can paste below and
// reformat to Go format
var servers = []string{
	"au1-wireguard",
	"at1-wireguard",
	"be1-wireguard",
	"bg1-wireguard",
	"ca3-wireguard",
	"ca1-wireguard",
	"ca2-wireguard",
	"cz1-wireguard",
	"dk1-wireguard",
	"fi1-wireguard",
	"fr1-wireguard",
	"de3-wireguard",
	"de1-wireguard",
	"de2-wireguard",
	"hk1-wireguard",
	"it1-wireguard",
	"jp1-wireguard",
	"nl1-wireguard",
	"nl2-wireguard",
	"nl3-wireguard",
	"no1-wireguard",
	"pl1-wireguard",
	"ro1-wireguard",
	"sg1-wireguard",
	"sk1-wireguard",
	"es1-wireguard",
	"se3-wireguard",
	"se5-wireguard",
	"se4-wireguard",
	"se2-wireguard",
	"se6-wireguard",
	"ch1-wireguard",
	"ch2-wireguard",
	"gb1-wireguard",
	"gb2-wireguard",
	"gb3-wireguard",
	"us6-wireguard",
	"us4-wireguard",
	"us7-wireguard",
	"us11-wireguard",
	"us12-wireguard",
	"us2-wireguard",
	"us3-wireguard",
	"us1-wireguard",
	"us13-wireguard",
	"us8-wireguard",
	"us10-wireguard",
	"us9-wireguard",
	"us5-wireguard",
}

var mullvadAddr = ".mullvad.net"

var results = map[string]int64{}

func main() {
	for _, server := range servers {
		currentServer := server + mullvadAddr
		pinger, err := ping.NewPinger(currentServer)
		if err != nil {
			panic(err)
		}
		pinger.Count = 10
		pinger.OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}
		pinger.Run()                 // blocks until finished
		stats := pinger.Statistics() // get send/receive/rtt stats
		fmt.Printf("%s = %s\n", currentServer, stats.AvgRtt.String())
		results[currentServer] = stats.AvgRtt.Nanoseconds()
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
