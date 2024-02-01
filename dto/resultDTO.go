package dto

// ResultDTO 返回给前端的响应类
type ResultDTO struct {
	//响应码
	Code    int
	Message string
	//响应信息
	Data any
	//响应数据 以json格式发送
}

// SuccessResp 方法设置 ResultDTO 为一个标准的成功响应
func (r *ResultDTO) SuccessResp(code int, message string, data interface{}) *ResultDTO {
	r.Code = code
	r.Message = message
	r.Data = data
	return r
}

// Failed 方法设置 ResultDTO 为一个标准的失败响应 跟成功是一样的 只是名字不同 controller中好分辨
func (r *ResultDTO) FailResp(code int, message string, data interface{}) *ResultDTO {
	r.Code = code
	r.Message = message
	r.Data = data
	return r
}
