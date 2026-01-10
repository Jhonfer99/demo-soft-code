package infrastructure

import (
	"log"

	models "github.com/fsangopanta/demo-soft-code/modules/chat/domain/models"
	usecases "github.com/fsangopanta/demo-soft-code/modules/chat/domain/usecases"
	"github.com/gin-gonic/gin"
	// interfaces "github.com/fsangopanta/demo-soft-code/modules/chat/domain/usecases/interfaces"
)

type GoogleController interface {
	// HandleAddWorkspace(ctx *gin.Context)
	HandleMessage(ctx *gin.Context)
}

type googleController struct {
	processIncomingUC *usecases.ProcessIncomingMessageUseCase
}

func NewGoogleController(
	processIncomingUC *usecases.ProcessIncomingMessageUseCase,
) GoogleController {
	return &googleController{
		processIncomingUC: processIncomingUC,
	}
}

func (c *googleController) HandleMessage(ctx *gin.Context) {
	var req models.ChatEvent
	var msg models.IncomingMessage
	ctx.Writer.Header().Set("ngrok-skip-browser-warning", "true")
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	msg.Text = req.Message.Text
	msg.UserId = req.User.DisplayName

	log.Println("Argumentos recibidos")
	log.Println(req)
	out, err := c.processIncomingUC.Handle(
		ctx.Request.Context(),
		msg,
		nil,
	)

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, out)
}

// func (c *googleController) HandleAddWorkspace(ctx *gin.Context) {
// 	var req inbound.AddWorkspaceRequest

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err := c.registerWorkspaceUC.Execute(ctx.Request.Context(), req)
// 	if err != nil {
// 		ctx.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(201, gin.H{"message": "workspace registered"})
// }
