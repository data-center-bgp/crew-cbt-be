package models

import (
	"quiz-crew/config"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nama         string        `gorm:"type:varchar(100);not null"`
	Nik          string        `gorm:"type:varchar(100);not null"`
	Jabatan      string        `gorm:"type:varchar(100);not null"`
	Perusahaan   string        `gorm:"type:varchar(100);not null"`
	QuizAttempts []QuizAttempt `gorm:"foreignKey:UserID"`
}

type QuizAttempt struct {
	gorm.Model
	UserID         uint `gorm:"not null"`
	QuizCategoryID uint `gorm:"not null"`
	Score          int  `gorm:"not null"`
	PassingStatus  bool `gorm:"not null"`
}

type QuizCategory struct {
	gorm.Model
	Nama      string     `gorm:"type:varchar(100);not null"`
	Questions []Question `gorm:"foreignKey:QuizCategoryID"`
}

type Question struct {
	gorm.Model
	QuizCategoryID uint     `gorm:"not null"`
	Text           string   `gorm:"type:text;not null"`
	ImageUrl       string   `gorm:"default:null"`
	Answers        []Answer `gorm:"foreignKey:QuestionID;constraint:onDelete:CASCADE"`
}

type Answer struct {
	gorm.Model
	QuestionID uint   `gorm:"not null"`
	Text       string `gorm:"type:text;not null"`
	IsCorrect  bool   `gorm:"not null"`
}

func Migrate() {
	db := config.DB
	db.AutoMigrate((&User{}), (&QuizAttempt{}), (&QuizCategory{}), (&Question{}), (&Answer{}))
}
