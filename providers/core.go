package providers

import (
	"github.com/samber/do"

	serviceDeskController "github.com/fsangopanta/demo-soft-code/modules/service_desk/controller"
)



func RegisterDependencies(injector *do.Injector) {
	do.Provide(
		injector, func(i *do.Injector) (serviceDeskController.ServiceDeskController, error) {
			return serviceDeskController.NewServiceDeskController(i), nil
		},
	)
}