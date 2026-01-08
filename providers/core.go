package providers

import (
	context "context"

	"github.com/samber/do"

	serviceGoogleController "github.com/fsangopanta/demo-soft-code/modules/google_workspace/controller"
	serviceGoogle "github.com/fsangopanta/demo-soft-code/modules/google_workspace/service"
	serviceDeskController "github.com/fsangopanta/demo-soft-code/modules/service_desk/controller"
)



func RegisterDependencies(injector *do.Injector) {

	do.Provide(injector, func(i *do.Injector) (serviceGoogle.Messaging, error) {
		ctx := context.Background() 
		return serviceGoogle.NewChatService(ctx)
	})

	do.Provide(
		injector, func(i *do.Injector) (serviceDeskController.ServiceDeskController, error) {
			return serviceDeskController.NewServiceDeskController(i), nil
		},
	)



	do.Provide(
		injector, func(i *do.Injector) (serviceGoogleController.GoogleWorkspaceController, error) {
				chatService := do.MustInvoke[serviceGoogle.Messaging](i)
			return serviceGoogleController.NewWorkSpaceController(i, chatService), nil
		},
	)

	
}