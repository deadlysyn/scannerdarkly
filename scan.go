package main

import (
	"fmt"
	"net"
)

func scan() {
	for _, recs := range DB {
		for _, rec := range recs {
			scanTCP(rec)
		}
	}
}

func scanTCP(rec dnsRecord) {
	for _, v := range rec.Values {
		for _, p := range PORTS {
			host := fmt.Sprintf("%v:%v", v, p)
			fmt.Printf("Scanning %v...\n", host)
			conn, err := net.DialTimeout("tcp", host, TIMEOUT)
			if err != nil {
				continue
			}
			defer conn.Close()
			rec.Active = append(rec.Active, host)
		}
	}
}
