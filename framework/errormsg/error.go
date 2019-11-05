package errormsg

// ReturnType 返回类型
type ReturnType int32

const (
	// ReturnTypeSUCCESS 成功
	ReturnTypeSUCCESS ReturnType = 0
	// ReturnTypeSERVERBUSY 服务器忙
	ReturnTypeSERVERBUSY ReturnType = 1
	// ReturnTypeTOKENINVALID token无效
	ReturnTypeTOKENINVALID ReturnType = 2
	// ReturnTypeTOKENPARAMERR token参数错误
	ReturnTypeTOKENPARAMERR ReturnType = 3
	// ReturnTypeFAILRELOGIN 重登录失败
	ReturnTypeFAILRELOGIN ReturnType = 4
)
