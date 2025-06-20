package models

type KategoriUser struct {
	ID        string `gorm:"primary_key;column:id"`
	LevelUser string `gorm:"primary_key;column:level_user"`
	User      []User `gorm:"foreignKey:id_level;references:id"` //1 to n
	//relasi many to one
	//beberapa ID yang sama misalnya kasir bisa dimiliki beberapa user
}

func (k *KategoriUser) TableName() string {
	return "kategori_user" //nama table pada db
}
