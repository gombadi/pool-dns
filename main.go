package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	debug  bool // Display debug output
	ipv4   bool // Support ipv4 requests
	ipv6   bool // Support ipv6 requests
	quiet  bool // Suppress warning messages
	retErr int  // Any error to return to op system
)

type csvLine struct {
	ts        int64
	location  string
	dnsServer string
	hostname  string
	ip        net.IP
	rtt       int64
}

func main() {

	var (
		servers   string     // csv value of dns servers to use for requests. ip addresses
		hostnames string     // csv value of hostnames to lookup. 2.pool.ntp.org,2.asia.pool.ntp.org etc
		location  string     // Location string to include in the csv output
		header    bool       // Do you want the CSV headers printed
		allLines  []*csvLine // stores the csvlines returned from dns
	)

	flag.StringVar(&servers, "dns-servers", "", "csv of dns servers to use for lookups.")
	flag.StringVar(&hostnames, "hostnames", "", "csv of hostnames to lookup.")
	flag.StringVar(&location, "location", "", "Location string to include in csv.")
	flag.BoolVar(&ipv4, "4", false, "Check ipv4 hosts.")
	flag.BoolVar(&ipv6, "6", false, "Check ipv6 hosts.")
	flag.BoolVar(&quiet, "quiet", false, "Suppress warning messages.")
	flag.BoolVar(&debug, "debug", false, "Display debug info.")
	flag.BoolVar(&header, "header", false, "Display CSV header.")
	flag.Parse()

	serverList := strings.Split(servers, ",")
	if len(serverList) == 0 {
		log.Fatal("Error: No DNS servers supplied")
	}

	hostList := strings.Split(hostnames, ",")
	if len(hostList) == 0 {
		log.Fatal("Error: No hostnames to lookup supplied")
	}

	// loop across each server/hostname combination
	for _, server := range serverList {
		for _, hostname := range hostList {
			if debug {
				fmt.Printf("Debug: getIPs with %s, %s\n", server, hostname)
			}
			someLines, err := getCSVLines(server, hostname, location)
			if err != nil {
				fmt.Printf("Error: getting someIPs. server: %s, hostname: %s err: %v\n", server, hostname, err)
			} else {
				allLines = append(allLines, someLines...)
			}
		}
	}

	// Now we have all the ip address information we can compute the ntp rtt
	updateRTT(allLines)

	if header {
		fmt.Println("time_stamp,location_of_request,dns_server_used,hostname_asked_for,ip_address_returned,round_trip_time_ms")
	}

	// Last task is to output the results in csv format
	for _, v := range allLines {
		fmt.Printf("%d,%s,%s,%s,%s,%d\n", v.ts, v.location, v.dnsServer, v.hostname, v.ip.String(), v.rtt)
	}

	os.Exit(retErr)
}
