package middleware

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
	"net/http"
)

func DefaultCorsConfig(h *server.Hertz) {
	h.Use(
		cors.New(
			cors.Config{
				AllowAllOrigins:  true,
				AllowCredentials: true,
				AllowHeaders:     []string{"Content-Type", " Content-Length", " Accept-Encoding", " X-CSRF-Token", "Authorization", "accept", "origin", " Cache-Control", "X-Requested-With", "token"},
				AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodDelete, http.MethodPatch},
			}),

		//TODO REMOVE
		//func(ctx context.Context, c *app.RequestContext) {
		//	if consts.MethodOptions == string(c.Request.Method()[0]) {
		//		c.AbortWithStatus(204)
		//	} else {
		//		c.Next(ctx)
		//	}
		//},
	)
}
