package infrastructure

import (
	googleController "github.com/fsangopanta/demo-soft-code/modules/chat/infrastructure/google/rest"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	googleController := do.MustInvoke[googleController.GoogleController](injector)

	userRoutes := server.Group("/api/v1/google_chat")
	{
		userRoutes.POST("/messaging", googleController.HandleMessage)
		// userRoutes.POST("", serviceDeskControl.GoogleWorkspaceHandler)

	}
}
