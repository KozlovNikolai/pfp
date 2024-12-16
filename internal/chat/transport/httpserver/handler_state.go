package httpserver

import (
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/pkg/config"
	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
	"github.com/gin-gonic/gin"
)

// GetStates is ...
// GetStatesTags 		godoc
// @Summary			Получить список всех состояний.
// @Description		Return states list.
// @Tags			State
// @Security		BearerAuth
// @Param        	limit  query   string  true  "limit records on page" example(10) default(10)
// @Param       	offset  query   string  true  "start of record output" example(0) default(0)
// @Produce      	json
// @Success			200 {object} []StateResponse
// @failure			404 {string} err.Error()
// @Router			/admin/states [get]
func (h HTTPServer) GetStates(c *gin.Context) {
	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}
	userCtx, err = h.userService.GetUserByID(c, userCtx.ID())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNotFound.Error()})
		return
	}
	// fmt.Printf("\nuserCtx: %+v\n", userCtx)
	// check admin
	if userCtx.UserType() != config.AdminRole {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"invalid user id or role": domain.ErrAccessDenied.Error()},
		)
		return
	}
	limitQuery := c.Query("limit")
	offsetQuery := c.Query("offset")

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-limit": err.Error()})
		return
	}

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-offset": err.Error()})
		return
	}
	if limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"limit-must-be-greater-then-zero": ""})
		return
	}
	if offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"offset-must-be-greater-or-equal-then-zero": ""})
		return
	}

	states := h.stateService.GetAllStates(c)

	response := make([]StateResponse, 0, len(states))
	for _, state := range states {
		response = append(response, toResponseState(state))
	}

	c.JSON(http.StatusOK, response)
}
