package models

import "gorm.io/gorm"

type User struct {
	ID           string       `gorm:"primary_key;column:id"`
	IdLevel      string       `gorm:"column:id_level"`
	Password     string       `gorm:"column:password"`
	KategoriUser KategoriUser `gorm:"foreignKey:id_level;references:id"`
	//Transaksi    []Transaksi  `gorm:"foreignKey:id;references:id"`
	//relasi one to many
	//1 user menangani banyak transaksi
}

func (u *User) TableName() string {
	return "user" //nama table pada db nya adalah user_logs
}

func TampilkanUser(db *gorm.DB) []User { //return value slice [] of Barang
	var user []User

	err := db.Model(&User{}).Preload("KategoriUser").Find(&user).Error
	if err != nil {
		panic(err)
	}

	return user
}

func ValidasiUser(db *gorm.DB, idBarang string, password string) (string, string) { //returnkan username dan password berdasarkan login
	var user []User

	err := db.Model(&User{}).Where("id = ? and password =?", idBarang, password).Find(&user).Error
	if err != nil {
		panic(err)
	}

	if len(user) > 0 {
		//jika data user ketemu
		idUser := user[0].ID
		pwd := user[0].Password

		return idUser, pwd
	} else {
		//jika data user tidak ketemu
		return "", ""
	}
}
