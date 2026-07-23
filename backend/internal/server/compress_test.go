package server

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompressPayloadIfLarge_SmallPayloadNotCompressed(t *testing.T) {
	small := []byte("hello world small payload")
	compressed, isCompressed, err := CompressPayloadIfLarge(small)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isCompressed {
		t.Fatalf("expected small payload to not be compressed")
	}
	if !bytes.Equal(compressed, small) {
		t.Fatalf("payload modified when not compressed")
	}
}

func TestCompressPayloadIfLarge_LargePayloadCompressedAndDecompressed(t *testing.T) {
	largeText := strings.Repeat("Exodus Subscription Page Test Payload 1234567890 ", 50)
	large := []byte(largeText)

	if len(large) <= CompressionThresholdBytes {
		t.Fatalf("test payload too small: %d", len(large))
	}

	compressed, isCompressed, err := CompressPayloadIfLarge(large)
	if err != nil {
		t.Fatalf("unexpected compression error: %v", err)
	}
	if !isCompressed {
		t.Fatalf("expected large payload to be compressed")
	}
	if len(compressed) >= len(large) {
		t.Fatalf("expected compressed size (%d) to be smaller than original (%d)", len(compressed), len(large))
	}

	decompressed, err := DecompressPayloadIfNeeded(compressed, true)
	if err != nil {
		t.Fatalf("unexpected decompression error: %v", err)
	}

	if !bytes.Equal(decompressed, large) {
		t.Fatalf("decompressed payload does not match original")
	}
}
