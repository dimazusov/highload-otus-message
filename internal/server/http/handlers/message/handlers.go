package message

import (
	"context"
	"log"
	"message/internal/pkg/pagination"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/minipkg/selection_condition"

	"message/internal/app"
	"message/internal/domain/message"
	"message/internal/pkg/apperror"
	"message/internal/server/http/api_error"
)

// @Summary get messages
// @Description get messages by params
// @ID get-messages
// @Accept json
// @Produce json
// @Param with_organization query boolean false "with organization"
// @Param per_page query int false "per page"
// @Param page query int false "page"
// @Success 200 {object} MessagesList
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_error.Error
// @Router /messages [get]
func GetMessagesHandler(c *gin.Context, app *app.App) {
	msg := message.Message{}
	if err := c.ShouldBindQuery(&msg); err != nil {
		c.JSON(http.StatusBadRequest, api_error.New(err))
		return
	}

	pag := pagination.Pagination{
		PerPage: pagination.DefaultPerPage,
		Page:    pagination.DefaultPage,
	}
	if err := c.ShouldBindWith(&pag, binding.Query); err != nil {
		c.JSON(http.StatusBadRequest, api_error.New(err))
		return
	}

	messages, err := app.Domain.Message.Service.GetMessages(context.Background(), &msg, &pag)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, api_error.New(apperror.ErrInternal))
		return
	}

	count, err := app.Domain.Message.Service.Count(context.Background(), &msg)
	if err != nil {
		log.Println(err)

		c.JSON(http.StatusInternalServerError, api_error.New(apperror.ErrInternal))
		return
	}

	c.JSON(http.StatusOK, MessagesList{
		Items: messages,
		Count: count,
	})
}

// @Summary get message by id
// @Description get message by id
// @ID get-message-by-id
// @Accept json
// @Produce json
// @Param id path boolean true "id"
// @Success 200 {object} message.Message
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_error.Error
// @Router /message/{id} [get]
func GetMessageHandler(c *gin.Context, app *app.App) {
	messageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, api_error.New(err))
		return
	}

	msg, err := app.Domain.Message.Service.GetMessage(context.Background(), uint(messageID))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, api_error.New(apperror.ErrInternal))
		return
	}

	c.JSON(http.StatusOK, msg)
}

// @Summary update message
// @Description update message
// @ID update-message
// @Accept json
// @Produce json
// @Param message body message.Message true "updatable message"
// @Success 200 {string} success
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_error.Error
// @Router /message [put]
func UpdateMessageHandler(c *gin.Context, app *app.App) {
	bdg := message.Message{}
	if err := c.ShouldBindJSON(&bdg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := app.Domain.Message.Service.Update(context.Background(), &bdg)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary create message
// @Description create message
// @ID create-message
// @Accept json
// @Produce json
// @Param message body message.Message true "creatable message"
// @Success 200 {object} Message
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_error.Error
// @Router /message/send [post]
func SendMessageHandler(c *gin.Context, app *app.App) {
	msg := message.Message{}
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := app.Domain.Message.Service.SendMessage(context.Background(), msg.FromUserID, msg.ToUserID, msg.Text)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary delete message by id
// @Description delete message by id
// @ID delete-message-by-id
// @Accept json
// @Produce json
// @Param id path boolean true "id"
// @Success 200 {string} success
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_error.Error
// @Router /message/{id} [delete]
func DeleteMessageHandler(c *gin.Context, app *app.App) {
	messageID, err := selection_condition.ParseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.Domain.Message.Service.Delete(context.Background(), messageID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
