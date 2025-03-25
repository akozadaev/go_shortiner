package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go_shurtiner/internal/app/authentication"
	"go_shurtiner/internal/app/model"
	repository "go_shurtiner/internal/app/repository"
	"go_shurtiner/internal/http/handler"
	"go_shurtiner/internal/http/helper"
	"go_shurtiner/internal/http/middleware"
	"go_shurtiner/pkg/config"
	"go_shurtiner/pkg/logging"
	trace "go_shurtiner/pkg/trace"
	"io"
	"net/http"
)

type Handler struct {
	shortenRepository repository.ShortenRepository
	userRepository    repository.UserRepository
	cfg               *config.Config
}

func NewHandler(
	hortenRepository repository.ShortenRepository,
	userRepository repository.UserRepository,
	cfg *config.Config) *Handler {
	return &Handler{
		shortenRepository: hortenRepository,
		userRepository:    userRepository,
		cfg:               cfg,
	}
}

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

func RouteV1(cfg *config.Config, h *Handler, r *gin.Engine, auth authentication.Authentication) {
	v1 := r.Group("v1")

	client, err := trace.NewTraceClient()
	authMiddleware := middleware.AuthenticationMiddleware(auth)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	v1.Use(client.MiddleWareTrace())
	{
		v1.GET("/short/:link", h.getLink)
		v1.POST("/short", h.createLink)
		v1.Use(authMiddleware).GET("/user/:id", h.getUser)
		v1.Use(authMiddleware).GET("/users", h.getUsers)
	}
}

func RouteV2(cfg *config.Config, h *Handler, r *gin.Engine, auth authentication.Authentication) {
	v2 := r.Group("v2")

	client, err := trace.NewTraceClient()
	authMiddleware := middleware.AuthenticationMiddleware(auth)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	v2.Use(client.MiddleWareTrace())
	{
		v2.Use(authMiddleware).POST("/short", h.createLinkByUser)
		v2.Use(authMiddleware).GET("/user/:id", h.getUser)
		v2.Use(authMiddleware).GET("/users", h.getUsers)
	}

}
