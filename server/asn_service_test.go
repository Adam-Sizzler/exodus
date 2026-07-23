package server

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/Adam-Sizzler/lmdb-go/lmdb"
	"github.com/vmihailenco/msgpack/v5"
)

func TestAsnLmdbService(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "asn-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "asn-prefixes.lmdb")

	// Prepare a dummy LMDB database with AS13335 (Cloudflare) and AS15169 (Google)
	env, err := lmdb.NewEnv()
	if err != nil {
		t.Fatalf("Failed to create env: %v", err)
	}
	defer env.Close()

	if err := env.SetMapSize(10 * 1024 * 1024); err != nil {
		t.Fatalf("Failed to set map size: %v", err)
	}
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		t.Fatalf("Failed to mkdir dbPath: %v", err)
	}
	if err := env.Open(dbPath, 0, 0644); err != nil {
		t.Fatalf("Failed to open env: %v", err)
	}

	err = env.Update(func(txn *lmdb.Txn) error {
		dbi, openErr := txn.OpenRoot(0)
		if openErr != nil {
			return openErr
		}

		// AS13335 data
		cfData, _ := msgpack.Marshal(AsnPrefixes{
			IPv4: []string{"1.1.1.0/24", "1.0.0.0/24"},
			IPv6: []string{"2606:4700::/32"},
		})
		cfKey := testEncodeOrderedBinaryNumberKey(13335)
		if putErr := txn.Put(dbi, cfKey, cfData, 0); putErr != nil {
			return putErr
		}

		// AS15169 data
		gData, _ := msgpack.Marshal(AsnPrefixes{
			IPv4: []string{"8.8.8.0/24"},
			IPv6: []string{"2001:4860::/32"},
		})
		gKey := testEncodeOrderedBinaryNumberKey(15169)
		if putErr := txn.Put(dbi, gKey, gData, 0); putErr != nil {
			return putErr
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Failed to write test LMDB data: %v", err)
	}
	env.Close()

	t.Setenv("ASN_LMDB_PATH", dbPath)
	service := NewAsnLmdbService(nil)
	defer service.Close()

	if !service.IsAvailable() {
		t.Fatalf("Expected AsnLmdbService to be available")
	}

	// Test ResolvePrefixes for Cloudflare AS13335
	v4, v6 := service.ResolvePrefixes(13335)
	if len(v4) != 2 || len(v6) != 1 {
		t.Errorf("Unexpected prefixes for AS13335: v4=%v, v6=%v", v4, v6)
	}

	// Test ResolveASNs for AS13335 + AS15169
	allV4, allV6 := service.ResolveASNs([]int{13335, 15169})
	if len(allV4) != 3 || len(allV6) != 2 {
		t.Errorf("Unexpected combined prefixes for ASNs [13335, 15169]: v4=%v, v6=%v", allV4, allV6)
	}

	// Test non-existent ASN
	noV4, noV6 := service.ResolvePrefixes(999999)
	if len(noV4) != 0 || len(noV6) != 0 {
		t.Errorf("Expected empty prefixes for non-existent ASN, got v4=%v, v6=%v", noV4, noV6)
	}
}

func testEncodeOrderedBinaryNumberKey(n uint32) []byte {
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
