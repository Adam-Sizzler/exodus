// Package reputation asks proxycheck.io what the internet thinks of your exit
// address: whether it is datacenter or residential space, whether it is a known
// VPN or proxy, and how risky it is considered.
//
// This answers a question the rest of geocheck cannot. The geolocation table
// says which country a service assigns you; the path measurement says how
// directly you reach it. Neither explains why a service that geolocates you
// correctly still refuses to serve you — and the usual reason is right here:
// the address is classed as hosting space, or carries a bad reputation.
//
// The service allows 100 queries a day without registration, and geocheck
// spends one per address per run. A free API key raises that to 1,000.
package reputation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"net/url"
	"strings"
	"time"

	"exodus-node/geocheck/internal/netx"
)

// Info is the parsed reputation of one address.
type Info struct {
	IP netip.Addr

	// Type is the address-space classification: "Hosting", "Residential",
	// "Business", "Education", "Wireless". The most actionable single field,
	// because services that refuse "disallowed ISP" are refusing Hosting.
	Type         string
	ASN          string
	Range        string
	Hostname     string
	Provider     string
	Organisation string

	City    string
	Region  string
	Country string
	Code    string

	// Risk is 0..100, and Confidence is how sure the service is of it.
	Risk       int
	Confidence int

	Proxy       bool
	VPN         bool
	Tor         bool
	Hosting     bool
	Scraper     bool
	Compromised bool
	Anonymous   bool

	FirstSeen string
	LastSeen  string

	// Operator names the VPN service when the address belongs to a known one.
	Operator     string
	OperatorURL  string
	Anonymity    string
	NoLogging    bool
	HasPolicies  bool
	DevicesAddr  int
	DevicesSubnt int

	// Warning carries the service's own caveat when it answered with one.
	Warning string
}

// Flags lists the detections that fired, in the order worth reading.
func (i *Info) Flags() []string {
	var out []string
	for _, f := range []struct {
		on   bool
		name string
	}{
		{i.Tor, "Tor"},
		{i.VPN, "VPN"},
		{i.Proxy, "proxy"},
		{i.Compromised, "compromised"},
		{i.Scraper, "scraper"},
		{i.Hosting, "hosting"},
		{i.Anonymous, "anonymous"},
	} {
		if f.on {
			out = append(out, f.name)
		}
	}
	return out
}

// Clean reports whether nothing was flagged at all.
func (i *Info) Clean() bool { return len(i.Flags()) == 0 }

// Residential reports whether the address is end-user space rather than a
// datacenter, which is what most consumer services want to see.
func (i *Info) Residential() bool {
	switch strings.ToLower(i.Type) {
	case "residential", "wireless", "business", "education":
		return true
	default:
		return false
	}
}

// ErrQuotaExceeded means the daily allowance is used up.
var ErrQuotaExceeded = errors.New("proxycheck.io daily query allowance exhausted")

// wire mirrors the parts of the v3 response we read.
type wire struct {
	Network struct {
		ASN          string `json:"asn"`
		Range        string `json:"range"`
		Hostname     string `json:"hostname"`
		Provider     string `json:"provider"`
		Organisation string `json:"organisation"`
		Type         string `json:"type"`
	} `json:"network"`
	Location struct {
		City    string `json:"city_name"`
		Region  string `json:"region_name"`
		Country string `json:"country_name"`
		Code    string `json:"country_code"`
	} `json:"location"`
	DeviceEstimate struct {
		Address int `json:"address"`
		Subnet  int `json:"subnet"`
	} `json:"device_estimate"`
	Detections struct {
		Proxy       bool   `json:"proxy"`
		VPN         bool   `json:"vpn"`
		Compromised bool   `json:"compromised"`
		Scraper     bool   `json:"scraper"`
		Tor         bool   `json:"tor"`
		Hosting     bool   `json:"hosting"`
		Anonymous   bool   `json:"anonymous"`
		Risk        int    `json:"risk"`
		Confidence  int    `json:"confidence"`
		FirstSeen   string `json:"first_seen"`
		LastSeen    string `json:"last_seen"`
	} `json:"detections"`
	Operator *struct {
		Name      string `json:"name"`
		URL       string `json:"url"`
		Anonymity string `json:"anonymity"`
		Policies  *struct {
			Logging bool `json:"logging"`
		} `json:"policies"`
	} `json:"operator"`
}

// Lookup fetches the reputation of ip. An empty apiKey uses the unregistered
// allowance, which is enough for interactive use.
func Lookup(ctx context.Context, stack *netx.Stack, f netx.Family, ip netip.Addr, apiKey string) (*Info, error) {
	if !ip.IsValid() {
		return nil, errors.New("no address to look up")
	}

	endpoint := "https://proxycheck.io/v3/" + ip.String()
	if apiKey != "" {
		endpoint += "?key=" + url.QueryEscape(apiKey)
	}

	ctx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	resp, err := stack.Do(ctx, f, netx.Request{
		URL:       endpoint,
		UserAgent: "geocheck (+https://github.com/adam-sizzler/exodus)",
		Headers:   map[string]string{"Accept": "application/json"},
	})
	if err != nil {
		return nil, err
	}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(resp.Body, &envelope); err != nil {
		return nil, fmt.Errorf("proxycheck.io: unreadable response: %w", err)
	}

	var status string
	if raw, ok := envelope["status"]; ok {
		_ = json.Unmarshal(raw, &status)
	}
	switch status {
	case "denied":
		return nil, ErrQuotaExceeded
	case "error":
		var msg string
		if raw, ok := envelope["message"]; ok {
			_ = json.Unmarshal(raw, &msg)
		}
		if msg == "" {
			msg = "request rejected"
		}
		return nil, fmt.Errorf("proxycheck.io: %s", msg)
	}

	raw, ok := envelope[ip.String()]
	if !ok {
		return nil, fmt.Errorf("proxycheck.io: no entry for %s", ip)
	}
	var w wire
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, fmt.Errorf("proxycheck.io: %w", err)
	}

	info := &Info{
		IP:           ip,
		Type:         w.Network.Type,
		ASN:          w.Network.ASN,
		Range:        w.Network.Range,
		Hostname:     w.Network.Hostname,
		Provider:     w.Network.Provider,
		Organisation: w.Network.Organisation,
		City:         w.Location.City,
		Region:       w.Location.Region,
		Country:      w.Location.Country,
		Code:         w.Location.Code,
		Risk:         w.Detections.Risk,
		Confidence:   w.Detections.Confidence,
		Proxy:        w.Detections.Proxy,
		VPN:          w.Detections.VPN,
		Tor:          w.Detections.Tor,
		Hosting:      w.Detections.Hosting,
		Scraper:      w.Detections.Scraper,
		Compromised:  w.Detections.Compromised,
		Anonymous:    w.Detections.Anonymous,
		FirstSeen:    shortDate(w.Detections.FirstSeen),
		LastSeen:     shortDate(w.Detections.LastSeen),
		DevicesAddr:  w.DeviceEstimate.Address,
		DevicesSubnt: w.DeviceEstimate.Subnet,
	}
	if w.Operator != nil {
		info.Operator = w.Operator.Name
		info.OperatorURL = w.Operator.URL
		info.Anonymity = w.Operator.Anonymity
		if w.Operator.Policies != nil {
			info.HasPolicies = true
			info.NoLogging = !w.Operator.Policies.Logging
		}
	}
	if status == "warning" {
		var msg string
		if raw, ok := envelope["message"]; ok {
			_ = json.Unmarshal(raw, &msg)
		}
		info.Warning = msg
	}
	return info, nil
}

// shortDate trims an ISO 8601 timestamp to its date, which is all the display
// needs and all that is meaningful for a "first seen" figure.
func shortDate(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}
