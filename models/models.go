package models

import (
	"gorm.io/gorm"

	"mahayoga.org/services"
)

var db *gorm.DB

type Cards struct {
	gorm.Model
	UID         string `gorm:"primaryKey;not null;unique;size:40" json:"uid"`
	Title       string `gorm:"not null;size:150" json:"title"`
	Language    int    `gorm:"not null" json:"language"`
	Description string `gorm:"not null;size:500" json:"description"`
	Cover       []byte `gorm:"not null;" json:"cover"`
	Content     []byte `gorm:"not null;" json:"content"`
}

func init() {
	services.ConnectDB()
	db = services.GetDB()
	db.AutoMigrate(&Cards{})
}

func (card *Cards) CreateCard() error {
	return db.Create(&card).Error
}

// update cards where uid = cardID
func (card *Cards) UpdateCard(cardID string) error {
	return db.Where("uid = ?", cardID).Updates(card).Error
}

func DeleteCard(cardID string) error {
	return db.Where("uid = ?", cardID).Delete(&Cards{}).Error
}

func GetCards() ([]Cards, error) {
	var cards []Cards
	err := db.Find(&cards).Error
	return cards, err
}

func GetCardsUnscoped(modified string) ([]Cards, error) {
	var cards []Cards

	if modified != "" {
		err := db.Unscoped().Where("updated_at >= ?", modified).Find(&cards).Error
		return cards, err
	}

	err := db.Unscoped().Find(&cards).Error
	return cards, err
}
