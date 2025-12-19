package controller

import (
	"net/http"

	"github.com/fsangopanta/demo-soft-code/common/utils"
	googleworkspace "github.com/fsangopanta/demo-soft-code/modules/google_workspace/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)


type (
	GoogleWorkspaceController interface {
		CreateMessage(ctx *gin.Context)
	}

	googleWorkspaceController struct {
		chatService googleworkspace.ChatService
	}
)


func NewWorkSpaceController(injector *do.Injector, chatService googleworkspace.ChatService ) *googleWorkspaceController{
	return &googleWorkspaceController{
		chatService: chatService,
	}
}

func (c *googleWorkspaceController) CreateMessage(ctx * gin.Context){
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