package shorten

import (
	"github.com/gin-gonic/gin"
	"go_shurtiner/internal/http/handler"
	shortenRepository "go_shurtiner/internal/shorten/datebase"
	"go_shurtiner/pkg/config"
	"go_shurtiner/pkg/logging"
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
		return handler.NewSuccessResponse(
			http.StatusOK,
			NewLinkResponse(h.cfg.ServerConfig.Host),
		)
	})
}

func (h *Handler) createLink(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		logger := logging.FromContext(c)
		logger.Debugw("creating link", "create")
		return handler.NewSuccessResponse(
			http.StatusOK,
			NewLinkResponse(h.cfg.ServerConfig.Host),
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
