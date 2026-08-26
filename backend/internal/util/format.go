// Package util contains small, dependency-free helpers shared across the
// backend. FormatBytes and FormatBitrate are the single source of truth for
// human-readable byte/bitrate formatting - Telegram notifications, HTTP
// responses, and subscription placeholder rendering all call into this
// package instead of keeping their own copies.
package util

import (
	"fmt"
	"math/big"
)

var byteUnits = [...]string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

// FormatBytes converts a byte count into a human-readable IEC string
// (e.g. "372.53 GiB"), auto-scaling the unit up to EiB. Negative values are
// formatted with a leading "-" against their absolute value.
func FormatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	sign := ""
	value := float64(bytes)
	if value < 0 {
		sign = "-"
		value = -value
	}

	const step = 1024.0
	unit := 0
	for value >= step && unit < len(byteUnits)-1 {
		value /= step
		unit++
	}

	if unit == 0 {
		return fmt.Sprintf("%s%d %s", sign, int64(value), byteUnits[unit])
	}
	return fmt.Sprintf("%s%.2f %s", sign, value, byteUnits[unit])
}

// FormatBigBytes is the *big.Int counterpart of FormatBytes, for aggregate
// values (e.g. panel-wide traffic totals) that may exceed the int64 range.
func FormatBigBytes(value *big.Int) string {
	if value == nil || value.Sign() == 0 {
		return "0 B"
	}

	sign := ""
	abs := new(big.Int).Set(value)
	if abs.Sign() < 0 {
		sign = "-"
		abs.Abs(abs)
	}

	floatValue, _ := new(big.Float).SetInt(abs).Float64()

	const step = 1024.0
	unit := 0
	for floatValue >= step && unit < len(byteUnits)-1 {
		floatValue /= step
		unit++
	}

	return fmt.Sprintf("%s%.2f %s", sign, floatValue, byteUnits[unit])
}

var bitrateUnits = [...]string{"bps", "Kbps", "Mbps", "Gbps", "Tbps"}

// FormatBitrate converts a bits-per-second value into a human-readable
// string (e.g. "1.35 Gbps"), auto-scaling the unit up to Tbps. Negative
// values are formatted with a leading "-" against their absolute value.
func FormatBitrate(bitsPerSecond float64) string {
	if bitsPerSecond == 0 {
		return "0 bps"
	}

	sign := ""
	value := bitsPerSecond
	if value < 0 {
		sign = "-"
		value = -value
	}

	const step = 1000.0
	unit := 0
	for value >= step && unit < len(bitrateUnits)-1 {
		value /= step
		unit++
	}

	if unit == 0 {
		return fmt.Sprintf("%s%.0f %s", sign, value, bitrateUnits[unit])
	}
	return fmt.Sprintf("%s%.2f %s", sign, value, bitrateUnits[unit])
}
