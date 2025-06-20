package models

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

type Transaksi struct {
	ID        string    `gorm:"primary_key;column:id_transaksi;autoIncrement"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"` //gorm tag untuk autocreatetime
	UserId    string    `gorm:"column:user_id"`
	User      User      `gorm:"foreignKey:user_id;references:id"` //relasi 1 to many: 1 user menangani beberapa transaksi. Buat foreign key nya, dan buat var User dari struct User serperti ini
	//relasi many to many
	BarangTerjual []Barang `gorm:"many2many:detil_transaksi;foreignKey:id_transaksi;joinForeignKey:id_transaksi;references:id;joinReferences:id_barang"`
	//format: tabel_detil;foreignKey:PK_tabel_ini;joinForeignKey:nama_field_PK_di_tabel_detil;references:PK_tabel_master_lainnya;joinReferences:nama_field_PK_di_tabel_detil
}

func (k *Transaksi) TableName() string {
	return "transaksi" //nama table pada db
}

// Save Data Master Transaksi
func SaveMasterTransaksi(db *gorm.DB, userID string) string {
	transaksi := Transaksi{
		UserId: userID,
	}

	//query: INSERT INTO `products` (`id`,`name`,`price`,`created_at`,`updated_at`) VALUES ('P001','Laptop ASUS',10250000,'2024-12-06 15:15:51.069','2024-12-06 15:15:51.069')
	err := db.Create(&transaksi).Error
	if err != nil {
		panic(err)
	}

	idTransaksi := transaksi.ID //simpan ID transaksi untuk add more detil_jual
	return idTransaksi          //return value untuk detil transaksi
}

// cetak laporan transaksi full
func CetakLaporanTransaksiSummary(db *gorm.DB) {
	var total, totalItem int

	// 	Header
	header := []string{"ID", "Tanggal", "#Item(s)", "Total"}

	// Column widths
	w := []float64{12.0, 28.0, 20.0, 30.0}
	wSum := 0.0
	for _, v := range w {
		wSum += v
	}

	//setting orientation and size
	pdf := gofpdf.New("P", "mm", "A4", "")
	//set font style
	pdf.SetFont("Arial", "B", 16)
	//create a new page
	pdf.AddPage()

	//pdf.Cell(40, 10, "Laporan Stok Minimarket")
	pdf.WriteAligned(0, 15, "Laporan Penjualan Minimarket", "C")
	pdf.Ln(0)

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//set font style for header
	pdf.SetFont("Arial", "B", 12)
	left := (210 - wSum) / 2
	pdf.SetX(left)

	//deklarasi struct dan query
	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int       `gorm:"column:id_transaksi"`
		CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
		Item        int       `gorm:"column:item"`
		Subtotal    int       `gorm:"column:subtotal"`
	}

	var detil []DetilTrans

	err := db.Table("detil_transaksi").Select("transaksi.created_at, detil_transaksi.id_transaksi, sum(jumlah) as item, sum(harga*jumlah) as subtotal").Joins("join barang on detil_transaksi.id_barang=barang.id").Joins("join transaksi on detil_transaksi.id_transaksi=transaksi.id_transaksi").Order("detil_transaksi.id_transaksi asc").Group("detil_transaksi.id_transaksi").Find(&detil).Error
	if err != nil {
		panic(err)
	}

	pdf.SetX(left)
	//print header
	for j, str := range header {
		pdf.CellFormat(w[j], 8, str, "1", 0, "C,M", false, 0, "")
	}
	pdf.Ln(-1)

	//set font style for data (not bold)
	pdf.SetFont("Arial", "", 12)

	// Data
	for _, row := range detil {
		pdf.SetX(left)
		pdf.CellFormat(w[0], 8, strconv.Itoa(row.IdTransaksi), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[1], 8, row.CreatedAt.Format("2006-01-02"), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[2], 8, FormatAngka(row.Item), "1", 0, "R,M", false, 0, "")
		pdf.CellFormat(w[3], 8, FormatAngka(row.Subtotal), "1", 0, "R,M", false, 0, "")

		//hitung total item dan subtotal secara incremental
		total += row.Subtotal
		totalItem++

		pdf.Ln(-1)
	}
	//Summary

	pdf.SetX(left)
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//harus selalu di set agar bisa nyambung dengan row sebelumnya
	pdf.SetX(left)
	pdf.CellFormat(60, 8, "Total Nilai Penjualan: ", "1", 0, "L,M", false, 0, "")
	pdf.CellFormat(30, 8, FormatAngka(total), "1", 0, "R,M", false, 0, "")

	//footer test
	pdf.SetFooterFunc(func() {
		currentTime := time.Now() //var tanggal dan waktu saat ini

		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 11)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Printed on: %v", currentTime.Format("2006-01-02 15:04:05")),
			"", 0, "L", false, 0, "")
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "R", false, 0, "")
	})

	pdf.SetX(left)

	//errs -> untuk cetak report
	errs := pdf.OutputFileAndClose("assets/laporan_transaksi_summary.pdf") //URL silakan disetting
	if errs != nil {
		panic(errs)
	} else {
		fmt.Println("Laporan transaksi telah dibuat! (laporan_transaksi_summary.pdf)")
	}

}

type BarangDetilTransaksi struct {
	IdTransaksi string `gorm:"column:id_transaksi"`
	Id          int    `gorm:"column:id_barang"`
	Jumlah      int    `gorm:"column:jumlah"`
}

var brg []BarangDetilTransaksi //slice untuk memasukkan data

func TambahDetilTransaksi(idTransaksi string, idBarang int, jumlah int) {
	//tambah barang pada slice BarangDetilTransaksi untuk diinput sekaligus
	brg = append(brg, BarangDetilTransaksi{idTransaksi, idBarang, jumlah})
	//pengecekan barang masuk ke dalam slice detil_order
	fmt.Println(brg)
}

func SaveDetilTransaksi(db *gorm.DB) { //save semua row barang ke dalam tabel detil_transaksi
	//insert row(s) ke tabel detil transaksi
	result := db.Table("detil_transaksi").Create(brg)
	if result.Error != nil {
		panic(result.Error)
	}

	//update stok tiap row barang
	barang := Barang{}
	fmt.Println("Kurangi stok tiap barang")

	for i := range brg {
		barang.ID = strconv.Itoa(brg[i].Id)
		log.Println("ID barang : " + barang.ID)

		_ = db.First(&barang, "id = ?", barang.ID) //ambil 1 row dengan ID tertentu

		log.Println("Stok barang : " + strconv.Itoa(barang.Stok))

		barang.Stok = barang.Stok - brg[i].Jumlah //update stok berdasarkan qty terjual

		log.Println("Stok dikurangi : " + strconv.Itoa(brg[i].Jumlah))
		log.Println("Stok barang setelah dipotong : " + strconv.Itoa(barang.Stok))

		_ = db.Save(&barang) //update data ke database
	}

	//kosongkan kembali slice data detil barang
	brg = []BarangDetilTransaksi{}

	fmt.Println("Data master-detil transaksi telah ditambahkan!")
}

// Cetak struk belanja

func CetakStruk(db *gorm.DB, idTransaksi string, userName string) {
	var total int
	var totalItem int
	var subtotal int

	//var transaksi []models.Transaksi

	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int    `gorm:"column:id_transaksi"`
		IdBarang    int    `gorm:"column:id_barang"`
		Jumlah      int    `gorm:"column:jumlah"`
		NamaBarang  string `gorm:"column:nama_barang"`
		Harga       int    `gorm:"column:harga"`
		//NamaBarang  string
		//Harga       int
	}

	var detil []DetilTrans

	err := db.Table("detil_transaksi").Select("detil_transaksi.id_transaksi, detil_transaksi.id_barang, detil_transaksi.jumlah, barang.nama_barang, barang.harga").Joins("join barang on detil_transaksi.id_barang=barang.id").Where("id_transaksi = ?", idTransaksi).Find(&detil).Error
	if err != nil {
		panic(err)
	}

	currentTime := time.Now() //var tanggal dan waktu saat ini
	fmt.Println("====================Struk belanja=======================")
	fmt.Println("ID Transaksi: " + idTransaksi + "               Kasir: " + userName) //idTransaksi ambil dari klausa where di atas
	fmt.Println("Tgl/Waktu: " + currentTime.Format("2006-01-02 15:04:05"))
	fmt.Println("========================================================")
	fmt.Printf("%-3s %-25s %-8s %-5s %-14s\n", "ID", "Nama Barang", "Harga", "Qty", "Subtotal")
	fmt.Println("========================================================")
	//pengecekan per row
	for _, row := range detil {
		subtotal = row.Harga * row.Jumlah
		fmt.Printf("%-3s %-25s %-8s %-5s %-14s\n", strconv.Itoa(row.IdBarang), row.NamaBarang, FormatAngka(row.Harga), strconv.Itoa(row.Jumlah), FormatAngka(subtotal))
		//fmt.Println( /*strconv.Itoa(row.IdTransaksi)+" | "+*/ strconv.Itoa(row.IdBarang) + " | " + row.NamaBarang + " | " + FormatAngka(row.Harga) + " | " + strconv.Itoa(row.Jumlah) + " | " + FormatAngka(subtotal))
		total += subtotal
		totalItem++
	}

	fmt.Println("========================================================")
	fmt.Printf("%-22s %-3s\n", "Total item purchased: ", FormatAngka(totalItem))
	//fmt.Println("Total item purchased: " + FormatAngka(totalItem))
	fmt.Printf("%-22s %-3s\n", "Total purchase: ", FormatRupiah(total))
	//fmt.Println("Total purchase: " + FormatRupiah(total))
	fmt.Println("========================================================")
}

// Cetak laporan transaksi hari ini
func CetakLaporanTransaksiSummaryToday(db *gorm.DB) {
	var today = time.Now().Format("2006-01-02")
	var todayString = time.Now().Format("2 Jan 2006")

	var total, totalItem int

	// 	Header
	header := []string{"ID", "Tanggal", "#Item(s)", "Total"}

	// Column widths
	w := []float64{12.0, 28.0, 20.0, 30.0}
	wSum := 0.0
	for _, v := range w {
		wSum += v
	}

	//setting orientation and size
	pdf := gofpdf.New("P", "mm", "A4", "")
	//set font style
	pdf.SetFont("Arial", "B", 16)
	//create a new page
	pdf.AddPage()

	//pdf.Cell(40, 10, "Laporan Stok Minimarket")
	pdf.WriteAligned(0, 15, "Laporan Transaksi Tanggal "+todayString, "C")
	pdf.Ln(0)

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//set font style for header
	pdf.SetFont("Arial", "B", 12)
	left := (210 - wSum) / 2
	pdf.SetX(left)

	//deklarasi struct dan query
	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int       `gorm:"column:id_transaksi"`
		CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
		Item        int       `gorm:"column:item"`
		Subtotal    int       `gorm:"column:subtotal"`
	}

	var detil []DetilTrans

	err := db.Table("detil_transaksi").Select("transaksi.created_at, detil_transaksi.id_transaksi, sum(jumlah) as item, sum(harga*jumlah) as subtotal").Joins("join barang on detil_transaksi.id_barang=barang.id").Joins("join transaksi on detil_transaksi.id_transaksi=transaksi.id_transaksi").Where("transaksi.created_at like ?", today+"%").Order("detil_transaksi.id_transaksi asc").Group("detil_transaksi.id_transaksi").Find(&detil).Error
	if err != nil {
		panic(err)
	}

	pdf.SetX(left)
	//print header
	for j, str := range header {
		pdf.CellFormat(w[j], 8, str, "1", 0, "C,M", false, 0, "")
	}
	pdf.Ln(-1)

	//set font style for data (not bold)
	pdf.SetFont("Arial", "", 12)

	// Data
	for _, row := range detil {
		pdf.SetX(left)
		pdf.CellFormat(w[0], 8, strconv.Itoa(row.IdTransaksi), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[1], 8, row.CreatedAt.Format("2006-01-02"), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[2], 8, FormatAngka(row.Item), "1", 0, "R,M", false, 0, "")
		pdf.CellFormat(w[3], 8, FormatAngka(row.Subtotal), "1", 0, "R,M", false, 0, "")

		//hitung total item dan subtotal secara incremental
		total += row.Subtotal
		totalItem++

		pdf.Ln(-1)
	}
	//Summary

	pdf.SetX(left)
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//harus selalu di set agar bisa nyambung dengan row sebelumnya
	pdf.SetX(left)
	pdf.CellFormat(60, 8, "Total Nilai Penjualan: ", "1", 0, "L,M", false, 0, "")
	pdf.CellFormat(30, 8, FormatAngka(total), "1", 0, "R,M", false, 0, "")

	//footer test
	pdf.SetFooterFunc(func() {
		currentTime := time.Now() //var tanggal dan waktu saat ini

		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 11)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Printed on: %v", currentTime.Format("2006-01-02 15:04:05")),
			"", 0, "L", false, 0, "")
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "R", false, 0, "")
	})

	pdf.SetX(left)

	//errs -> untuk cetak report
	errs := pdf.OutputFileAndClose("assets/laporan_transaksi_today.pdf") //URL silakan disetting

	if errs != nil {
		panic(errs)
	} else {
		fmt.Println("Laporan transaksi telah dibuat!" + "(laporan_transaksi_today.pdf)")
	}
}

// Cetak laporan transaksi per periode
func CetakLapTransaksiSummaryPeriode(db *gorm.DB, periode1Full string, periode2Full string, periode1 string, periode2 string) {
	var total, totalItem int

	// 	Header
	header := []string{"ID", "Tanggal", "#Item(s)", "Total"}

	// Column widths
	w := []float64{12.0, 28.0, 20.0, 27.0}
	wSum := 0.0
	for _, v := range w {
		wSum += v
	}

	//setting orientation and size
	pdf := gofpdf.New("P", "mm", "A4", "")
	//set font style
	pdf.SetFont("Arial", "B", 16)
	//create a new page
	pdf.AddPage()

	//pdf.Cell(40, 10, "Laporan Stok Minimarket")
	pdf.WriteAligned(0, 15, "Laporan Transaksi Dari "+periode1+" Sampai "+periode2, "C")
	pdf.Ln(0)

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//set font style for header
	pdf.SetFont("Arial", "B", 12)
	left := (210 - wSum) / 2
	pdf.SetX(left)

	//deklarasi struct dan query
	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int       `gorm:"column:id_transaksi"`
		CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
		Item        int       `gorm:"column:item"`
		Subtotal    int       `gorm:"column:subtotal"`
	}

	var detil []DetilTrans

	err := db.Table("detil_transaksi").Select("transaksi.created_at, detil_transaksi.id_transaksi, sum(jumlah) as item, sum(harga*jumlah) as subtotal").Joins("join barang on detil_transaksi.id_barang=barang.id").Joins("join transaksi on detil_transaksi.id_transaksi=transaksi.id_transaksi").Where("transaksi.created_at between ? and ?", periode1Full, periode2Full).Order("detil_transaksi.id_transaksi asc").Group("detil_transaksi.id_transaksi").Find(&detil).Error
	if err != nil {
		panic(err)
	}

	pdf.SetX(left)
	//print header
	for j, str := range header {
		pdf.CellFormat(w[j], 8, str, "1", 0, "C,M", false, 0, "")
	}
	pdf.Ln(-1)

	//set font style for data (not bold)
	pdf.SetFont("Arial", "", 12)

	// Data
	for _, row := range detil {
		pdf.SetX(left)
		pdf.CellFormat(w[0], 8, strconv.Itoa(row.IdTransaksi), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[1], 8, row.CreatedAt.Format("2006-01-02"), "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[2], 8, FormatAngka(row.Item), "1", 0, "R,M", false, 0, "")
		pdf.CellFormat(w[3], 8, FormatAngka(row.Subtotal), "1", 0, "R,M", false, 0, "")

		//hitung total item dan subtotal secara incremental
		total += row.Subtotal
		totalItem++

		pdf.Ln(-1)
	}
	//Summary

	pdf.SetX(left)
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//harus selalu di set agar bisa nyambung dengan row sebelumnya
	pdf.SetX(left)
	pdf.CellFormat(60, 8, "Total Nilai Penjualan: ", "1", 0, "L,M", false, 0, "")
	pdf.CellFormat(27, 8, FormatAngka(total), "1", 0, "R,M", false, 0, "")

	//footer test
	pdf.SetFooterFunc(func() {
		currentTime := time.Now() //var tanggal dan waktu saat ini

		// Position at 1.5 cm from bottom
		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 11)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Printed on: %v", currentTime.Format("2006-01-02 15:04:05")),
			"", 0, "L", false, 0, "")
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "R", false, 0, "")
	})

	pdf.SetX(left)

	//errs -> untuk cetak report
	errs := pdf.OutputFileAndClose("assets/laporan_transaksi_periode.pdf") //URL silakan disetting

	if errs != nil {
		panic(errs)
	}
}
