package gorm_model

import (
	"gorm.io/gorm"
)

type UserEditRecord struct {
	gorm.Model
	Username       string `json:"username"`
	EditUsername   string `json:"edited_username"`
	OldName        string `json:"old_Name"`
	NewName        string `json:"new_Name"`
	OldUsername    string `json:"old_Username"`
	NewUsername    string `json:"new_Username"`
	OldPassword    string `json:"old_Password"`
	NewPassword    string `json:"new_Password"`
	OldGender      string `json:"old_Gender"`
	NewGender      string `json:"new_Gender"`
	OldClass       string `json:"old_class"`
	NewClass       string `json:"new_class"`
	OldPlustime    string `json:"old_Plustime"`
	NewPlustime    string `json:"new_Plustime"`
	OldPhoneNumber string `json:"old_phone_number"`
	NewPhoneNumber string `json:"new_phone_number"`
}
