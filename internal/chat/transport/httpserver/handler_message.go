package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (h HTTPServer) SendMessage(c *gin.Context) {
	var msgRequest SendMessageRequest
	var err error
	if err = c.ShouldBindJSON(&msgRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}
	fmt.Printf("\nGetChatMsgs: %+v\n\n", msgRequest)
	if err = msgRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{invaldRequest: err.Error()})
		return
	}
	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}
	createdAt := time.Now().Unix()
	msgRequest.SenderID = userCtx.ID()
	msgRequest.CreatedAt = createdAt
	msgRequest.UpdatedAt = createdAt
	// log.Printf("model req Message: %+v", msgRequest)
	domainMsg := toDomainMessage(msgRequest)
	// log.Printf("domain    Message: %+v", domainMsg)
	if err := h.msgService.SaveMsg(c, domainMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	usersID, err := h.chatService.GetUserIDsByChatID(c, msgRequest.ChatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content := fmt.Sprintf("check msgs in chatID: %d", msgRequest.ChatID)
	action := "new-message"
	h.wsHandler.Broadcast(content, action, msgRequest.ChatID, msgRequest.SenderID, usersID)
	c.JSON(http.StatusCreated, "message sent")
}

func (h HTTPServer) GetMessages(c *gin.Context) {
	var msgsRequest GetMessagesRequest
	var err error
	if err = c.ShouldBindJSON(&msgsRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err = msgsRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{invaldRequest: err.Error()})
		return
	}
	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}

	msgsRequest.UserID = userCtx.ID()

	msgsDomain, err := h.msgService.GetMessagesByChatID(
		c,
		msgsRequest.ChatID,
		msgsRequest.Limit,
		msgsRequest.Offset,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msgsResponse := make([]MessageResponse, len(msgsDomain))
	for i, msgDomain := range msgsDomain {
		msgsResponse[i] = toResponseMessage(msgDomain)
	}
	c.JSON(http.StatusOK, msgsResponse)
}

func (h HTTPServer) GetChatMessages(c *gin.Context) {
	var msgsRequest GetChatMessagesRequest
	var err error
	if err = c.ShouldBindJSON(&msgsRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	fmt.Printf("\nGetChatMsgs: %+v\n\n", msgsRequest)

	if err = msgsRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{invaldRequest: err.Error()})
		return
	}
	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}
	_ = userCtx
	// msgsRequest.UserID = userCtx.ID()

	// msgsDomain, err := h.msgService.GetMessagesByChatID(
	// 	c,
	// 	msgsRequest.ChatID,
	// 	msgsRequest.Limit,
	// 	msgsRequest.Offset,
	// )

	msgsDomain, err := h.msgService.GetChatMessages(
		c,
		msgsRequest.ChatID,
		msgsRequest.InitialMsgID,
		msgsRequest.Before,
		msgsRequest.After,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msgsResponse := make([]MessageResponse, len(msgsDomain))
	for i, msgDomain := range msgsDomain {
		msgsResponse[i] = toResponseMessage(msgDomain)
	}
	c.JSON(http.StatusOK, msgsResponse)
}
