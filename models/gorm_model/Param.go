package gorm_model

import "gorm.io/gorm"

type Param struct {
	gorm.Model
	ParamsKey   string `json:"paramKey"`
	ParamsValue string `json:"paramValue"`
	MenuId      int    `json:"menuId"`
}
