package gorm_model

import "gorm.io/gorm"

type Advice struct {
	gorm.Model
	Username string `json:"username"`
	Advice   string `json:"advice"`
}
