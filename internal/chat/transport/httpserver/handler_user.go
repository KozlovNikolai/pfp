package httpserver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/pkg/config"
	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
)

const (
	errorGetUsers = "error get users"
)

// GetUser is ...
// GetUserTags 		godoc
// @Summary			Посмотреть пользователя по его id или логину.
// @Description		Получить пользователя по его id ли логину.
// @Param        	id  query   string  false  "id of the user" example(1) default(1)
// @Param        	login  query   string  false  "login of the user" example(cmd@cmd.ru) default(cmd@cmd.ru)
// @Tags			User
// @Security		BearerAuth
// @Success			200 {object} UserResponse
// @failure			404 {string} err.Error()
// @Router			/auth/user [get]
func (h HTTPServer) GetUser(c *gin.Context) {
	var userRequest UserRequest
	profileQuery := c.Query("profile")
	idQuery := c.Query("id")
	loginQuery := c.Query("login")

	// check auth
	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}

	//
	if loginQuery != "" {
		userRequest.Login = loginQuery
		userRequest.Password = "fake_password"
		if err := userRequest.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"invalid-user-login": err.Error()})
			return
		}

		domainUser, err := h.userService.GetUserByLogin(c, profileQuery, loginQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{errorGetUsers: err.Error()})
			return
		}
		// auth login
		if userCtx.Login() != domainUser.Login() && userCtx.UserType() != config.AdminRole {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"invalid user login or role": domain.ErrAccessDenied.Error()},
			)
			return
		}

		_, ok := h.stateService.GetState(c, domainUser.ID())
		status := "offline"
		if ok {
			status = "online"
		}

		response := toResponseUser(domainUser, status)
		c.JSON(http.StatusOK, response)
		return
	}
	//
	if idQuery != "" {
		userID, err := strconv.Atoi(idQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"invalid-user-id": domain.ErrValidation.Error()})
			return
		}
		if userID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id lower or equal zero"})
			return
		}

		// auth user id
		if userCtx.ID() != userID && userCtx.UserType() != config.AdminRole {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"invalid user id or role": domain.ErrAccessDenied.Error()},
			)
			return
		}
		user, err := h.userService.GetUserByID(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{errorGetUsers: err.Error()})
			return
		}
		_, ok := h.stateService.GetState(c, user.ID())
		status := "offline"
		if ok {
			status = "online"
		}
		response := toResponseUser(user, status)
		c.JSON(http.StatusOK, response)
		return
	}
}

// GetUsers is ...
// GetUsersTags 		godoc
// @Summary			Получить список всех пользователей.
// @Description		Return users list.
// @Tags			User
// @Security		BearerAuth
// @Param        	limit  query   string  true  "limit records on page" example(10) default(10)
// @Param       	offset  query   string  true  "start of record output" example(0) default(0)
// @Produce      	json
// @Success			200 {object} []UserResponse
// @failure			404 {string} err.Error()
// @Router			/admin/users [get]
func (h HTTPServer) GetUsers(c *gin.Context) {
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

	users, err := h.userService.GetUsers(c, userCtx, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get users": err.Error()})
		return
	}

	response := make([]UserResponse, 0, len(users))
	for _, user := range users {
		_, ok := h.stateService.GetState(c, user.ID())
		status := "offline"
		if ok {
			status = "online"
		}
		response = append(response, toResponseUser(user, status))
	}

	c.JSON(http.StatusOK, response)
}

// GetUsers is ...
// GetUsersTags 		godoc
// @Summary			Получить список всех пользователей.
// @Description		Return users list.
// @Tags			User
// @Security		BearerAuth
// @Param        	limit  query   string  true  "limit records on page" example(10) default(10)
// @Param       	offset  query   string  true  "start of record output" example(0) default(0)
// @Produce      	json
// @Success			200 {object} []UserResponse
// @failure			404 {string} err.Error()
// @Router			/admin/users [get]
func (h HTTPServer) AddContact(c *gin.Context) {
	var addContactRequest AddContactRequest

	if err := c.ShouldBindJSON(&addContactRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err := addContactRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{invaldRequest: err.Error()})
		return
	}

	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}

	err := h.userService.AddContact(c, userCtx, addContactRequest.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error add contact": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "contact added"})
}
