package networking

import (
	"errors"
	"net"
	"strings"
	"testing"
)

// skipIfDNSUnavailable skips the test when err is a DNS resolution failure,
// since these tests depend on outbound network/DNS access that may not be
// available in every environment (sandboxes, offline CI, etc.).
func skipIfDNSUnavailable(t *testing.T, err error) {
	t.Helper()
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		t.Skipf("skipping: DNS appears unavailable in this environment: %v", err)
	}
}

func TestGetDomainFromHostname(t *testing.T) {
	cases := []struct {
		name     string
		hostname string
		want     string
		wantErr  bool
	}{
		{"simple subdomain", "www.example.com", "example.com", false},
		{"deep subdomain", "a.b.c.example.com", "example.com", false},
		{"multi-label public suffix", "foo.bar.example.co.uk", "example.co.uk", false},
		{"already at eTLD+1", "example.co.uk", "example.co.uk", false},
		{"bare public suffix errors", "co.uk", "", true},
		{"bare TLD errors", "com", "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := GetDomainFromHostname(c.hostname)
			if c.wantErr {
				if err == nil {
					t.Fatalf("GetDomainFromHostname(%q) expected error, got nil (result %q)", c.hostname, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetDomainFromHostname(%q) unexpected error: %v", c.hostname, err)
			}
			if got != c.want {
				t.Errorf("GetDomainFromHostname(%q) = %q, want %q", c.hostname, got, c.want)
			}
		})
	}
}

func TestReverseLookupIPKnownHost(t *testing.T) {
	names, err := ReverseLookupIP("8.8.8.8")
	skipIfDNSUnavailable(t, err)
	if err != nil {
		t.Fatalf("ReverseLookupIP(8.8.8.8) unexpected error: %v", err)
	}
	if len(names) == 0 {
		t.Fatal("ReverseLookupIP(8.8.8.8) returned no names")
	}

	found := false
	for _, n := range names {
		if strings.Contains(n, "dns.google") {
			found = true
		}
		if strings.HasSuffix(n, ".") {
			t.Errorf("ReverseLookupIP returned name %q with an unstripped trailing dot", n)
		}
	}
	if !found {
		t.Errorf("ReverseLookupIP(8.8.8.8) = %v, expected a name containing %q", names, "dns.google")
	}
}

func TestReverseLookupIPNoRecord(t *testing.T) {
	// .invalid is reserved by RFC 2606 and is guaranteed to never resolve,
	// so this error path doesn't depend on outbound network access.
	_, err := ReverseLookupIP("definitely-not-a-real-host.invalid")
	if err == nil {
		t.Fatal("ReverseLookupIP() with a reserved-invalid hostname expected an error, got nil")
	}
}

func TestGetDomainFromIPKnownHost(t *testing.T) {
	name, tld, err := GetDomainFromIP("8.8.8.8")
	skipIfDNSUnavailable(t, err)
	if err != nil {
		t.Fatalf("GetDomainFromIP(8.8.8.8) unexpected error: %v", err)
	}
	if !strings.Contains(name, "dns.google") {
		t.Errorf("GetDomainFromIP(8.8.8.8) name = %q, want it to contain %q", name, "dns.google")
	}
	if tld != "dns.google" {
		t.Errorf("GetDomainFromIP(8.8.8.8) tld = %q, want %q", tld, "dns.google")
	}
}

func TestGetDomainFromIPNoRecord(t *testing.T) {
	_, _, err := GetDomainFromIP("999.999.999.999")
	if err == nil {
		t.Fatal("GetDomainFromIP() with a malformed IP expected an error, got nil")
	}
}
