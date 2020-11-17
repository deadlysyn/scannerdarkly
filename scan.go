package main

import (
	"fmt"
	"net"
)

func scan() {
	for id, recs := range DB {
		for idx := range recs {
			scanTCP(&DB[id][idx])
		}
	}
}

func scanTCP(rec *dnsRecord) {
	for _, v := range rec.Values {
		for _, p := range PORTS {
			var host string
			if rec.Type == "AAAA" && !rec.Alias {
				host = fmt.Sprintf("[%v]:%v", v, p)
			} else {
				host = fmt.Sprintf("%v:%v", v, p)
			}
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
