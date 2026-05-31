// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2026/03/15 22:54
// Original filename: networking/dns.go

package networking

import (
	"fmt"
	"net"
	"strings"
	"golang.org/x/net/publicsuffix"
)

// Returns a slice of hostnames (the PTR records returned by your system resolver)
// Returns an error if there is no PTR record, or if the resolver returns any kind of lookup error

func ReverseLookupIP(ip string) ([]string, error) {
	ptrRec := []string{}

	names, err := net.LookupAddr(ip)
	if err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("no PTR record found for %s", ip)
	}

	for _, name := range names {
		ptrRec = append(ptrRec, strings.TrimSuffix(name, "."))
	}
	return ptrRec, nil
}

// This one uses the above ReverseLookup to fetch the PTR records associated with the IP address
// Once done, it fetches the domain name in that PTR record

func GetDomainFromIP(ip string) (string, string, error) {
	var names []string
	var err error
	var tld string

	if names, err = ReverseLookupIP(ip); err != nil {
		return "", "", err
	}
	if tld, err = publicsuffix.EffectiveTLDPlusOne(names[0]); err != nil {
		return "", "", err
	}
	return names[0], tld, nil
}

// Similar to GetDomainFromIP, except that it accepts a hostname as a parameter instead of an IP address

func GetDomainFromHostname(hostname string) (string, error) {
	var tld string
	var err error

	if tld, err = publicsuffix.EffectiveTLDPlusOne(hostname); err != nil {
		return "", err
	}
	return tld, nil
}
