package models

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

type Barang struct {
	ID             string         `gorm:"primary_key;column:id;autoIncrement"`
	IdKategori     int            `gorm:"column:id_kategori"`
	NamaBarang     string         `gorm:"column:nama_barang"`
	Harga          int            `gorm:"column:harga"`
	Stok           int            `gorm:"column:stok"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"` //tipe datanya bukan time.Time tapi gorm.DeletedAt -> penanda soft delete
	KategoriBarang KategoriBarang `gorm:"foreignKey:id_kategori;references:id"`
	JualBarang     []Transaksi    `gorm:"many2many:detil_transaksi;foreignKey:id;joinForeignKey:id_barang;references:id_transaksi;joinReferences:id_transaksi"`
	//format: tabel_many_to_many;foreignKey:PK_tabel_ini;joinForeignKey:nama_field_PK_di_tabel_detil;references:PK_tabel_master_lainnya;joinReferences:nama_field_PK_di_tabel_detil
}

func (b *Barang) TableName() string {
	return "barang" //nama table pada db nya adalah user_logs
}

func TampilkanBarang(db *gorm.DB) []Barang { //return value slice [] of Barang
	var barang []Barang

	err := db.Model(&Barang{}).Preload("KategoriBarang").Find(&barang).Error
	if err != nil {
		panic(err)
	}

	return barang
}

func TampilkanBarangOrderByNama(db *gorm.DB) []Barang { //return value slice [] of Barang
	var barang []Barang

	err := db.Model(&Barang{}).Preload("KategoriBarang").Order("nama_barang asc").Find(&barang).Error
	if err != nil {
		panic(err)
	}

	return barang
}

func TampilkanBarangById(db *gorm.DB, idBarang string) ([]Barang, int) { //return value slice [] of Barang dan selected idKategori untuk dropdown
	var barang []Barang

	err := db.Model(&Barang{}).Preload("KategoriBarang").Where("id = ?", idBarang).Find(&barang).Error
	if err != nil {
		panic(err)
	}

	idKategori := barang[0].IdKategori //mengambil idKategori dari produk terpilih untuk dijadikan default value pada dropdown

	return barang, idKategori
}

func TampilkanBarangPerKategori(db *gorm.DB, idKategori int) {
	var barang []Barang
	/*
		2 query:
		SELECT * FROM `addresses` WHERE `addresses`.`user_id` IN ('1','20','50','10','11','12','13','14','2','21','3','4','5','6','7','8','9')
		SELECT `users`.`id`,`users`.`password`,`users`.`first_name`,`users`.`middle_name`,`users`.`last_name`,`users`.`created_at`,`users`.`updated_at`,`Wallet`.`id` AS `Wallet__id`,`Wallet`.`user_id` AS `Wallet__user_id`,`Wallet`.`balance` AS `Wallet__balance`,`Wallet`.`created_at` AS `Wallet__created_at`,`Wallet`.`updated_at` AS `Wallet__updated_at` FROM `users` LEFT JOIN `wallets` `Wallet` ON `users`.`id` = `Wallet`.`user_id`
	*/
	err := db.Model(&Barang{}).Preload("KategoriBarang").Where("barang.id_kategori = ?", idKategori).Find(&barang).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("==========================Data Stok barang==========================")
	fmt.Printf("%-5s %-30s %-10s %-8s %-12s\n", "ID", "Nama Barang", "Harga", "Stok", "Kategori")
	//fmt.Println("ID | Nama Barang | Harga | Stok | Kategori |")
	fmt.Println("====================================================================")

	for i := range barang {
		fmt.Printf("%-5s %-30s %-10s %-8s %-12s\n", barang[i].ID, barang[i].NamaBarang, FormatAngka(barang[i].Harga), strconv.Itoa(barang[i].Stok), barang[i].KategoriBarang.NamaKategori)
		//fmt.Println(barang[i].ID + " | " + barang[i].NamaBarang + " | " + models.FormatAngka(barang[i].Harga) + " | " + strconv.Itoa(barang[i].Stok) + " | " + barang[i].KategoriBarang.NamaKategori + " | ")
	}
}

func TambahBarang(db *gorm.DB, namaBarang string, harga int, idKategori int) {
	barang := Barang{ //masukkan data (single) pada struct
		//ID:         id, auto_increment
		NamaBarang: namaBarang,
		Harga:      harga,
		//Stok:       0, default value di db sudah di set 0
		IdKategori: idKategori,
	}
	//db.Last()
	err := db.Create(&barang).Error
	if err != nil {
		panic(err)
	}
}

func UpdateBarang(db *gorm.DB, idBarang int, namaBarang string, harga int, stok int, idKategori int) {
	barang := Barang{}
	_ = db.First(&barang, "id = ?", idBarang) //ambil 1 row dengan ID pada parameter

	barang.NamaBarang = namaBarang
	barang.Harga = harga
	barang.Stok = stok
	barang.KategoriBarang.ID = idKategori

	err := db.Save(&barang).Error
	if err != nil {
		panic(err)
	}
}

func DeleteBarang(db *gorm.DB, idBarang string) {
	/* -----Cek apakah barang sudah pernah masuk detil transaksi -------
	type DetilTrans struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
		IdTransaksi int `gorm:"column:id_transaksi"`
		IdBarang    int `gorm:"column:id_barang"`
	}

	var detil []DetilTrans

	//cek apakah id barang ini ada di tabel detil transaksi
	err := db.Table("detil_transaksi").Select("id_barang").Where("id_barang = ?", idBarang).Limit(1).Find(&detil).Error
	if err != nil {
		panic(err)
	}

	if len(detil) < 1 {
		//Jika data barang tidak ditemukan di detil_transaksi, hapus data barangnya
		err = db.Where("id = ?", idBarang).Delete(&Barang{}).Error
		if err != nil {
			panic(err)
		}
	} else {
		//jika id barang ditemukan di salah satu row pada detil_transaksi
		log.Println("Gagal menghapus data barang " + idBarang + "! (Sudah pernah transaksi)")
	}
	-----Akhir pengecekan dan delete data barang ------	*/

	//Menggunakan metode soft delete
	//Field deleted at akan terisi, record barang masih ada, tidak dihapus
	err := db.Where("id = ?", idBarang).Delete(&Barang{}).Error
	if err != nil {
		panic(err)
	}
}

// Update stok barang yang ada pada detil transaksi
func UpdateStokBarangDetilTransaksi(db *gorm.DB) { //save semua row barang ke dalam tabel detil_transaksi
	barang := Barang{}
	//fmt.Println("Kurangi stok tiap barang")

	for i := range brg {
		barang.ID = strconv.Itoa(brg[i].Id)
		log.Println("ID barang pertama: " + barang.ID)

		_ = db.First(&barang, "id = ?", barang.ID) //ambil 1 row dengan ID tertentu

		log.Println("Stok barang : " + strconv.Itoa(barang.Stok))

		barang.Stok = barang.Stok - brg[i].Jumlah //update stok berdasarkan qty terjual

		log.Println("Stok dikurangi : " + strconv.Itoa(brg[i].Jumlah))
		log.Println("Stok barang setelah dipotong : " + strconv.Itoa(barang.Stok))

		_ = db.Save(&barang) //update data ke database
		//fmt.Println("Stok barang " + barang.ID + " telah diupdate menjadi: " + strconv.Itoa(barang.Stok))
	}
}

func TampilkanBarangSedikit(db *gorm.DB) []Barang {
	var barang []Barang

	//tampilkan barang yang stoknya < 10
	result := db.Model(Barang{}).Preload("KategoriBarang").Where("stok < ?", "10").Find(&barang).Error
	if result != nil {
		panic(result)
	}

	return barang

	/* Untuk versi CLI
	fmt.Println("======================Data Stok barang menipis======================")
	fmt.Printf("%-5s %-30s %-10s %-8s %-12s\n", "ID", "Nama Barang", "Harga", "Stok", "Kategori")
	//fmt.Println("ID | Nama Barang | Harga | Stok | Kategori |")
	fmt.Println("====================================================================")

	for i := range barang {
		fmt.Printf("%-5s %-30s %-10s %-8s %-12s\n", barang[i].ID, barang[i].NamaBarang, FormatAngka(barang[i].Harga), strconv.Itoa(barang[i].Stok), barang[i].KategoriBarang.NamaKategori)

		//fmt.Println(barang[i].ID + " | " + barang[i].NamaBarang + " | " + FormatAngka(barang[i].Harga) + " | " + strconv.Itoa(barang[i].Stok) + " | " + barang[i].KategoriBarang.NamaKategori)
	}
	*/
}

// Cetak laporan stok barang
func CetakLaporanStokBarang(db *gorm.DB) {
	var todayString = time.Now().Format("2 Jan 2006")
	// 	Header
	header := []string{"ID", "Nama Barang", "Kategori", "Harga", "Stok"}

	// Column widths
	w := []float64{12.0, 88.0, 30.0, 22.0, 18.0}
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
	pdf.WriteAligned(0, 15, "Laporan Stok Barang Per "+todayString, "C")
	pdf.Ln(0)

	//pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.Ln(-1) //line spacing

	//set font style for header
	pdf.SetFont("Arial", "B", 12)
	left := (210 - wSum) / 2
	pdf.SetX(left)
	var barang []Barang

	//select * from barang + barang.nama_kategori
	err := db.Model(Barang{}).Preload("KategoriBarang").Find(&barang).Error
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
	for _, b := range barang {
		pdf.SetX(left)
		pdf.CellFormat(w[0], 8, b.ID, "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[1], 8, b.NamaBarang, "1", 0, "L,M", false, 0, "")
		pdf.CellFormat(w[2], 8, b.KategoriBarang.NamaKategori, "1", 0, "C,M", false, 0, "")
		pdf.CellFormat(w[3], 8, FormatAngka(b.Harga), "1", 0, "R,M", false, 0, "")
		pdf.CellFormat(w[4], 8, FormatAngka(b.Stok), "1", 0, "R,M", false, 0, "")

		pdf.Ln(-1)
	}
	pdf.SetX(left)
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")

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
	errs := pdf.OutputFileAndClose("assets/laporan_stok_barang.pdf") //URL silakan disetting
	if errs != nil {
		panic(errs)
	} else {
		fmt.Println("Laporan stok barang telah dibuat! (laporan_stok_barang.pdf)")
	}
}
