package models

import (
	"strings"

	"github.com/dustin/go-humanize"
)

func FormatRupiah(amount int) string {
	// Menggunakan humanize.Comma untuk memformat angka dengan koma sebagai pemisah ribuan
	formatted := humanize.Comma(int64(amount))
	// Ganti koma dengan titik untuk format Rupiah
	return "Rp " + strings.ReplaceAll(formatted, ",", ".")
}

func FormatAngka(amount int) string {
	// Menggunakan humanize.Comma untuk memformat angka dengan koma sebagai pemisah ribuan
	formatted := humanize.Comma(int64(amount))
	// Ganti koma dengan titik untuk format Rupiah
	return strings.ReplaceAll(formatted, ",", ".")
}
