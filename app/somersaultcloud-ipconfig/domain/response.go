package domain

type IpConfigResp struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data"`
}

func SuccessResp(a any) *IpConfigResp {
	return &IpConfigResp{
		Message: "ok",
		Code:    0,
		Data:    a,
	}
}
