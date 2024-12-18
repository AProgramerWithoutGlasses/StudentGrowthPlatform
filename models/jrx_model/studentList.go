package jrx_model

import "time"

// 学生信息表（为贴合apifox的字段，备用）en
type StuMesStruct struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Class       string `json:"class"`
	Year        string `json:"year"`
	Gender      string `json:"gender"`
	Telephone   string `json:"telephone"`
	Ban         bool   `json:"ban"`
	ManagerType string `json:"manager_type"`
}

// 用于入学年份下拉框
type YearStruct struct {
	Id_Year string `json:"value"`
	Year    string `json:"label"`
}

// 用于班级下拉框
type ClassStruct struct {
	Id_class string `json:"value"`
	Class    string `json:"label"`
}

// ResponseStruct 返回查询结果给前端
type ResponseStruct struct {
	Role            string         `json:"role"`
	Year            []YearStruct   `json:"year"`
	Class           []ClassStruct  `json:"class"`
	StuInfo         []StuMesStruct `json:"stuInfo"`
	AllStudentCount int            `json:"allStudentCount"`
}

// queryParmaStruct 用于获取查询参数
type QueryParmaStruct struct {
	Year          string `json:"year"`
	Class         string `json:"class"`
	Gender        string `json:"gender"`
	IsDisable     bool   `json:"is_disable"`
	SearchSelect  string `json:"search_select"`
	SearchMessage string `json:"search_message"`
	IsManager     bool   `json:"is_manager"`
}

// 修改学生信息
type ChangeStuMesStruct struct {
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Gender      string    `json:"gender"`
	Class       string    `json:"class"`
	PhoneNumber string    `json:"telephone"`
	Year        string    `json:"year"`      // 入学年份,负责接收请求参数，然后转换后赋值给下方 PlusTime
	PlusTime    time.Time `json:"plus_time"` // 实际对应数据库的入学年份
	OldUsername string    `json:"oldUsername"`
}

// 学生信息表（为贴合apifox的字段，备用）(year int)
type StuMesYearIntStruct struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Class     string `json:"class"`
	Year      int    `json:"year"`
	Gender    string `json:"gender"`
	Telephone string `json:"telephone"`
	Ban       bool   `json:"ban"`
	IsManager bool   `json:"isManager"`
}

type SelectedStuMesStruct struct {
	Selected_students []StuMesStruct `json:"selected_students"`
}

type MyTokenMes struct {
	MyUsername string
	MyId       int
	MyRole     string
	MyClass    string
}

type SetStuManagerModel struct {
	Student     SetStuManagerInnerModel `json:"student"`
	ManagerType string                  `json:"managerType"`
}

type SetStuManagerInnerModel struct {
	Username string `json:"username" binding:"required"`
	Year     string `json:"year" binding:"required"`
}
