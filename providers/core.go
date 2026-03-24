package providers

import (
	"github.com/fsangopanta/demo-soft-code/modules/chat/domain/usecases"
	"github.com/fsangopanta/demo-soft-code/modules/chat/infrastructure"
	chatInfrastructure "github.com/fsangopanta/demo-soft-code/modules/chat/infrastructure"
	serviceDeskController "github.com/fsangopanta/demo-soft-code/modules/service_desk/controller"
	"github.com/samber/do"
)

func RegisterDependencies(injector *do.Injector) {

	// do.Provide(injector, func(i *do.Injector) (serviceGoogle.Messaging, error) {
	// 	ctx := context.Background()
	// 	return serviceGoogle.NewChatService(ctx)
	// })

	do.Provide(
		injector, func(i *do.Injector) (serviceDeskController.ServiceDeskController, error) {
			return serviceDeskController.NewServiceDeskController(i), nil
		},
	)

	do.Provide(injector, func(i *do.Injector) (usecases.Processor, error) {
		return &chatInfrastructure.LocalProcessor{}, nil
	})

	do.Provide(injector, func(i *do.Injector) (chatInfrastructure.GoogleController, error) {
		uc := do.MustInvoke[*usecases.ProcessIncomingMessageUseCase](i)
		return infrastructure.NewGoogleController(uc), nil
	})

	do.Provide(injector, func(i *do.Injector) (*usecases.ProcessIncomingMessageUseCase, error) {
		return usecases.NewProcessIncomingMessageUseCase(
			do.MustInvoke[usecases.Processor](i),
		), nil
	})

}
