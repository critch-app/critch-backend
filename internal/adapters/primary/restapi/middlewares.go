package restapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (api *Adapter) authenticate(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	userId, err := api.app.ValidateJWTToken(token)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("user_id", userId)

	ctx.Next()
}
