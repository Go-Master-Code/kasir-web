package config

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Koneksi database

func OpenConnectionMaster() *gorm.DB { //Open connection isinya diambil dari web gorm
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dialect := "root:root@tcp(127.0.0.1:3306)/kasir?charset=utf8mb4&parseTime=True&loc=Local" //root:root artinya username=root, password=root
	db, err := gorm.Open(mysql.Open(dialect), &gorm.Config{
		//comment Logger di bawah agar perintah sql nya tidak muncul
		Logger: logger.Default.LogMode(logger.Info), //ubah logger pada gorm: semua perintah SQL akan di log sebagai info

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
