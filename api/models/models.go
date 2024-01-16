package models

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("alphanumwithspace", alphaNumWithSpace)
}

func alphaNumWithSpace(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	for _, char := range value {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) && !unicode.IsSpace(char) {
			return false
		}
	}

	return true
}

type List struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"type:varchar(100);not null" json:"title" form:"title" validate:"required,max=100,alphanumwithspace"`
	Description string    `gorm:"type:text;not null" json:"description" form:"description" validate:"required,max=1000"`
	Files       []File    `gorm:"foreignKey:ListID"`
	Sublists    []Sublist `gorm:"foreignKey:ListID;onDelete:CASCADE"`
}

type Sublist struct {
	ID          uint   `gorm:"primaryKey"`
	ListID      uint   `gorm:"index"`
	Title       string `gorm:"type:varchar(100);not null" json:"title" form:"title" validate:"required,max=100,alphanumwithspace"`
	Description string `gorm:"type:text;not null" json:"description" form:"description" validate:"required,max=1000"`
	Files       []File `gorm:"foreignKey:SublistID"`
}

type File struct {
	ID        uint   `gorm:"primaryKey"`
	ListID    *uint  `gorm:"index;onDelete:CASCADE"`
	SublistID *uint  `gorm:"index;onDelete:CASCADE"`
	FileName  string `gorm:"type:varchar(255);not null" json:"file_name" form:"file" validate:"omitempty,fileExtension=txt|pdf"`
}

func (l *List) Validate() error {
	return validate.Struct(l)
}

func (s *Sublist) Validate() error {
	return validate.Struct(s)
}

func (f *File) Validate() error {
	return validate.Struct(f)
}
