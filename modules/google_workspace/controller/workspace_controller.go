package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fsangopanta/demo-soft-code/common/utils"
	"github.com/fsangopanta/demo-soft-code/modules/google_workspace/domain/dto"
	googleworkspace "github.com/fsangopanta/demo-soft-code/modules/google_workspace/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type (
	GoogleWorkspaceController interface {
		CreateMessage(ctx *gin.Context)
		GoogleWorkspaceHandler(ctx *gin.Context)
	}

	googleWorkspaceController struct {
		chatService googleworkspace.Messaging
	}
)

func NewWorkSpaceController(injector *do.Injector, chatService googleworkspace.Messaging) *googleWorkspaceController {
	return &googleWorkspaceController{
		chatService: chatService,
	}
}

func (c *googleWorkspaceController) CreateMessage(ctx *gin.Context) {
	msg, err := c.chatService.SendMessage(
		ctx.Request.Context(),
		"AAQAz-oaQ8g",
		"Mensaje que se debe enviar",
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	text := msg.GetFormattedText()
	if text == "" {
		text = msg.GetText()
	}
	res := utils.BuildResponseSuccess(text, msg)
	ctx.JSON(http.StatusAccepted, res)
}

func (c *googleWorkspaceController) GoogleWorkspaceHandler(ctx *gin.Context) {
	var event dto.WorkspaceEvent

	if err := ctx.ShouldBindJSON(&event); err != nil {
		log.Println(" Error parsing event:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"text": "Error leyendo evento",
		})
		return
	}

	log.Println("Evento recibido desde Chat")

	// if event.Chat != nil && event.Chat.MessagePayload != nil {
	handleMessage(ctx, event)
	// return
	// }

	// log.Println(" Evento no manejado (no messagePayload)")
	// ctx.JSON(http.StatusOK, dto.MessageResponse{
	// 	Text: "EVENTO RECIBIDO",
	// })
}

func handleMessage(ctx *gin.Context, event dto.WorkspaceEvent) {
	msg := event.Chat.MessagePayload.Message
	user := event.Chat.User

	if user.Type == "BOT" {
		ctx.JSON(http.StatusOK, gin.H{"text": ""})
		return
	}

	responseMessage := dto.MessageResponse{
		Text: "Hola " + user.DisplayName + "! Recibí tu mensaje: " + msg.Text,
		// Thread: dto.ThreadResponse{
		// 	Name: msg.Thread.Name,
		// },
	}

	jsonBytes, _ := json.MarshalIndent(responseMessage, "", "  ")
	log.Println(string(jsonBytes))

	// Enviar respuesta
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.JSON(http.StatusOK, responseMessage)
}
