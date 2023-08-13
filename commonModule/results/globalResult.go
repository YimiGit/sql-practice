package results

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// JsonResult 统一返回结果
type JsonResult struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
	Error   error  `json:"error"`
}

// Success 成功
func Success(message string, data any) *JsonResult {
	return &JsonResult{Message: message, Status: http.StatusOK, Data: data}
}

// Error 请求错误
func Error(message string, data any) *JsonResult {
	return &JsonResult{Message: message, Status: http.StatusBadRequest, Data: data}
}

// Fail 响应错误
func Fail(message string, err error) *JsonResult {
	return &JsonResult{Message: message, Status: http.StatusInternalServerError, Error: err}
}

// ResultHandle Controller统一处理Service返回的结果
func ResultHandle(ctx *gin.Context, result *JsonResult) {
	ctx.JSON(result.Status, result)
	if result.Error != nil {
		log.Println(result.Error)
		ctx.Abort()
	}
}
