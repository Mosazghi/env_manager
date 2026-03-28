package handler

import "github.com/gin-gonic/gin"

func ToResponse(sucess bool, msg string, data any) gin.H {
	return gin.H{"sucess": sucess, "message": msg, "data": data}
}
