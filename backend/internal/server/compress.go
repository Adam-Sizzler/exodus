package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

const CompressionThresholdBytes = 1024 // 1 KB

// CompressPayloadIfLarge compresses payload with GZIP if size exceeds CompressionThresholdBytes (1 KB).
// Returns original payload if size <= 1 KB.
func CompressPayloadIfLarge(payload []byte) ([]byte, bool, error) {
	if len(payload) <= CompressionThresholdBytes {
		return payload, false, nil
	}

	var buf bytes.Buffer
	gw, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return nil, false, fmt.Errorf("create gzip writer: %w", err)
	}

	if _, err := gw.Write(payload); err != nil {
		_ = gw.Close()
		return nil, false, fmt.Errorf("gzip write: %w", err)
	}

	if err := gw.Close(); err != nil {
		return nil, false, fmt.Errorf("gzip close: %w", err)
	}

	compressed := buf.Bytes()
	// Only return compressed if it actually saved space
	if len(compressed) >= len(payload) {
		return payload, false, nil
	}

	return compressed, true, nil
}

// DecompressPayloadIfNeeded decompresses GZIP payload if compressed flag is true.
func DecompressPayloadIfNeeded(payload []byte, compressed bool) ([]byte, error) {
	if !compressed || len(payload) == 0 {
		return payload, nil
	}

	gr, err := gzip.NewReader(bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create gzip reader: %w", err)
	}
	defer gr.Close()

	decompressed, err := io.ReadAll(gr)
	if err != nil {
		return nil, fmt.Errorf("read decompressed payload: %w", err)
	}

	return decompressed, nil
}
