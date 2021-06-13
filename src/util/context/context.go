package context

import (
	"github.com/labstack/echo"
	"net/http"
)

// Response 返回值
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// Success 成功
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Status: 200,
		Data:   data,
	})
}

// Error 错误
func Error(c echo.Context, status int, data string) error {
	return c.JSON(status, Response{
		Status: status,
		Data:   data,
	})
}
