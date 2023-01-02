package user

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
	"message/internal/domain/user"
	"message/internal/pkg/apperror"
	"message/internal/server/http/api_error"
)

// @Summary get users
// @Description get users by params
// @ID get-users
// @Accept json
// @Produce json
// @Success 200 {object} UsersList
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_error.Error
// @Router /auth [get]
func AuthHandler(c *gin.Context, app *app.App) {
	cond := user.User{}
	if err := c.ShouldBindJSON(&cond); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := app.Domain.User.Service.First(context.Background(), &cond)
	if err != nil {
		if err == apperror.ErrNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": apperror.ErrNotFound.Error()})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal.Error()})
		return
	}

	jwtToken, err := app.Domain.AuthToken.Service.Create(context.Background(), u.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId": u.ID,
		"token":  jwtToken,
	})
}

// @Summary get users
// @Description get users by params
// @ID get-users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_error.Error
// @Router /registration [get]
func RegisterHandler(c *gin.Context, app *app.App) {
	u := user.User{}
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := app.Domain.User.Service.Create(context.Background(), &u)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal.Error()})
		return
	}

	jwtToken, err := app.Domain.AuthToken.Service.Create(context.Background(), userID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  jwtToken,
		"userId": userID,
	})
}

// @Summary get users
// @Description get users by params
// @ID get-users
// @Accept json
// @Produce json
// @Param with_organization query boolean false "with organization"
// @Param per_page query int false "per page"
// @Param page query int false "page"
// @Success 200 {object} UsersList
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_wwerror.Error
// @Router /users [get]
func GetUsersHandler(c *gin.Context, app *app.App) {
	user := user.User{}
	if err := c.ShouldBindQuery(&user); err != nil {
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

	users, err := app.Domain.User.Service.Query(context.Background(), &user, &pag)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, api_error.New(apperror.ErrInternal))
		return
	}

	count, err := app.Domain.User.Service.Count(context.Background(), &user)
	if err != nil {
		log.Println(err)

		c.JSON(http.StatusInternalServerError, api_error.New(apperror.ErrInternal))
		return
	}

	c.JSON(http.StatusOK, UsersList{
		Items: users,
		Count: count,
	})
}

// @Summary get user by id
// @Description get user by id
// @ID get-user-by-id
// @Accept json
// @Produce json
// @Param id path boolean true "id"
// @Success 200 {object} User
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_error.Error
// @Router /user/{id} [get]
func GetUserHandler(c *gin.Context, app *app.App) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, api_error.New(err))
		return
	}

	u, err := app.Domain.User.Service.Get(context.Background(), uint(userID))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, api_error.New(apperror.ErrInternal))
		return
	}

	c.JSON(http.StatusOK, u)
}

// @Summary update user
// @Description update user
// @ID update-user
// @Accept json
// @Produce json
// @Param user body user.User true "updatable user"
// @Success 200 {string} success
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_error.Error
// @Router /user [put]
func UpdateUserHandler(c *gin.Context, app *app.App) {
	bdg := user.User{}
	if err := c.ShouldBindJSON(&bdg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := app.Domain.User.Service.Update(context.Background(), &bdg)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary create user
// @Description create user
// @ID create-user
// @Accept json
// @Produce json
// @Param user body user.User true "creatable user"
// @Success 200 {object} User
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_error.Error
// @Router /user [post]
func CreateUserHandler(c *gin.Context, app *app.App) {
	bdg := user.User{}
	if err := c.ShouldBindJSON(&bdg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := app.Domain.User.Service.Create(context.Background(), &bdg)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal.Error()})
		return
	}
	bdg.ID = userID

	c.JSON(http.StatusOK, bdg)
}

// @Summary delete user by id
// @Description delete user by id
// @ID delete-user-by-id
// @Accept json
// @Produce json
// @Param id path boolean true "id"
// @Success 200 {string} success
// @Failure 400 {object} api_error.Error
// @Failure 500 {object} api_error.Error
// @Router /user/{id} [delete]
func DeleteUserHandler(c *gin.Context, app *app.App) {
	userID, err := selection_condition.ParseUintParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.Domain.User.Service.Delete(context.Background(), userID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
