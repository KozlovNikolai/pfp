package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	errorGetAccount = "error get account"
)

// func (h HTTPServer) GetAccountByUser(c *gin.Context) {
// 	var accountRequest AccountRequest
// 	accountNameQuery := c.Query("name")
// 	accountIDQuery := c.Query("id")

// 	// check auth
// 	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
// 		return
// 	}

// 	//
// 	if accountIDQuery != "" {
// 		accountRequest.Login = loginQuery
// 		accountRequest.Password = "fake_password"
// 		if err := accountRequest.Validate(); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"invalid-user-login": err.Error()})
// 			return
// 		}

// 		domainUser, err := h.userService.GetUserByLogin(c, accountNameQuery, loginQuery)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{errorGetUsers: err.Error()})
// 			return
// 		}
// 		// auth login
// 		if userCtx.Login() != domainUser.Login() && userCtx.UserType() != config.AdminRole {
// 			c.JSON(
// 				http.StatusUnauthorized,
// 				gin.H{"invalid user login or role": domain.ErrAccessDenied.Error()},
// 			)
// 			return
// 		}

// 		_, ok := h.stateService.GetState(c, domainUser.ID())
// 		status := "offline"
// 		if ok {
// 			status = "online"
// 		}

// 		response := toResponseUser(domainUser, status)
// 		c.JSON(http.StatusOK, response)
// 		return
// 	}
// 	//
// 	if accountNameQuery != "" {
// 		accountID, err := strconv.Atoi(accountIDQuery)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"invalid-account-id": domain.ErrValidation.Error()})
// 			return
// 		}
// 		if accountID <= 0 {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "id lower or equal zero"})
// 			return
// 		}

// 		// auth user id
// 		if userCtx.ID() != accountID && userCtx.UserType() != config.AdminRole {
// 			c.JSON(
// 				http.StatusUnauthorized,
// 				gin.H{"invalid user id or role": domain.ErrAccessDenied.Error()},
// 			)
// 			return
// 		}
// 		user, err := h.userService.GetUserByID(c, accountID)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{errorGetUsers: err.Error()})
// 			return
// 		}
// 		_, ok := h.stateService.GetState(c, user.ID())
// 		status := "offline"
// 		if ok {
// 			status = "online"
// 		}
// 		response := toResponseUser(user, status)
// 		c.JSON(http.StatusOK, response)
// 		return
// 	}
// }

// func (h HTTPServer) GetUserByAccount(c *gin.Context) {
// 	var accountRequest AccountRequest
// 	accountNameQuery := c.Query("name")
// 	accountIDQuery := c.Query("id")

// 	// check auth
// 	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
// 		return
// 	}

// 	//
// 	if accountIDQuery != "" {
// 		accountRequest.Login = loginQuery
// 		accountRequest.Password = "fake_password"
// 		if err := accountRequest.Validate(); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"invalid-user-login": err.Error()})
// 			return
// 		}

// 		domainUser, err := h.userService.GetUserByLogin(c, accountNameQuery, loginQuery)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{errorGetUsers: err.Error()})
// 			return
// 		}
// 		// auth login
// 		if userCtx.Login() != domainUser.Login() && userCtx.UserType() != config.AdminRole {
// 			c.JSON(
// 				http.StatusUnauthorized,
// 				gin.H{"invalid user login or role": domain.ErrAccessDenied.Error()},
// 			)
// 			return
// 		}

// 		_, ok := h.stateService.GetState(c, domainUser.ID())
// 		status := "offline"
// 		if ok {
// 			status = "online"
// 		}

// 		response := toResponseUser(domainUser, status)
// 		c.JSON(http.StatusOK, response)
// 		return
// 	}
// 	//
// 	if accountNameQuery != "" {
// 		accountID, err := strconv.Atoi(accountIDQuery)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"invalid-account-id": domain.ErrValidation.Error()})
// 			return
// 		}
// 		if accountID <= 0 {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "id lower or equal zero"})
// 			return
// 		}

// 		// auth user id
// 		if userCtx.ID() != accountID && userCtx.UserType() != config.AdminRole {
// 			c.JSON(
// 				http.StatusUnauthorized,
// 				gin.H{"invalid user id or role": domain.ErrAccessDenied.Error()},
// 			)
// 			return
// 		}
// 		user, err := h.userService.GetUserByID(c, accountID)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{errorGetUsers: err.Error()})
// 			return
// 		}
// 		_, ok := h.stateService.GetState(c, user.ID())
// 		status := "offline"
// 		if ok {
// 			status = "online"
// 		}
// 		response := toResponseUser(user, status)
// 		c.JSON(http.StatusOK, response)
// 		return
// 	}
// }

func (h HTTPServer) CreateAccount(c *gin.Context) {
	var accountCreateRequest AccountRequest
	var err error
	if err = c.ShouldBindJSON(&accountCreateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err = accountCreateRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{invaldRequest: err.Error()})
		return
	}

	domainAccount := toDomainAccount(accountCreateRequest)

	createdAccount, err := h.accountService.CreateAccount(c, domainAccount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error service Account": err.Error()})
		return
	}
	response := toResponseAccount(createdAccount)
	c.JSON(http.StatusOK, response)
	return
}

// // GetUsers is ...
// // GetUsersTags 		godoc
// // @Summary			Получить список всех пользователей.
// // @Description		Return users list.
// // @Tags			User
// // @Security		BearerAuth
// // @Param        	limit  query   string  true  "limit records on page" example(10) default(10)
// // @Param       	offset  query   string  true  "start of record output" example(0) default(0)
// // @Produce      	json
// // @Success			200 {object} []UserResponse
// // @failure			404 {string} err.Error()
// // @Router			/admin/users [get]
// func (h HTTPServer) GetUsers(c *gin.Context) {
// 	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
// 		return
// 	}
// 	userCtx, err = h.userService.GetUserByID(c, userCtx.ID())
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNotFound.Error()})
// 		return
// 	}
// 	// fmt.Printf("\nuserCtx: %+v\n", userCtx)
// 	// check admin
// 	if userCtx.UserType() != config.AdminRole {
// 		c.JSON(
// 			http.StatusUnauthorized,
// 			gin.H{"invalid user id or role": domain.ErrAccessDenied.Error()},
// 		)
// 		return
// 	}
// 	limitQuery := c.Query("limit")
// 	offsetQuery := c.Query("offset")

// 	limit, err := strconv.Atoi(limitQuery)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"invalid-limit": err.Error()})
// 		return
// 	}

// 	offset, err := strconv.Atoi(offsetQuery)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"invalid-offset": err.Error()})
// 		return
// 	}
// 	if limit < 1 {
// 		c.JSON(http.StatusBadRequest, gin.H{"limit-must-be-greater-then-zero": ""})
// 		return
// 	}
// 	if offset < 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"offset-must-be-greater-or-equal-then-zero": ""})
// 		return
// 	}

// 	users, err := h.userService.GetUsers(c, userCtx, limit, offset)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error get users": err.Error()})
// 		return
// 	}

// 	response := make([]UserResponse, 0, len(users))
// 	for _, user := range users {
// 		_, ok := h.stateService.GetState(c, user.ID())
// 		status := "offline"
// 		if ok {
// 			status = "online"
// 		}
// 		response = append(response, toResponseUser(user, status))
// 	}

// 	c.JSON(http.StatusOK, response)
// }
