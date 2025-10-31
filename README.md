# pool-dns
NTP pool dns tool

This tool allows you to monitor the ip addresses given out by the NTP pool DNS system and/or those returned by the large anycast DNS providers such as Google, Cloudflare and Quad9. This program can be run on a regular basis and the output stored in a file to give you a picture of the ip addresses returned by the pool DNS system over time.

Example usage
```
./pool-dns-darwin-arm64 -h
Usage of ./pool-dns-darwin-arm64:
  -4	Check ipv4 hosts.
  -6	Check ipv6 hosts.
  -debug
    	Display debug info.
  -dns-servers string
    	csv of dns servers to use for lookups.
  -header
    	Display CSV header.
  -hostnames string
    	csv of hostnames to lookup.
  -location string
    	Location string to include in csv.
  -quiet
    	Suppress warning messages.
```

The CSV field headers in the output are:
```
/path/to/app/pool-dns -4 -6 -dns-servers 1,1.1.1.1,8.8.8.8,9.9.9.9 -hostnames 2.pool.ntp.org,2.asia.pool.ntp.org -location sg-sin -quiet -header
time_stamp,location_of_request,dns_server_used,hostname_asked_for,ip_address_returned,round_trip_time_ms
1761940572,sg-sin,1.1.1.1,2.pool.ntp.org,15.235.181.37,0
1761940572,sg-sin-o1,8.8.8.8,2.asia.pool.ntp.org,2402:1f00:8300:821::123,69
...
```

