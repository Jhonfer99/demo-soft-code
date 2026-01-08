package googleworkspace

import (
	googleWorkspaceController "github.com/fsangopanta/demo-soft-code/modules/google_workspace/controller"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	serviceDeskControl := do.MustInvoke[googleWorkspaceController.GoogleWorkspaceController](injector)
	

	userRoutes := server.Group("/api/v1/google_workspace")
	{
		userRoutes.GET("", serviceDeskControl.CreateMessage)
		userRoutes.POST("", serviceDeskControl.GoogleWorkspaceHandler)
		
	}
}