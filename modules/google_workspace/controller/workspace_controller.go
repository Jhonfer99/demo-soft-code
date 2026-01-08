package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/fsangopanta/demo-soft-code/common/utils"
	"github.com/fsangopanta/demo-soft-code/modules/google_workspace/domain/dto"
	googleworkspace "github.com/fsangopanta/demo-soft-code/modules/google_workspace/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"golang.org/x/oauth2/google"
	chat "google.golang.org/api/chat/v1"
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


func NewWorkSpaceController(injector *do.Injector, chatService googleworkspace.Messaging ) *googleWorkspaceController{
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

// func handleMessage(ctx *gin.Context, event dto.WorkspaceEvent) {
//     log.Println("=== DEBUG INICIO - Respuesta simplificada ===")
//     // Asegurarse de que el Content-Type esté establecido explícitamente para esta prueba
//     ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	ctx.Writer.Header().Set("ngrok-skip-browser-warning", "true")

//     // Respuesta estática y mínima
//     responseMessage := dto.MessageResponse{
//         Text: "Hola desde Soft Code Bot!",
//     }

//     jsonBytes, _ := json.MarshalIndent(responseMessage, "", "  ")
//     log.Println("📤 JSON simplificado que se enviará:")
//     log.Println(string(jsonBytes))

//     ctx.JSON(http.StatusOK, responseMessage)
//     log.Println("✅ Respuesta simplificada síncrona enviada a Google Chat")
// }

// func handleMessage(ctx *gin.Context) {
//     // 1. Headers críticos
//     ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
//     ctx.Writer.Header().Set("ngrok-skip-browser-warning", "true")

//     // 2. Respuesta estándar de Chat App (ya no Add-on)
//     response := gin.H{
//         "text": "Hola desde conexion sincrona",
//         "actionResponse": gin.H{
//             "type": "NEW_MESSAGE",
//         },
//     }

//     // 3. Usar PureJSON para evitar problemas de codificación
//     ctx.PureJSON(http.StatusOK, response)
// }


func handleMessage(ctx *gin.Context, event dto.WorkspaceEvent) {
	msg := event.Chat.MessagePayload.Message
	user := event.Chat.User

	log.Println("=== DEBUG INICIO ===")
	log.Printf("💬 Texto recibido: %s", msg.Text)
	log.Printf("👤 Usuario: %s (Type: %s)", user.DisplayName, user.Type)
	log.Printf("🧵 Thread Name: %s", msg.Thread.Name)
	log.Printf("📝 Message Name: %s", msg.Name)

	if user.Type == "BOT" {
		log.Println("Ignorando mensaje de bot")
		ctx.JSON(http.StatusOK, gin.H{"text": ""})
		return
	}

	// Crear respuesta
	responseMessage := dto.MessageResponse{
		Text: "Hola " + user.DisplayName + "! Recibí tu mensaje: " + msg.Text,
		// Thread: dto.ThreadResponse{
		// 	Name: msg.Thread.Name,
		// },
	}

	// DEBUG: Ver qué estás enviando realmente
	jsonBytes, _ := json.MarshalIndent(responseMessage, "", "  ")
	log.Println("📤 JSON que se enviará:")
	log.Println(string(jsonBytes))

	// DEBUG: Ver estructura completa
	log.Printf("📦 Estructura responseMessage: %+v", responseMessage)
	// log.Printf("🧵 Thread en respuesta: %+v", responseMessage.Thread)

	// Enviar respuesta
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	log.Printf("Encabezados de respuesta: %+v", ctx.Writer.Header())
	ctx.JSON(http.StatusOK, responseMessage)
	log.Println("✅ Respuesta síncrona enviada a Google Chat")
}

func sendGoogleMessage(ctx *gin.Context, user dto.User, msg dto.Message, space dto.Space){
	if user.Type == "BOT" {
		log.Println("Ignorando mensaje de bot")
		ctx.JSON(http.StatusOK, gin.H{"text": ""})
		return
	}

	botToken := getBearerToken("credentials.json") 

	reqBody := map[string]any{
		"text": "Hola " + user.DisplayName,
		"thread": map[string]string{
			"name": msg.Thread.Name,
		},
	}

	b, _ := json.Marshal(reqBody)
	url := "https://chat.googleapis.com/v1/" + space.Name + "/messages"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		log.Println(" Error creando request:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"text": "Error interno"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+botToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(" Error enviando mensaje:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"text": "Error interno"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error API Google Chat, status:", resp.Status)
	}
}

func getBearerToken(jsonPath string) string {
	ctx := context.Background()

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf(" Error leyendo archivo JSON: %v", err)
	}

	creds, err := google.CredentialsFromJSON(ctx, data, chat.ChatBotScope)
	if err != nil {
		log.Fatalf(" Error creando credenciales: %v", err)
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		log.Fatalf(" Error obteniendo token: %v", err)
	}

	return token.AccessToken

}
