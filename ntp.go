package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func updateRTT(allLines []*csvLine) error {
	var err error
	var response *ntp.Response
	var ipMap = map[string]int64{} // empty map to cache rtt results for reuse

	// scan through all the dns generated csvLines and get rtt for each remote once only
	for _, aLine := range allLines {

		// Need to check if this ip has already been done
		if rtt, ok := ipMap[aLine.ip.String()]; ok {
			aLine.rtt = rtt
			if debug {
				fmt.Fprintf(os.Stderr, "Debug: Read rtt from cache for %s\n", aLine.ip.String())
			}
			continue
		}

		response, err = queryNTP(aLine.ip.String())
		if err != nil {
			aLine.rtt = -1
			if !quiet {
				fmt.Fprintf(os.Stderr, "Error: Unable to get rtt for %s err: %v\n", aLine.ip.String(), err)
			} else {
				retErr = 64
			}
			continue
		}
		aLine.rtt = response.RTT.Milliseconds()
		// store the new value in the cache map
		ipMap[aLine.ip.String()] = aLine.rtt
		if debug {
			fmt.Fprintf(os.Stderr, "Debug: Stored rtt to cache for %s\n", aLine.ip.String())
		}
	}
	return nil
}

func queryNTP(ip string) (*ntp.Response, error) {

	var response *ntp.Response
	var err error

	for i := 0; i < 3; i++ {

		response, err = ntp.Query(ip)

		if err != nil {
			if !quiet {
				fmt.Fprintf(os.Stderr, "Warning: error with ntp request to %s. err: %v\n", ip, err)
			} else {
				retErr = 65
			}
			// if there was an error the wait before trying again
			// sleep 1, 2, 3 seconds plus the ntp.Query 5 second timeout
			time.Sleep(time.Duration(i+1) * time.Second)
		} else {
			// all good so carry on
			return response, nil
		}
	}
	return nil, err
}
