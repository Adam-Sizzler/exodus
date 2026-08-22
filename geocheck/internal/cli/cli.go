// Package cli wires flags, the network stack and the checks together.
package cli

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/netip"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"exodus-node/geocheck/internal/access"
	"exodus-node/geocheck/internal/ai"
	"exodus-node/geocheck/internal/asn"
	"exodus-node/geocheck/internal/detect"
	"exodus-node/geocheck/internal/geo"
	"exodus-node/geocheck/internal/mtr"
	"exodus-node/geocheck/internal/netx"
	"exodus-node/geocheck/internal/portal"
	"exodus-node/geocheck/internal/render"
	"exodus-node/geocheck/internal/reputation"
	"exodus-node/geocheck/internal/version"
)

type options struct {
	ipv4Only bool
	ipv6Only bool
	iface    string
	proxy    string
	doh      string
	timeout  int
	jsonOut  bool
	group    string
	noGeo    bool
	noPath   bool
	noDetect bool
	noPortal bool
	noAccess bool
	noAI     bool
	noRep    bool
	repKey   string
	portals  string
	targets  string
	detail   bool
	rounds   int
	maxTTL   int
	noRDNS   bool
	mask     bool
	quiet    bool
	showVer  bool
	demo     bool
	svgOut   string
	svgB64   bool
	svgURI   bool
}

// Run executes the tool and returns a process exit code.
func Run(args []string) int {
	opts, err := parse(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			usage()
			return 0
		}
		fmt.Fprintln(os.Stderr, "geocheck:", err)
		return 2
	}
	if opts.showVer {
		fmt.Println("geocheck", version.String())
		return 0
	}

	if opts.demo {
		if err := runDemo(opts); err != nil {
			fmt.Fprintln(os.Stderr, "geocheck:", err)
			return 1
		}
		return 0
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, opts); err != nil {
		if ctx.Err() != nil {
			fmt.Fprintln(os.Stderr, "\ngeocheck: interrupted")
			return 130
		}
		fmt.Fprintln(os.Stderr, "geocheck:", err)
		return 1
	}
	return 0
}

// RunJSONBytes executes geocheck in-process with the provided CLI args and returns the JSON output.
func RunJSONBytes(ctx context.Context, args []string) ([]byte, error) {
	opts, err := parse(args)
	if err != nil {
		return nil, err
	}
	var byteBuf bytes.Buffer
	if err := runToWriter(ctx, opts, &byteBuf); err != nil {
		return nil, err
	}
	return byteBuf.Bytes(), nil
}

func parse(args []string) (*options, error) {
	o := &options{}
	fs := flag.NewFlagSet("geocheck", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	bind := func(p *bool, v bool, usage string, names ...string) {
		for _, n := range names {
			fs.BoolVar(p, n, v, usage)
		}
	}
	bindStr := func(p *string, v string, usage string, names ...string) {
		for _, n := range names {
			fs.StringVar(p, n, v, usage)
		}
	}
	bindInt := func(p *int, v int, usage string, names ...string) {
		for _, n := range names {
			fs.IntVar(p, n, v, usage)
		}
	}

	bind(&o.ipv4Only, false, "test IPv4 only", "4", "ipv4")
	bind(&o.ipv6Only, false, "test IPv6 only", "6", "ipv6")
	bindStr(&o.iface, "", "bind all traffic to an interface name or a local source address", "i", "interface")
	bindStr(&o.proxy, "", "route checks through a SOCKS5 proxy (host:port)", "p", "proxy")
	bindStr(&o.doh, "auto", "DNS-over-HTTPS resolver: auto, off, or an https:// URL", "doh")
	bindInt(&o.timeout, 8, "per-request timeout in seconds", "t", "timeout")
	bind(&o.jsonOut, false, "emit JSON instead of a report", "j", "json")
	bindStr(&o.group, "all", "which geo groups to run: all, services, geoip, cdn", "g", "group")
	bind(&o.noGeo, false, "skip the geolocation checks", "no-geo")
	bind(&o.noPath, false, "skip the connectivity (MTR) checks", "no-mtr")
	bind(&o.noDetect, false, "skip the tunnel and DNS interception checks", "no-detect")
	bind(&o.noPortal, false, "skip the captive-portal connectivity checks", "no-portal")
	bind(&o.noAccess, false, "skip the Stash service-availability checks", "no-access")
	bind(&o.noAI, false, "skip the AI endpoint reachability checks", "no-ai")
	bind(&o.noRep, false, "skip the proxycheck.io address reputation lookup", "no-reputation")
	bindStr(&o.repKey, os.Getenv("PROXYCHECK_API_KEY"), "proxycheck.io API key (raises the daily allowance to 1000)", "proxycheck-key")
	bindStr(&o.portals, "default", "connectivity-check set: a tag, an id, or 'all'", "portal")
	bindStr(&o.targets, "default", "MTR target set: a tag, an id, 'all', or a comma-separated list", "T", "targets")
	bind(&o.detail, false, "print the full per-hop table for every target", "d", "detail")
	bindInt(&o.rounds, 5, "probes per hop", "rounds")
	bindInt(&o.maxTTL, 30, "maximum TTL to probe", "max-ttl")
	bind(&o.noRDNS, false, "skip reverse DNS for hops", "no-rdns")
	bind(&o.mask, false, "mask the public address in the output", "mask")
	bind(&o.quiet, false, "suppress progress output", "q", "quiet")
	bind(&o.showVer, false, "print the version and exit", "V", "version")
	bind(&o.demo, false, "render a sample report from invented data", "demo")
	bindStr(&o.svgOut, "", "write the report as an SVG to a file, or - for stdout", "svg")
	bind(&o.svgB64, false, "write the report as a base64-encoded SVG to stdout", "svg-base64")
	bind(&o.svgURI, false, "write the report as a data: URI, ready to paste into an <img>", "svg-data-uri")

	fs.Usage = usage
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if o.ipv4Only && o.ipv6Only {
		return nil, errors.New("-4 and -6 are mutually exclusive")
	}
	if o.timeout <= 0 {
		return nil, errors.New("timeout must be positive")
	}
	if o.jsonOut && o.svgOut == "-" {
		return nil, errors.New(
			"--json and --svg - would both write to stdout, and two documents in one stream cannot be parsed.\n" +
				"  Use --svg-base64 to carry the picture inside the JSON, or give --svg a file path")
	}
	return o, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `geocheck %s

Determines where the internet thinks you are, and how directly you reach the
networks that matter.

Usage:
  geocheck [options]

Options:
  -4, --ipv4              test IPv4 only
  -6, --ipv6              test IPv6 only
  -i, --interface IF|IP   bind all traffic to an interface name or a local
                          source address (e.g. eth0, or 203.0.113.10)
  -p, --proxy HOST:PORT   route checks through a SOCKS5 proxy
      --doh MODE          DoH resolver: auto (default), off, or an https:// URL
  -t, --timeout SEC       per-request timeout (default 8)
  -g, --group GROUP       geo groups to run: all, services, geoip, cdn
  -T, --targets SET       MTR targets: %s, all, or an id
  -d, --detail            print the full per-hop table for every target
      --rounds N          probes per hop (default 5)
      --max-ttl N         maximum TTL (default 30)
      --no-geo            skip geolocation checks
      --no-mtr            skip connectivity checks
      --no-detect         skip tunnel and DNS interception checks
      --portal SET        connectivity-check set: %s, all, or an id
      --no-portal         skip the captive-portal connectivity checks
      --no-access         skip the Stash service-availability checks
      --no-ai             skip the AI endpoint reachability checks
      --no-reputation     skip the proxycheck.io address reputation lookup
      --proxycheck-key K  proxycheck.io API key ($PROXYCHECK_API_KEY)
      --no-rdns           skip reverse DNS for hops
      --mask              mask the public address in the output
  -j, --json              emit JSON
      --svg FILE          write the report as a self-contained SVG (- for stdout)
      --svg-base64        write that SVG base64-encoded to stdout
      --svg-data-uri      write it as data:image/svg+xml;base64,... ready to paste
  -q, --quiet             suppress progress output
  -V, --version           print the version
  -h, --help              show this help

Examples:
  geocheck                        # everything, both address families
  geocheck -4 -g services         # IPv4, consumer services only
  geocheck -i wg0 -d              # measure through a specific interface, full paths
  geocheck -T all -d              # trace every known target
  geocheck -p 127.0.0.1:1080 -j   # through a SOCKS5 proxy, as JSON

Notes:
  Hop-by-hop tracing needs a raw socket. Without one geocheck still measures
  destination latency over TCP, and says so.
`, version.String(), strings.Join(mtr.Tags(), ", "), strings.Join(portal.Tags(), ", "))
}

// emitSVG writes the report as a picture when a flag asked for one, and reports
// whether it took stdout. Writing to a file does not: stdout stays free for the
// terminal report or the JSON document, so `--json --svg out.svg` yields both.
// Both the demo and a real run go through it, so the two cannot drift apart.
func emitSVG(o *options, report render.Report, findings []detect.Finding) (bool, error) {
	// With --json the picture belongs inside the document, not next to it on
	// the same stream; JSON() embeds it and this leaves stdout alone.
	if (o.svgB64 || o.svgURI) && !o.jsonOut {
		if o.svgURI {
			if _, err := io.WriteString(os.Stdout, render.SVGDataURIPrefix); err != nil {
				return true, err
			}
		}
		if err := render.SVGBase64(os.Stdout, report, findings); err != nil {
			return true, err
		}
		_, err := fmt.Fprintln(os.Stdout)
		return true, err
	}
	switch o.svgOut {
	case "-":
		return true, render.SVG(os.Stdout, report, findings)
	case "":
		return false, nil
	}

	f, err := os.Create(o.svgOut)
	if err != nil {
		return true, err
	}
	if err := render.SVG(f, report, findings); err != nil {
		_ = f.Close()
		_ = os.Remove(o.svgOut)
		return true, err
	}
	if err := f.Close(); err != nil {
		return true, err
	}
	// To stderr, so `--svg -` and a file path behave the same on stdout.
	fmt.Fprintln(os.Stderr, "wrote", o.svgOut)
	// The file took nothing from stdout, so whatever else was asked for still
	// gets written there.
	return false, nil
}

// runDemo prints a report built from invented data. It opens no sockets, so
// it renders identically everywhere and can be recorded for documentation
// without exposing whoever runs it.
func runDemo(o *options) error {
	report := render.DemoReport(version.String())
	report.MaskIP = o.mask
	report.TraceDetail = o.detail

	// The skip flags apply here too, so documentation can show one section at
	// a time instead of a screenshot too tall to read.
	if o.noGeo {
		report.Geo = nil
	}
	if o.noPath {
		report.Trace = nil
	}
	if o.noPortal {
		report.Portal = nil
	}
	if o.noAccess {
		report.Access = nil
	}
	if o.noAI {
		report.AI = nil
	}
	if o.noRep {
		report.Reputation = nil
	}

	report.EmbedSVG = o.jsonOut && (o.svgB64 || o.svgURI)
	if done, err := emitSVG(o, report, nil); done {
		return err
	}
	if o.jsonOut {
		return render.JSON(os.Stdout, report, nil, demoTimestamp)
	}
	render.NewOutput(os.Stdout).Print(report)
	return nil
}

// demoTimestamp is fixed so that regenerating the sample JSON produces no diff
// unless the report itself changed.
var demoTimestamp = time.Date(2026, time.August, 17, 12, 0, 0, 0, time.UTC)

func run(ctx context.Context, o *options) error {
	return runToWriter(ctx, o, os.Stdout)
}

func runToWriter(ctx context.Context, o *options, w io.Writer) error {
	start := time.Now()
	prog := render.NewProgress(os.Stderr, !o.quiet && !o.jsonOut)

	stack, err := netx.New(ctx, netx.Options{
		Interface: o.iface,
		Proxy:     o.proxy,
		Timeout:   time.Duration(o.timeout) * time.Second,
		DoH:       o.doh,
	})
	if err != nil {
		return err
	}
	defer stack.Close()

	asnResolver := asn.New(stack.Resolver())

	// Which address families can actually reach the internet.
	prog.Set("detecting public addresses")
	families, v4, v6 := detectFamilies(ctx, stack, o)
	if len(families) == 0 {
		prog.Stop()
		return errors.New("no working IPv4 or IPv6 connectivity was found")
	}

	env := &geo.Env{Stack: stack, IPv4: v4, IPv6: v6}

	var (
		wg       sync.WaitGroup
		identity geo.Identity
		findings []detect.Finding
		results  []geo.Result
		portals  []portal.Result
		accesses []access.Result
		aiRes    []ai.Result
		rep      *reputation.Info
		repErr   error
		traces   []*mtr.Report
		traceCap mtr.Capability
		traceErr error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		identity = geo.Describe(ctx, stack, asnResolver, v4, v6)
	}()

	if !o.noDetect {
		wg.Add(1)
		go func() {
			defer wg.Done()
			findings = detect.Run(ctx, detect.Options{
				Timeout: 3 * time.Second,
				// Plain port-53 probes would sidestep the proxy and describe
				// the host's own DNS instead of the tunnel's.
				SkipDNS: o.proxy != "",
			})
		}()
	}

	if !o.noRep {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rep, repErr = reputation.Lookup(ctx, stack, families[0], env.PublicIP(families[0]), o.repKey)
		}()
	}

	if !o.noAccess {
		wg.Add(1)
		go func() {
			defer wg.Done()
			accesses = access.Run(ctx, stack, families[0], access.Checks(), 6)
		}()
	}

	if !o.noAI {
		wg.Add(1)
		go func() {
			defer wg.Done()
			aiRes = ai.Run(ctx, stack, families[0], ai.Checks(), 6)
		}()
	}

	if !o.noPortal {
		wg.Add(1)
		go func() {
			defer wg.Done()
			portals = portal.Run(ctx, stack, families[0], selectPortals(o.portals), 8)
		}()
	}

	// Geolocation and path measurement use different resources, so they run
	// together; the run takes as long as the slower of the two rather than both.
	if !o.noGeo {
		wg.Add(1)
		go func() {
			defer wg.Done()
			checks := selectGeoChecks(o.group)
			runner := &geo.Runner{
				Env: env, Families: families, Concurrency: 12,
				OnProgress: func(done, total int, name string) {
					prog.Set(fmt.Sprintf("geo %d/%d  %s", done, total, name))
				},
			}
			results = runner.Run(ctx, checks)
		}()
	}

	if !o.noPath {
		wg.Add(1)
		go func() {
			defer wg.Done()
			traces, traceCap, traceErr = tracePaths(ctx, o, stack, asnResolver, families[0], prog)
		}()
	}

	wg.Wait()
	prog.Stop()

	if traceErr != nil && len(traces) == 0 && !o.noPath {
		fmt.Fprintln(os.Stderr, "geocheck: connectivity checks unavailable:", traceErr)
	}

	report := render.Report{
		Version:       version.String(),
		Identity:      identity,
		Resolver:      stack.Resolver().Active(),
		Interface:     o.iface,
		Proxy:         o.proxy,
		Families:      families,
		Geo:           results,
		Portal:        portals,
		Access:        accesses,
		AI:            aiRes,
		Reputation:    rep,
		ReputationErr: repErr,
		Trace:         traces,
		TraceCap:      traceCap,
		Duration:      time.Since(start),
		MaskIP:        o.mask,
		TraceDetail:   o.detail,
	}

	report.EmbedSVG = o.jsonOut && (o.svgB64 || o.svgURI)
	if done, err := emitSVG(o, report, findings); done {
		return err
	}
	if o.jsonOut {
		return render.JSON(w, report, findings, time.Now())
	}
	out := render.NewOutput(os.Stdout)
	out.PrintFindings(findings)
	out.Print(report)
	return nil
}

// detectFamilies finds the public address of each family the host can use.
func detectFamilies(ctx context.Context, stack *netx.Stack, o *options) ([]netx.Family, netip.Addr, netip.Addr) {
	var (
		families []netx.Family
		v4, v6   netip.Addr
		wg       sync.WaitGroup
		mu       sync.Mutex
	)

	probe := func(f netx.Family) {
		defer wg.Done()
		if !netx.HasFamily(f) {
			return
		}
		addr, ok := geo.PublicIP(ctx, stack, f)
		if !ok {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if f == netx.V6 {
			v6 = addr
		} else {
			v4 = addr
		}
	}

	// A literal source address in --interface fixes the family; probing the
	// other one would measure a different address than the one requested.
	wantV4, wantV6 := !o.ipv6Only, !o.ipv4Only
	if pinned, ok := stack.PinnedFamily(); ok {
		wantV4, wantV6 = pinned == netx.V4, pinned == netx.V6
	}

	if wantV4 {
		wg.Add(1)
		go probe(netx.V4)
	}
	if wantV6 {
		wg.Add(1)
		go probe(netx.V6)
	}
	wg.Wait()

	if v4.IsValid() {
		families = append(families, netx.V4)
	}
	if v6.IsValid() {
		families = append(families, netx.V6)
	}
	return families, v4, v6
}

func selectGeoChecks(group string) []geo.Check {
	switch strings.ToLower(group) {
	case "services", "custom":
		return geo.ServiceChecks()
	case "geoip", "primary":
		return geo.DatabaseChecks()
	case "cdn":
		return geo.CDNChecks()
	default:
		checks := geo.ServiceChecks()
		checks = append(checks, geo.DatabaseChecks()...)
		return append(checks, geo.CDNChecks()...)
	}
}

func tracePaths(
	ctx context.Context, o *options, stack *netx.Stack,
	ar *asn.Resolver, family netx.Family, prog *render.Progress,
) ([]*mtr.Report, mtr.Capability, error) {
	targets := selectTargets(o.targets)
	if len(targets) == 0 {
		return nil, mtr.Capability{}, fmt.Errorf("no targets match %q", o.targets)
	}

	tracer, err := mtr.NewTracer(mtr.Config{
		Family:   family,
		SourceIP: stack.LocalAddr(family),
		// The resolved device name, which may have come from an address.
		Interface:  stack.BindDevice(),
		MaxTTL:     o.maxTTL,
		Rounds:     o.rounds,
		Timeout:    1500 * time.Millisecond,
		Targets:    21,
		Resolver:   stack.Resolver(),
		ASN:        ar,
		ReverseDNS: !o.noRDNS,
	})
	if err != nil {
		return nil, mtr.Capability{}, err
	}
	defer func() { _ = tracer.Close() }()

	capability := tracer.Capability()
	capability.Hint = strings.ReplaceAll(capability.Hint, mtr.BinaryPathPlaceholder, executablePath())

	var done int
	var mu sync.Mutex
	reports := tracer.Run(ctx, targets, func(r *mtr.Report) {
		mu.Lock()
		done++
		n := done
		mu.Unlock()
		prog.Set(fmt.Sprintf("path %d/%d  %s", n, len(targets), r.Target.Name))
	})
	return reports, capability, nil
}

func selectTargets(spec string) []mtr.Target {
	parts := strings.Split(spec, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return mtr.Select(parts...)
}

func executablePath() string {
	if p, err := os.Executable(); err == nil {
		return p
	}
	return "geocheck"
}

// selectPortals resolves the --portal value to a set of endpoints, falling back
// to the default set rather than silently checking nothing.
func selectPortals(spec string) []portal.Endpoint {
	parts := strings.Split(spec, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	if eps := portal.Select(parts...); len(eps) > 0 {
		return eps
	}
	return portal.DefaultEndpoints()
}
