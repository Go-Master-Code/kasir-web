package models

import (
	"gorm.io/gorm"
)

type KategoriBarang struct {
	ID           int      `gorm:"primary_key;column:id;autoIncrement"`
	NamaKategori string   `gorm:"primary_key;column:nama_kategori"`
	Barang       []Barang `gorm:"foreignKey:id;references:id"`
	//relasi many to one terhadap barang
	//ID kategori barang yang sama dapat dimiliki beberapa barang
}

func (k *KategoriBarang) TableName() string {
	return "kategori_barang" //nama table pada db
}

// kategori barang
func TambahKategoriBarang(db *gorm.DB, kategori string) { //untuk memasukkan 1 row (data) ke db
	kategoribarang := KategoriBarang{ //masukkan data (single) pada struct
		NamaKategori: kategori, //ID tidak didefinisikan karena autoincrement
	}

	err := db.Create(&kategoribarang).Error
	if err != nil {
		panic(err)
	}
}

func TampilkanKategoriBarang(db *gorm.DB) []KategoriBarang { //return value: slice of KategoriBarang
	var kategoriBarang []KategoriBarang

	err := db.Find(&kategoriBarang).Error
	if err != nil {
		panic(err)
	}

	return kategoriBarang
}
