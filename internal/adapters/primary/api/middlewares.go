package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (api *Adapter) authenticate(ctx *gin.Context) {
	authHeader := strings.Split(ctx.GetHeader("Authorization"), " ")
	if len(authHeader) < 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := authHeader[1]

	userId, err := api.app.ValidateJWTToken(token)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("user_id", userId)

	ctx.Next()
}
