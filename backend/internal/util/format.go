package util

import "fmt"

// FormatData formats a numerical traffic or speed value.
func FormatData(value float64, unit string) string {
	const (
		gbit = 1_000_000_000
		mbit = 1_000_000
		kbit = 1_000
		GiB  = 1 << 30
		MiB  = 1 << 20
		KiB  = 1 << 10
	)
	dividers := map[string]float64{
		"bit":  1,
		"kbit": kbit,
		"mbit": mbit,
		"gbit": gbit,
		"Byte": 1,
		"KiB":  KiB,
		"MiB":  MiB,
		"GiB":  GiB,
	}

	if unit == "bps" {
		switch {
		case value >= gbit:
			return fmt.Sprintf("%.2f Gbps", value/gbit)
		case value >= mbit:
			return fmt.Sprintf("%.2f Mbps", value/mbit)
		case value >= kbit:
			return fmt.Sprintf("%.2f Kbps", value/kbit)
		default:
			return fmt.Sprintf("%.0f bps", value)
		}
	} else if unit == "byte" {
		switch {
		case value >= GiB:
			return fmt.Sprintf("%.2f GiB", value/GiB)
		case value >= MiB:
			return fmt.Sprintf("%.2f MiB", value/MiB)
		case value >= KiB:
			return fmt.Sprintf("%.2f KiB", value/KiB)
		default:
			return fmt.Sprintf("%.0f B", value)
		}
	} else if divider, ok := dividers[unit]; ok {
		return fmt.Sprintf("%.2f %s", value/divider, unit)
	} else {
		return fmt.Sprintf("%.0f", value)
	}
}
