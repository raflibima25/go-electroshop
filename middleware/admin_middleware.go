package middleware

import (
	"go-electroshop/internal/payload/response"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, response.SuccessResponse{
				ResponseStatus:  false,
				ResponseMessage: "Unauthorized",
				Data:            nil,
			})
			ctx.Abort()
			return
		}

		if mapClaims, ok := claims.(jwt.MapClaims); ok {
			isAdmin, exists := mapClaims["is_admin"]
			if !exists || !isAdmin.(bool) {
				ctx.JSON(http.StatusForbidden, response.SuccessResponse{
					ResponseStatus:  false,
					ResponseMessage: "Access denied: Admin privileges required",
					Data:            nil,
				})
				ctx.Abort()
				return
			}

			adminBool, ok := isAdmin.(bool)
			if !ok || !adminBool {
				ctx.JSON(http.StatusForbidden, response.SuccessResponse{
					ResponseStatus:  false,
					ResponseMessage: "Access denied: Admin privileges required",
					Data:            nil,
				})
				ctx.Abort()
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, response.SuccessResponse{
				ResponseStatus:  false,
				ResponseMessage: "Invalid claims format",
				Data:            nil,
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
