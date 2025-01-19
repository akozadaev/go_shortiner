package shorten

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_shurtiner/internal/http/handler"
	"go_shurtiner/internal/http/helper"
	shortenRepository "go_shurtiner/internal/shorten/datebase"
	"go_shurtiner/internal/shorten/model"
	"go_shurtiner/pkg/config"
	"go_shurtiner/pkg/logging"
	"io"
	"net/http"
)

type Handler struct {
	hortenRepository shortenRepository.ShortenRepository
	cfg              *config.Config
}

func NewHandler(hortenRepository shortenRepository.ShortenRepository,
	cfg *config.Config) *Handler {
	return &Handler{
		hortenRepository: hortenRepository,
		cfg:              cfg,
	}
}

func (h *Handler) getLink(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("getting link", "get")
		return handler.NewInternalErrorResponse(errors.New("sdsd"))
		/*		return handler.NewSuccessResponse(
				http.StatusOK,
				NewLinkResponse(h.cfg.ServerConfig. ),
			)*/
	})
}

func (h *Handler) createLink(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("creating link", "create")
		//ctx := c.Request.Context()

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
		fmt.Println(links)
		return handler.NewSuccessResponse(
			http.StatusOK,
			NewLinkResponse(links, h.cfg.ServerConfig.Host),
		)
	})
}

func RouteV1(cfg *config.Config, h *Handler, r *gin.Engine) {
	v1 := r.Group("v1")

	{
		v1.GET("/short", h.getLink)
		v1.POST("/short", h.createLink)
	}

}
