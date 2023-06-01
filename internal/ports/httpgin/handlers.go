package httpgin

import (
	"errors"
	"github.com/TobbyMax/validator"
	"github.com/gin-gonic/gin"
	"homework10/internal/app"
	"io"
	"net/http"
	"strconv"
)

var ErrParameterNotFound = errors.New("necessary parameters not provided")

type Handler = func(*gin.Context) error

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)

		if err != nil {
			switch {
			case errors.As(err, &validator.ValidationErrors{}):
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			case errors.Is(err, app.ErrUserNotFound):
				c.JSON(http.StatusFailedDependency, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adIDStr := c.Param("ad_id")
		adID, err := strconv.Atoi(adIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.ChangeAdStatus(c, int64(adID), reqBody.UserID, reqBody.Published)

		if err != nil {
			switch {
			case errors.Is(err, app.ErrForbidden):
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
			case errors.Is(err, app.ErrAdNotFound):
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adIDStr := c.Param("ad_id")
		adID, err := strconv.Atoi(adIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.UpdateAd(c, int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)

		if err != nil {
			switch {
			case errors.As(err, &validator.ValidationErrors{}):
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			case errors.Is(err, app.ErrForbidden):
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
			case errors.Is(err, app.ErrAdNotFound):
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для получения объявления по id
func getAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adIDStr := c.Param("ad_id")
		adID, err := strconv.Atoi(adIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.GetAd(c, int64(adID))

		if err != nil {
			switch {
			case errors.Is(err, app.ErrAdNotFound):
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для удаления объявления по id
func deleteAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adIDStr := c.Param("ad_id")
		adID, err := strconv.Atoi(adIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		userIDStr, ok := c.GetQuery("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, AdErrorResponse(ErrParameterNotFound))
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		err = a.DeleteAd(c, int64(adID), int64(userID))

		if err != nil {
			switch {
			case errors.Is(err, app.ErrAdNotFound):
				c.JSON(http.StatusNotFound, AdErrorResponse(err))
			case errors.Is(err, app.ErrForbidden):
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, DeletionSuccessResponse())
	}
}

// Метод для получения списка объявлений с фильтрами
func listAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody listAdsRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil && err != io.EOF {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		date, err := app.ParseDate(reqBody.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		al, err := a.ListAds(c, app.ListAdsParams{
			Published: reqBody.Published,
			Uid:       reqBody.UserID,
			Date:      date,
			Title:     reqBody.Title,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}
		//if len(al.Data) == 0 {
		//	c.JSON(http.StatusNoContent, AdListSuccessResponse(al))
		//	return
		//}
		c.JSON(http.StatusOK, AdListSuccessResponse(al))
	}
}

// Метод для создания пользователя
func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		u, err := a.CreateUser(c, reqBody.Nickname, reqBody.Email)

		if err != nil {
			switch {
			case errors.As(err, &validator.ValidationErrors{}):
				c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(u))
	}
}

// Метод для обновления имени(Nickname) или почты(Email) пользователя
func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		userIDStr := c.Param("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		u, err := a.UpdateUser(c, int64(userID), reqBody.Nickname, reqBody.Email)

		if err != nil {
			switch {
			case errors.As(err, &validator.ValidationErrors{}):
				c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			case errors.Is(err, app.ErrUserNotFound):
				c.JSON(http.StatusNotFound, UserErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(u))
	}
}

// Метод для получения пользователя по id
func getUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		u, err := a.GetUser(c, int64(userID))

		if err != nil {
			switch {
			case errors.Is(err, app.ErrUserNotFound):
				c.JSON(http.StatusNotFound, UserErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(u))
	}
}

// Метод для удаления пользователя по id
func deleteUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		err = a.DeleteUser(c, int64(userID))

		if err != nil {
			switch {
			case errors.Is(err, app.ErrUserNotFound):
				c.JSON(http.StatusNotFound, UserErrorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			}
			return
		}
		c.JSON(http.StatusOK, DeletionSuccessResponse())
	}
}
