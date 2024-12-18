package response

type Code int

// 常量初始化code值
const (
	SuccessCode Code = 200 + iota
	ParamFail
	ServerErrorCode     Code = 501
	UnprocessableEntity Code = 422 //客户端数据超出业务范围
	TokenError          Code = 401
)

// map用于存储每个code对应的提示信息
var codeMsgMap = map[Code]string{
	SuccessCode:     "操作成功",
	ParamFail:       "参数错误",
	ServerErrorCode: "服务端错误",
	TokenError:      "token错误",
}

// 用于获取code对应的提示信息
func (c Code) Msg() string {
	return codeMsgMap[c]
}
