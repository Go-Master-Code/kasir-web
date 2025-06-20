package test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Go-Master-Code/kasir-web/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenConnectionMaster() *gorm.DB { //Open connection isinya diambil dari web gorm
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dialect := "root:root@tcp(127.0.0.1:3306)/kasir?charset=utf8mb4&parseTime=True&loc=Local" //root:root artinya username=root, password=root
	db, err := gorm.Open(mysql.Open(dialect), &gorm.Config{
		//comment Logger di bawah agar perintah sql nya tidak muncul
		//Logger: logger.Default.LogMode(logger.Info), //ubah logger pada gorm: semua perintah SQL akan di log sebagai info

		//SkipDefaultTransaction: true,
		//PrepareStmt: true,
	})

	if err != nil {
		panic(err)
	}

	//connection pool
	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute) //maksimal digunakan 30 mins
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  //max waktu connection nganggur

	return db
}

var db = OpenConnectionMaster()

func TestBarang(t *testing.T) {
	fmt.Println("Cetak data barang")
	barang := CetakBarang() //tampung struct yang direturn dalam var listBarang

	for i := range barang {
		fmt.Printf("%-5s %-30s %-10s %-8s %-12s\n", barang[i].ID, barang[i].NamaBarang, FormatAngka(barang[i].Harga), strconv.Itoa(barang[i].Stok), barang[i].KategoriBarang.NamaKategori)
		//fmt.Println(barang[i].ID + " | " + barang[i].NamaBarang + " | " + models.FormatAngka(barang[i].Harga) + " | " + strconv.Itoa(barang[i].Stok) + " | " + barang[i].KategoriBarang.NamaKategori + " | ")
	}
}

func CetakBarang() []models.Barang {
	var barang []models.Barang

	err := db.Model(&models.Barang{}).Preload("KategoriBarang").Find(&barang).Error
	if err != nil {
		panic(err)
	}

	// looping di bawah untuk cetak data
	// for i := range barang {
	// 	fmt.Printf("%-5s %-30s %-10s %-8s %-12s\n", barang[i].ID, barang[i].NamaBarang, FormatAngka(barang[i].Harga), strconv.Itoa(barang[i].Stok), barang[i].KategoriBarang.NamaKategori)
	// 	//fmt.Println(barang[i].ID + " | " + barang[i].NamaBarang + " | " + models.FormatAngka(barang[i].Harga) + " | " + strconv.Itoa(barang[i].Stok) + " | " + barang[i].KategoriBarang.NamaKategori + " | ")
	// }
	return barang
}
