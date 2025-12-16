package user

import (
	serviceDeskController "github.com/fsangopanta/demo-soft-code/modules/service_desk/controller"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	serviceDeskControl := do.MustInvoke[serviceDeskController.ServiceDeskController](injector)
	

	userRoutes := server.Group("/api/v1/service-desk")
	{
		userRoutes.GET("", serviceDeskControl.Hello)
		
	}
}