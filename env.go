package main

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// PORTS ports to check on DNS targets (default: 80 8080 443 8443)
var PORTS []string

// TIMEOUT net.Dial timeout for network tests (default 10 seconds)
var TIMEOUT time.Duration

// ZONES Route53 zone IDs to audit (default: all in account)
var ZONES []string

func parseEnv() {
	ports := strings.TrimSpace(os.Getenv("PORTS"))
	if len(ports) > 0 {
		for _, p := range strings.Split(ports, " ") {
			PORTS = append(PORTS, p)
		}
	} else {
		PORTS = []string{
			"80",
			"443",
		}
	}

	timeout := strings.TrimSpace(os.Getenv(("TIMEOUT")))
	if len(timeout) > 0 {
		t, _ := strconv.Atoi(timeout)
		TIMEOUT = time.Duration(t) * time.Second
	} else {
		TIMEOUT = 5 * time.Second
	}

	zones := strings.TrimSpace(os.Getenv("ZONES"))
	if len(zones) > 0 {
		for _, z := range strings.Split(zones, " ") {
			ZONES = append(ZONES, z)
		}
	}
}
