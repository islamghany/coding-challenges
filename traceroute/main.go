package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/ipv4"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: cctraceroute <hostname>")
		return
	}

	hostname := os.Args[1]

	ips, err := net.LookupIP(hostname)
	if err != nil || len(ips) == 0 {
		fmt.Printf("Failed to resolve hostname: %s\n", hostname)
		return
	}

	maxHops := 64
	packetSize := 32
	packetContent := "islamghany traceroute pack route"
	destIp := ips[0]
	fmt.Printf("traceroute to %s (%s), %d hops max, %d byte packets\n", hostname, destIp, maxHops, packetSize)

	// Create a UDP socket for sending packets
	udpConn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   destIp,
		Port: 33434, // Traceroute uses ports in the range 33434–33534
	})

	if err != nil {
		fmt.Println("Failed to create UDP connection")
		return
	}
	defer udpConn.Close()

	// Create a raw socket for receiving ICMP messages
	icmpConn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println("Failed to create ICMP connection", err)
		return
	}
	defer icmpConn.Close()

	// loop for each hop
	for ttl := 1; ttl <= maxHops; ttl++ {

		// Set the TTL field of the UDP packet
		p := ipv4.NewPacketConn(udpConn)
		if err := p.SetTTL(ttl); err != nil {
			fmt.Println("Failed to set TTL")
			return
		}

		start := time.Now()
		// send a test UDP packet to the destination
		_, err = udpConn.Write([]byte(packetContent))
		if err != nil {
			fmt.Println("Failed to send test packet")
			return
		}

		// Set a read deadline on the ICMP connection
		icmpConn.SetDeadline(time.Now().Add(2 * time.Second))

		// Read the ICMP response from the first hop
		buff := make([]byte, 1500) // Buffer to hold the ICMP response
		n, add, err := icmpConn.ReadFrom(buff)
		if err != nil {

			// fmt.Println("Failed to read from ICMP connection", err)
			fmt.Printf("%d * * *\n", ttl)
			continue
		}
		rtt := time.Since(start).Milliseconds()
		// Resolve the hostname of the first hop
		hopIp := add.String()
		hostnames, err := net.LookupAddr(hopIp)
		if err != nil || len(hostnames) == 0 {
			fmt.Printf("%d %s (%s) (%d bytes) %d\n", ttl, hopIp, hopIp, n, rtt)
		} else {
			fmt.Printf("%d %s (%s) (%d bytes) %d\n", ttl, hostnames[0], hopIp, n, rtt)
		}

		// Check if the destination is reached
		currentHopIp := net.ParseIP(hopIp)

		if currentHopIp.Equal(ips[0]) {
			fmt.Println("Destination reached")
			return
		}

	}
}
