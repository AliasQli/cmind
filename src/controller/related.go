package controller

import (
	"conWord/src/model"
	"conWord/src/util/context"
	"github.com/labstack/echo"
	"net/http"
)

func GetSingleWordRelated(c echo.Context) error {
	word := c.QueryParam("word")
	data, err := model.GetRelatedWordList(word)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, err.Error())
	}
	return context.Success(c, data)
}
