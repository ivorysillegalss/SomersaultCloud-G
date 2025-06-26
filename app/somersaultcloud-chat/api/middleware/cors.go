package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/cors"
	"net/http"
)

func DefaultCorsConfig() app.HandlerFunc {
	return cors.New(
		cors.Config{
			AllowAllOrigins:  true,
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", " Content-Length", " Accept-Encoding", " X-CSRF-Token", "Authorization", "accept", "origin", " Cache-Control", "X-Requested-With", "token"},
			AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodDelete, http.MethodPatch},
		})
}

// CustomCorsConfig 自定义跨域配置
func CustomCorsConfig(allowAllOrigins bool, allowCredentials bool, allowHeaders []string, allowMethods []string) app.HandlerFunc {
	return cors.New(
		cors.Config{
			AllowAllOrigins:  allowAllOrigins,
			AllowCredentials: allowCredentials,
			AllowHeaders:     allowHeaders,
			AllowMethods:     allowMethods,
		},
	)
}
