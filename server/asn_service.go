package server

import (
	"archive/tar"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"exodus-node/config"

	"github.com/Adam-Sizzler/lmdb-go/lmdb"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	DefaultAsnLmdbPath    = "/usr/local/share/asn/asn-prefixes.lmdb"
	FallbackAsnLmdbPath   = "/var/lib/exodus-node/asn/asn-prefixes.lmdb"
	DefaultAsnReleaseURL  = "https://github.com/Adam-Sizzler/lmdb-go/releases/download/latest/asn-prefixes-lmdb.tar.gz"
	DefaultAsnJSONRelease = "https://github.com/Adam-Sizzler/lmdb-go/releases/download/latest/asn-prefixes.json"
)

// AsnPrefixes represents the IPv4 and IPv6 subnets associated with an ASN.
type AsnPrefixes struct {
	IPv4 []string `json:"ipv4" msgpack:"ipv4"`
	IPv6 []string `json:"ipv6" msgpack:"ipv6"`
}

// AsnLmdbService provides high-performance lookup of IP prefixes by ASN number from LMDB.
type AsnLmdbService struct {
	mu          sync.RWMutex
	dbPath      string
	env         *lmdb.Env
	isAvailable bool
	logger      interface {
		Log(string, ...any)
		Warn(string, ...any)
		Error(string, ...any)
		Debug(string, ...any)
		Info(string, ...any)
	}
	stopChan chan struct{}
}

// NewAsnLmdbService initializes the ASN LMDB lookup service.
func NewAsnLmdbService(cfg *config.NodeConfig) *AsnLmdbService {
	service := &AsnLmdbService{
		stopChan: make(chan struct{}),
	}
	if cfg != nil && cfg.Logger != nil {
		service.logger = cfg.LoggerFor("AsnLmdbService")
	}

	lmdbPath := os.Getenv("ASN_LMDB_PATH")
	if lmdbPath == "" {
		lmdbPath = DefaultAsnLmdbPath
	}
	service.dbPath = lmdbPath

	if err := service.openDB(); err != nil {
		// Try fallback location if primary path doesn't exist
		if service.dbPath != FallbackAsnLmdbPath {
			service.dbPath = FallbackAsnLmdbPath
			_ = service.openDB()
		}
	}

	if service.isAvailable {
		if service.logger != nil {
			service.logger.Info("[OK] ASN LMDB database opened successfully: " + service.dbPath)
		}
	} else {
		if service.logger != nil {
			service.logger.Info(fmt.Sprintf("ASN LMDB database not found at %s — downloading dataset in background...", service.dbPath))
		}
		go func() {
			if err := service.FetchAndLoadDataset(DefaultAsnReleaseURL); err != nil {
				if service.logger != nil {
					service.logger.Warn(fmt.Sprintf("Failed to download initial ASN dataset: %v — ASN lookup disabled (graceful degradation)", err))
				}
			} else {
				if service.logger != nil {
					service.logger.Info("[OK] ASN LMDB database downloaded and loaded successfully: " + service.dbPath)
				}
			}
		}()
	}

	go service.startPeriodicUpdater()

	return service
}

func (s *AsnLmdbService) startPeriodicUpdater() {
	ticker := time.NewTicker(7 * 24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			if err := s.FetchAndLoadDataset(DefaultAsnReleaseURL); err != nil {
				if s.logger != nil {
					s.logger.Warn(fmt.Sprintf("Periodic ASN dataset update failed: %v", err))
				}
			} else {
				if s.logger != nil {
					s.logger.Info("[OK] Periodic ASN dataset updated successfully")
				}
			}
		}
	}
}

func (s *AsnLmdbService) openDB() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(s.dbPath); os.IsNotExist(err) {
		s.isAvailable = false
		return fmt.Errorf("lmdb path does not exist: %s", s.dbPath)
	}

	env, err := lmdb.NewEnv()
	if err != nil {
		s.isAvailable = false
		return fmt.Errorf("create lmdb env: %w", err)
	}

	_ = env.SetMaxDBs(10)
	if err := env.Open(s.dbPath, lmdb.Readonly, 0644); err != nil {
		env.Close()
		s.isAvailable = false
		return fmt.Errorf("open lmdb env at %s: %w", s.dbPath, err)
	}

	if s.env != nil {
		s.env.Close()
	}
	s.env = env
	s.isAvailable = true
	return nil
}

// Close releases the LMDB environment resources.
func (s *AsnLmdbService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.stopChan:
	default:
		close(s.stopChan)
	}

	if s.env != nil {
		err := s.env.Close()
		s.env = nil
		s.isAvailable = false
		return err
	}
	return nil
}

// IsAvailable returns whether the ASN LMDB database is currently loaded.
func (s *AsnLmdbService) IsAvailable() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isAvailable
}

// GetByASN returns the AsnPrefixes entry for a given ASN number.
func (s *AsnLmdbService) GetByASN(asn int) (*AsnPrefixes, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isAvailable || s.env == nil || asn <= 0 {
		return nil, nil
	}

	var entry AsnPrefixes
	found := false

	err := s.env.View(func(txn *lmdb.Txn) error {
		dbi, openErr := txn.OpenRoot(0)
		if openErr != nil {
			return openErr
		}

		// Keys in LMDB may be encoded as ordered-binary double numbers or raw uint32
		keyCandidates := [][]byte{
			encodeOrderedBinaryNumberKey(uint32(asn)),
			uint32ToBytes(uint32(asn)),
		}

		for _, keyBytes := range keyCandidates {
			valBytes, err := txn.Get(dbi, keyBytes)
			if err == nil && len(valBytes) > 0 {
				if unmarshalErr := msgpack.Unmarshal(valBytes, &entry); unmarshalErr == nil {
					found = true
					return nil
				}
				if unmarshalErr := json.Unmarshal(valBytes, &entry); unmarshalErr == nil {
					found = true
					return nil
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &entry, nil
}

// ResolvePrefixes returns IPv4 and IPv6 prefixes for a given ASN.
func (s *AsnLmdbService) ResolvePrefixes(asn int) (ipv4, ipv6 []string) {
	entry, err := s.GetByASN(asn)
	if err != nil || entry == nil {
		return nil, nil
	}
	return entry.IPv4, entry.IPv6
}

// ResolveASNs resolves a list of ASNs into deduplicated and sorted lists of IPv4 and IPv6 CIDRs.
func (s *AsnLmdbService) ResolveASNs(asns []int) (ipv4, ipv6 []string) {
	if len(asns) == 0 {
		return nil, nil
	}

	seenV4 := make(map[string]struct{})
	seenV6 := make(map[string]struct{})

	for _, asn := range asns {
		v4List, v6List := s.ResolvePrefixes(asn)
		for _, cidr := range v4List {
			seenV4[cidr] = struct{}{}
		}
		for _, cidr := range v6List {
			seenV6[cidr] = struct{}{}
		}
	}

	outV4 := make([]string, 0, len(seenV4))
	for cidr := range seenV4 {
		outV4 = append(outV4, cidr)
	}
	sort.Strings(outV4)

	outV6 := make([]string, 0, len(seenV6))
	for cidr := range seenV6 {
		outV6 = append(outV6, cidr)
	}
	sort.Strings(outV6)

	return outV4, outV6
}

// FetchAndLoadDataset downloads and extracts the daily ASN dataset archive if missing or requested.
func (s *AsnLmdbService) FetchAndLoadDataset(urlStr string) error {
	if urlStr == "" {
		urlStr = DefaultAsnReleaseURL
	}

	targetDir := filepath.Dir(s.dbPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("create asn dir: %w", err)
	}

	if s.logger != nil {
		s.logger.Info("Downloading ASN dataset from " + urlStr + "...")
	}

	client := &http.Client{Timeout: 120 * time.Second}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}
	req.Header.Set("User-Agent", "exodus-node/asn-updater")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http get %s failed: %w", urlStr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http get %s returned status %d", urlStr, resp.StatusCode)
	}

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("read gzip stream: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar entry: %w", err)
		}

		targetPath := filepath.Join(targetDir, filepath.Base(header.Name))
		if header.Typeflag == tar.TypeDir {
			_ = os.MkdirAll(targetPath, 0755)
			continue
		}

		_ = os.MkdirAll(filepath.Dir(targetPath), 0755)
		out, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("create file %s: %w", targetPath, err)
		}

		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return fmt.Errorf("write file %s: %w", targetPath, err)
		}
		out.Close()
	}

	// Re-open LMDB database
	return s.openDB()
}

func encodeOrderedBinaryNumberKey(n uint32) []byte {
	f := float64(n)
	bits := math.Float64bits(f)
	highInt, lowInt := uint32(bits>>32), uint32(bits)

	var length int
	if (lowInt & 0xf) != 0 {
		length = 9
	} else if (lowInt & 0xfffff) != 0 {
		length = 8
	} else if lowInt != 0 || (highInt&0xf) != 0 {
		length = 6
	} else {
		length = 4
	}

	b0_3 := (highInt >> 4) | 0x10000000
	b4_7 := (lowInt >> 4) | (highInt << 28)
	b8 := uint8((lowInt & 0xf) << 4)

	buf := make([]byte, 9)
	binary.BigEndian.PutUint32(buf[0:4], b0_3)
	binary.BigEndian.PutUint32(buf[4:8], b4_7)
	buf[8] = b8

	return buf[:length]
}

func uint32ToBytes(n uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, n)
	return buf
}
