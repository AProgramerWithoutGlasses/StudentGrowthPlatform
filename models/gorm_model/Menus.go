package gorm_model

import (
	"gorm.io/gorm"
)

type Menus struct {
	gorm.Model
	ParentId      int    `json:"parentId"`
	TreePath      string `json:"treePath"`
	Name          string `json:"name"`
	Type          int    `json:"type"`
	RouteName     string `json:"routeName"`
	Path          string `json:"path"`
	Component     string `json:"component"`
	Perm          string `json:"perm"`
	Visible       int    `json:"visible"`
	Sort          int    `json:"sort"`
	Icon          string `json:"icon"`
	Redirect      string `json:"redirect"`
	Roles         string `json:"roles"`
	RequestUrl    string `json:"requestUrl "`
	RequestMethod string `json:"requestMethod"`
}
