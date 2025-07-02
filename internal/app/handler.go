package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	// docs "github.com/go_shurtiner/docs"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docsV1 "go_shurtiner/docs/v1"
	docsV2 "go_shurtiner/docs/v2"
	"go_shurtiner/internal/app/authentication"
	"go_shurtiner/internal/app/model"
	repository "go_shurtiner/internal/app/repository"
	"go_shurtiner/internal/http/handler"
	"go_shurtiner/internal/http/helper"
	"go_shurtiner/internal/http/middleware"
	"go_shurtiner/internal/queue/service"
	"go_shurtiner/pkg/config"
	"go_shurtiner/pkg/logging"
	trace "go_shurtiner/pkg/trace"

	"io"
	"net/http"
	// "time"
)

type Handler struct {
	shortenRepository  repository.ShortenRepository
	userRepository     repository.UserRepository
	queueJobRepository service.QueueRepository
	cfg                *config.Config
}

func NewHandler(
	shortenRepository repository.ShortenRepository,
	userRepository repository.UserRepository,
	queueJobRepository service.QueueRepository,
	cfg *config.Config) *Handler {
	return &Handler{
		shortenRepository:  shortenRepository,
		userRepository:     userRepository,
		queueJobRepository: queueJobRepository,
		cfg:                cfg,
	}
}

// getLink godoc
// @Summary Получить информацию по сокращённой ссылке
// @Description Поиск оригинальной ссылки по короткому идентификатору.
// @Tags link
// @Accept json
// @Produce json
// @Param link path string true "Идентификатор сокращённой ссылки"
// @Success 200 {object} model.Link "Успешный ответ с полной информацией о ссылке"
// @Failure 400 {object} handler.ErrorResponse "Внутренняя ошибка сервера"
// @Router /short/{link} [get]
func (h *Handler) getLink(c *gin.Context) {
	handler.HandleRequest(c, h.cfg.ServerConfig.GoroutineTimeout, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("getting link", "get")
		linkStr := c.Param("link")
		res, err := h.shortenRepository.FindLink(c.Request.Context(), linkStr)
		if err != nil {
			return handler.NewInternalErrorResponse(err)
		}

		res.Shortened = fmt.Sprintf("%s%s", h.cfg.ServerConfig.Host, res.Shortened)
		return handler.NewSuccessResponse(
			http.StatusOK,
			res,
		)
	})
}

func (h *Handler) getLinkByUser(c *gin.Context) {
	handler.HandleRequest(c, h.cfg.ServerConfig.GoroutineTimeout, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("getting link", "get")
		linkStr := c.Param("link")
		res, err := h.shortenRepository.FindLink(c.Request.Context(), linkStr)
		if err != nil {
			return handler.NewInternalErrorResponse(err)
		}

		res.Shortened = fmt.Sprintf("%s%s", h.cfg.ServerConfig.Host, res.Shortened)
		return handler.NewSuccessResponse(
			http.StatusOK,
			res,
		)
	})
}

func (h *Handler) getLinks(c *gin.Context) {
	handler.HandleRequest(c, h.cfg.ServerConfig.GoroutineTimeout, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("getting links", "get")
		res, err := h.shortenRepository.FetchLinks(c.Request.Context())
		if err != nil {
			return handler.NewInternalErrorResponse(err)
		}

		if err != nil {
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(
			http.StatusOK,
			res,
		)
	})
}

// getUser godoc
// @Summary Получить информацию по зарегистрированным пользователям.
// @Description Информация по всем зарегистрированным пользователям.
// @Tags user
// @Accept json
// @Produce json
// @Param offset query int false "Смещение пагинации"
// @Success 200 {object} []model.UserApi "Успешный ответ с полной информацией о пользователях"
// @Failure 400 {object} handler.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users [get]
func (h *Handler) getUsers(c *gin.Context) {
	handler.HandleRequest(c, h.cfg.ServerConfig.GoroutineTimeout, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("getting users", "getUsers")

		res, err := h.userRepository.FetchUsers(c.Request.Context())
		if err != nil {
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(
			http.StatusOK,
			res,
		)
	})
}

// getUser godoc
// @Summary Получить информацию по зарегистрированному пользователю.
// @Description Информация по зарегистрированному пользователю.
// @Tags user
// @Accept json
// @Produce json
// @Param offset query int false "Смещение пагинации"
// @Success 200 {object} model.UserApi "Успешный ответ с полной информацией о пользователях"
// @Failure 400 {object} handler.ErrorResponse "Внутренняя ошибка сервера"
// @Router /user/{id} [get]
func (h *Handler) getUser(c *gin.Context) {
	handler.HandleRequest(c, h.cfg.ServerConfig.GoroutineTimeout, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("getting user by id", "getUser")
		idStr := c.Param("id")
		res, err := h.userRepository.GetUserForApiById(c.Request.Context(), idStr)
		if err != nil {
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(
			http.StatusOK,
			res,
		)
	})
}

// getLink godoc
// @Summary Сократить ссылки
// @Description Сокращение ссылок из массива параметров.
// @Tags link
// @Accept json
// @Produce json
// @Param body body []model.CreateLink true "Тело запроса"
// @Success 200 {object} []model.Link "Успешный ответ с полной информацией о добавленных ссылках"
// @Failure 400 {object} handler.ErrorResponse "Внутренняя ошибка сервера"
// @Router /short [post]
func (h *Handler) createLink(c *gin.Context) {
	handler.HandleRequest(c, h.cfg.ServerConfig.GoroutineTimeout, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("creating link", "create")

		links := make([]model.CreateLink, 0)
		body, err := io.ReadAll(c.Request.Body)
		if !helper.RequestHasJsonArray(body) {
			var link model.CreateLink
			err = json.Unmarshal(body, &link)
			if err != nil {
				_ = c.Error(err)
				c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
				return handler.NewInternalErrorResponse(err)
			}
			links = append(links, link)
		} else {
			err = json.Unmarshal(body, &links)
			if err != nil {
				_ = c.Error(err)
				c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
				return handler.NewInternalErrorResponse(err)
			}
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return handler.NewInternalErrorResponse(err)
		}

		linkResponse := NewLinkResponse(links, h.cfg.ServerConfig.Host)
		for _, datum := range linkResponse.Data {
			link := model.Link{
				Source:    datum.Source,
				Shortened: datum.Shortened,
			}
			err = h.shortenRepository.SaveLink(c.Request.Context(), &link)
			if err != nil {
				return handler.NewInternalErrorResponse(err)
			}
		}

		return handler.NewSuccessResponse(
			http.StatusOK,
			linkResponse,
		)
	})
}

// getLink godoc
// @Summary Сократить ссылки от имени зарегистрированного пользователя
// @Description Сокращение ссылок из массива параметров.
// @Tags link
// @Accept json
// @Produce json
// @Param body body []model.CreateLink true "Тело запроса"
// @Success 200 {object} []model.Link "Успешный ответ с полной информацией о добавленных ссылках"
// @Failure 400 {object} handler.ErrorResponse "Внутренняя ошибка сервера"
// @Router /short [post]
func (h *Handler) createLinkByUser(c *gin.Context) {
	handler.HandleRequest(c, h.cfg.ServerConfig.GoroutineTimeout, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("creating link", "create")

		user, err := authentication.GetUser(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "")
			return handler.NewInternalErrorResponse(err)
		}

		links := make([]model.CreateLink, 0)
		body, _ := io.ReadAll(c.Request.Body)
		if !helper.RequestHasJsonArray(body) {
			var link model.CreateLink
			err = json.Unmarshal(body, &link)
			if err != nil {
				_ = c.Error(err)
				c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
				return handler.NewInternalErrorResponse(err)
			}
			links = append(links, link)
		} else {
			err = json.Unmarshal(body, &links)
			if err != nil {
				_ = c.Error(err)
				c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
				return handler.NewInternalErrorResponse(err)
			}
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return handler.NewInternalErrorResponse(err)
		}

		linkResponse := NewLinkResponse(links, h.cfg.ServerConfig.Host)

		userArray := make([]*model.User, 0)
		userArray = append(userArray, user)

		for _, datum := range linkResponse.Data {
			link := model.Link{
				Source:    datum.Source,
				Shortened: datum.Shortened,
				User:      userArray,
			}
			err = h.shortenRepository.SaveLink(c.Request.Context(), &link)
			if err != nil {
				return handler.NewInternalErrorResponse(err)
			}
		}

		return handler.NewSuccessResponse(
			http.StatusOK,
			linkResponse,
		)
	})
}

// getLink godoc
// @Summary Создать отложенную задачу генерации отчета.
// @Description Постановка задачи генерации отчета.
// @Tags link
// @Accept json
// @Produce json
// @Param body body model.Params true "Тело запроса"
// @Success 200 {object} []model.JobQueue "Успешный ответ с полной информацией о добавленных ссылках"
// @Failure 400 {object} handler.ErrorResponse "Внутренняя ошибка сервера"
// @Router /report [post]
func (h *Handler) setBackgroundTaskAsJob(c *gin.Context) {
	handler.HandleRequest(c, h.cfg.ServerConfig.GoroutineTimeout, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("creating background task", "create")

		_, err := authentication.GetUser(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "")
			return handler.NewInternalErrorResponse(err)
		}

		queueJobs := make([]model.JobQueue, 0)
		body, _ := io.ReadAll(c.Request.Body)
		if !helper.RequestHasJsonArray(body) {
			var queueJob model.JobQueue
			err = json.Unmarshal(body, &queueJob)
			if err != nil {
				_ = c.Error(err)
				c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
				return handler.NewInternalErrorResponse(err)
			}

			err = h.queueJobRepository.CreateJob(c, &queueJob)
			if err != nil {
				_ = c.Error(err)
				c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
				return handler.NewInternalErrorResponse(err)
			}
			queueJobs = append(queueJobs, queueJob)
		}

		return handler.NewSuccessResponse(
			http.StatusOK,
			queueJobs,
		)
	})
}

func RouteV1(cfg *config.Config, h *Handler, r *gin.Engine, auth authentication.Authentication) {
	v1 := r.Group("v1")

	client, err := trace.NewTraceClient()
	authMiddleware := middleware.AuthenticationMiddleware(auth)
	if err != nil {
		log.Error().Stack().Err(err)
	}

	docsV1.SwaggerInfov1.Schemes = []string{"http"}
	docsV1.SwaggerInfov1.BasePath = "/v1"
	docsV2.SwaggerInfov2.Schemes = []string{"http"}
	docsV2.SwaggerInfov2.BasePath = "/v2"

	r.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.InstanceName("v1")))
	r.GET("/swagger/v2/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.InstanceName("v2")))

	v1.Use(client.MiddleWareTrace())
	{

		v1.GET("/short/:link", h.getLink)
		v1.POST("/short", h.createLink)
		v1.Use(authMiddleware).GET("/user/:id", h.getUser)
		v1.Use(authMiddleware).GET("/users", h.getUsers)
	}
}

func RouteV2(cfg *config.Config, h *Handler, r *gin.Engine, auth authentication.Authentication) {
	authMiddleware := middleware.AuthenticationMiddleware(auth)
	client, err := trace.NewTraceClient()
	if err != nil {
		log.Error().Stack().Err(err)
	}
	v2 := r.Group("v2")
	v2.Use(authMiddleware)
	v2.Use(client.MiddleWareTrace())

	{
		// docs.SwaggerInfo.BasePath = "/v2"
		v2.POST("/short", h.createLinkByUser)
		v2.GET("/user/:id", h.getUser)
		v2.GET("/users", h.getUsers)
		v2.POST("/report", h.setBackgroundTaskAsJob)
	}

}
