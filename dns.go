package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/miekg/dns"
)

func getCSVLines(dnsServer, hostname, location string) ([]*csvLine, error) {

	ts := time.Now().Unix()
	newLines := make([]*csvLine, 0)

	// if using local resolver then
	if dnsServer == "127.0.0.1" {
		if foundIPs, err := net.LookupIP(hostname); err == nil {
			for _, ip := range foundIPs {
				// check if we can test this ip
				if x := ip.To4(); x != nil {
					// found a v4 address
					if !ipv4 {
						continue
					}
				} else {
					if !ipv6 {
						continue
					}
				}

				aLine := &csvLine{
					ts:        ts,
					location:  location,
					dnsServer: dnsServer,
					hostname:  hostname,
					ip:        ip,
				}
				newLines = append(newLines, aLine)
			}
		} else {
			if !quiet {
				fmt.Fprintf(os.Stderr, "Error: DNS using local resolver. err: %v\n", err)
			} else {
				retErr = 1
			}
		}
		return newLines, nil
	}

	// ask remote server for ip addresses

	c := dns.Client{}
	m := dns.Msg{}

	// first lookup AAAA
	// do we want ipv6 addresses
	if ipv6 {

		m.SetQuestion(hostname+".", dns.TypeAAAA)
		r, _, err := c.Exchange(&m, dnsServer+":53")

		if err != nil {
			if !quiet {
				fmt.Fprintf(os.Stderr, "Error: DNS AAAA record lookup: %v\n", err)
			} else {
				retErr = 2
			}

		} else {
			for _, ans := range r.Answer {
				newRecord6 := ans.(*dns.AAAA)
				aLine := &csvLine{
					ts:        ts,
					location:  location,
					dnsServer: dnsServer,
					hostname:  hostname,
					ip:        newRecord6.AAAA,
				}
				newLines = append(newLines, aLine)
			}
		}
	}

	// do we want ipv4 addresses
	if ipv4 {
		// reset and ask for A records
		m.SetQuestion(hostname+".", dns.TypeA)
		r, _, err := c.Exchange(&m, dnsServer+":53")

		if err != nil {
			if !quiet {
				fmt.Fprintf(os.Stderr, "Error: DNS A record lookup: %v\n", err)
			} else {
				retErr = 3
			}
		} else {
			for _, ans := range r.Answer {
				newRecord4 := ans.(*dns.A)
				aLine := &csvLine{
					ts:        ts,
					location:  location,
					dnsServer: dnsServer,
					hostname:  hostname,
					ip:        newRecord4.A,
				}
				newLines = append(newLines, aLine)
			}
		}
	}

	if debug {
		fmt.Fprintf(os.Stderr, "Debug: returning %v csvLines\n", len(newLines))
	}
	return newLines, nil
}
