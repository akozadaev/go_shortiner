package shorten

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go_shurtiner/internal/http/handler"
	"go_shurtiner/internal/http/helper"
	shortenRepository "go_shurtiner/internal/shorten/datebase"
	"go_shurtiner/internal/shorten/model"
	"go_shurtiner/pkg/config"
	"go_shurtiner/pkg/logging"
	trace "go_shurtiner/pkg/trace"
	"io"
	"net/http"
)

type Handler struct {
	shortenRepository shortenRepository.ShortenRepository
	cfg               *config.Config
}

func NewHandler(hortenRepository shortenRepository.ShortenRepository,
	cfg *config.Config) *Handler {
	return &Handler{
		shortenRepository: hortenRepository,
		cfg:               cfg,
	}
}

func (h *Handler) getLink(c *gin.Context) {
	handler.HandleRequest(c, h.cfg.ServerConfig.GoroutineTimeout, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("getting link", "get")
		linkStr := c.Param("link")
		fmt.Println(linkStr)
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

func RouteV1(cfg *config.Config, h *Handler, r *gin.Engine) {
	v1 := r.Group("v1")

	client, err := trace.NewTraceClient()
	if err != nil {
		log.Error().Stack().Err(err)
	}
	v1.Use(client.MiddleWareTrace())

	{
		v1.GET("/short/:link", h.getLink)
		v1.POST("/short", h.createLink)
	}

}
