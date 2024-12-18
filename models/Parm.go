package models

import (
	"github.com/dgrijalva/jwt-go"
	"studentGrow/models/gorm_model"
)

// Login 后台登录的结构体
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"verify" binging:""`
	Id       string `json:"verifyId"`
}

// Claims Token结构体
type Claims struct {
	User gorm_model.User
	jwt.StandardClaims
}

// TagAmount 文章tag和数量结构体
type TagAmount struct {
	Tag   int `json:"tag" gorm:"column:tag_id"`
	Count int
}

// StarBack 成长之星返回前端的结构体
type StarBack struct {
	Username           string `json:"username"`
	Frequency          int64  `json:"frequency"`
	Name               string `json:"name"`
	User_article_total int64  `json:"user_article_total"`
	Userfans           int64  `json:"userfans"`
	Score              int    `json:"score"`
	Quality            int    `json:"quality"`
	Hot                int    `json:"hot"`
	Status             bool   `json:"status"`
}

// StarClass 成长之星按班级分类的结构体
type StarClass struct {
	ClassName string   `json:"className"`
	ClassStar []string `json:"classStar"`
}

// StarGrade 年级之星数据的结构体
type StarGrade struct {
	GradeName  string `json:"gradeName"`
	GradeClass string `json:"gradeClass"`
}

// StarStu 前台成长之星数据的结构体
type StarStu struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	UserHeadshot string `json:"user_headshot"`
}

// Menu 返回前端侧边栏结构
type Menu struct {
	ID            int                `json:"id"`
	ParentId      int                `json:"parentId"`
	Name          string             `json:"menuName"`
	Type          int                `json:"type"`
	RouteName     string             `json:"routeName"`
	Path          string             `json:"routePath"`
	Perm          string             `json:"permissions"`
	Redirect      string             `json:"redirect"`
	Visible       int                `json:"isVisible"`
	Sort          int                `json:"sort"`
	FatherMenu    string             `json:"fatherMenu"`
	Component     string             `json:"componentPath"`
	RequestUrl    string             `json:"requestUrl"`
	RequestMethod string             `json:"requestMethod"`
	Params        []gorm_model.Param `json:"params"`
	Icon          string             `json:"icon"`
	Children      []Menu             `json:"children"`

	Status bool `json:"status"`
}

// Sidebar 返回前端侧边栏结构体
type Sidebar struct {
	Id        int                `json:"id"`
	ParentId  int                `json:"parentId"`
	Path      string             `json:"path"`
	Component string             `json:"component"`
	Redirect  string             `json:"redirect"`
	RouteName string             `json:"name"`
	Meta      Message            `json:"meta"`
	Params    []gorm_model.Param `json:"params"`
}

// Message 目录菜单信息
type Message struct {
	Name    string `json:"title"`
	Visible int    `json:"isVisible"`
	Icon    string `json:"icon"`
}

// RoleList 返回前端角色列表
type RoleList struct {
	Id       int    `json:"id"`
	RoleName string `json:"role_name"`
	RoleCode string `json:"role_Engname"`
}

// MenuList 返回前端菜单列表
type MenuList struct {
	Name     string     `json:"label"`
	Value    string     `json:"value"`
	Children []MenuList `json:"children"`
}
