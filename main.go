package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ak1ra24/multiping/utils"
	"github.com/sparrc/go-ping"
)

type PingOpt struct {
	Timeout    time.Duration
	Interval   time.Duration
	Count      int
	Privileged bool
	Pinger     *ping.Pinger
}

type Stats struct {
	Host *net.IPAddr
	Loss float64
	Rtt  time.Duration
	Avg  time.Duration
	Snt  int
}

func NewPingOpt(timeout, interval time.Duration, count int) *PingOpt {

	var privileged bool
	ostype := utils.DiscriminationOS()
	if ostype == "linux" {
		privileged = true
	} else if ostype == "windows" {
		privileged = false
	}

	return &PingOpt{
		Timeout:    timeout,
		Interval:   interval,
		Count:      count,
		Privileged: privileged,
	}

}

func main() {

	d := utils.ReadYaml()

	timeout := time.Second * 100000
	interval := time.Second
	count := -1

	pingopt := NewPingOpt(timeout, interval, count)

	var wg sync.WaitGroup
	wg.Add(len(d.Addresses))

	for _, pinglist := range d.Addresses {

		HostToPing := pinglist.Address

		go func(HostToPing string) {
			defer wg.Done()

			pingopt.PingCheck(HostToPing)
			// fmt.Println(stats)
		}(HostToPing)

		// s := <-OnRecv
		// fmt.Println("Stats: ", s)
	}

	wg.Wait()
}

func (popt *PingOpt) PingCheck(HostToPing string) {

	pinger, err := ping.NewPinger(HostToPing)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	// listen for ctrl-C signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pinger.Stop()
		}
	}()

	var total time.Duration

	pinger.OnRecv = func(pkt *ping.Packet) {
		loss := float64(pinger.PacketsSent-pinger.PacketsRecv) / float64(pinger.PacketsSent) * 100
		snt := pinger.PacketsSent
		total += pkt.Rtt
		avgrtt := total / time.Duration(snt)

		fmt.Printf("\x1b[32m|Host: %s | LOSS: %v | RTT: %v | AVG: %v | SNT: %v |\x1b[0m\n",
			pkt.IPAddr, loss, pkt.Rtt, avgrtt, pinger.PacketsSent)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	pinger.Count = popt.Count
	pinger.Interval = popt.Interval
	pinger.Timeout = popt.Timeout
	pinger.SetPrivileged(popt.Privileged)

	// fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.Run()
}
