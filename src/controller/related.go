package controller

import (
	"conWord/src/model"
	"conWord/src/util/context"
	"github.com/labstack/echo"
	"net/http"
)

func GetSingleWordRelated(c echo.Context) error {
	word := c.QueryParam("word")
	categorized := c.QueryParam("categorized") == "true"
	var data interface{}
	var err error
	if categorized {
		data, err = model.GetSingleRelatedCategorized(word)
	} else {
		data, err = model.GetSingleRelatedWordList(word)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, err.Error())
	}
	return context.Success(c, data)
}
