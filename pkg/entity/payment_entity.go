package entity

type Payment struct {
	ID            int `gorm:"primaryKey; type:int; not null"`
	Name          string
	AccountNumber string
	Description   string
	PictureID     string         `gorm:"type:varchar(36)"`
	Picture       PaymentPicture `gorm:"foreignKey:PictureID"`
}

type PaymentPicture struct {
	ID  string `gorm:"primaryKey; type:varchar(36); not null"`
	Url int
	Alt string
}
