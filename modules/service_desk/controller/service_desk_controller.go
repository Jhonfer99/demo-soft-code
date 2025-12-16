package controller

import (
	"net/http"

	"github.com/fsangopanta/demo-soft-code/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)


type (
	ServiceDeskController interface {
		Hello(ctx *gin.Context)
	}

	serviceDeskController struct {

	}
)


func NewServiceDeskController(injector *do.Injector) *serviceDeskController{
	return &serviceDeskController{}
}

func (c *serviceDeskController) Hello(ctx * gin.Context){
	res := utils.BuildResponseSuccess("MENSAJE HOLA MUNDO","WORKS WITH GINN!!")
	ctx.JSON(http.StatusAccepted, res)
}

