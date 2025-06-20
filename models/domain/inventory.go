package domain

import (
	"time"

	"github.com/Go-Master-Code/kasir-web/models"

	"gorm.io/gorm"
)

type Inventory struct {
	ID             string                `gorm:"primary_key;column:id;autoIncrement"`
	IdKategori     int                   `gorm:"column:id_kategori"`
	NamaBarang     string                `gorm:"column:nama_barang"`
	Harga          int                   `gorm:"column:harga"`
	Stok           int                   `gorm:"column:stok"`
	CreatedAt      time.Time             `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time             `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt        `gorm:"column:deleted_at"` //tipe datanya bukan time.Time tapi gorm.DeletedAt -> penanda soft delete
	KategoriBarang models.KategoriBarang `gorm:"foreignKey:id_kategori;references:id"`
	JualBarang     []models.Transaksi    `gorm:"many2many:detil_transaksi;foreignKey:id;joinForeignKey:id_barang;references:id_transaksi;joinReferences:id_transaksi"`
	//format: tabel_many_to_many;foreignKey:PK_tabel_ini;joinForeignKey:nama_field_PK_di_tabel_detil;references:PK_tabel_master_lainnya;joinReferences:nama_field_PK_di_tabel_detil
}
