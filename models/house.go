package models

import (
	"github.com/jinzhu/gorm"
)

// House ...
type House struct {
	gorm.Model
	Channel string
	Price   float64
	Owner   uint
}
