package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/repository"
	"github.com/barysh-vn/shortener/internal/service"
	"github.com/gin-gonic/gin"
)

type LinkHandler struct {
	LinkService   *service.LinkService
	RandomService *service.RandomService
}

func (h *LinkHandler) HandleGet(c *gin.Context) {
	alias := c.Param("id")
	if alias == "" {
		c.String(http.StatusBadRequest, "id is required")
		return
	}

	link, err := h.LinkService.GetLinkByAlias(alias)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, link.Url)
}

func (h *LinkHandler) HandlePost(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.String(http.StatusBadRequest, "Incorrect request body")
		return
	}

	url := string(body)
	link, err := h.LinkService.GetLinkByUrl(url)
	if err != nil && errors.Is(err, repository.ErrNotFoundError) {
		link = model.Link{Url: url, Alias: h.RandomService.GetRandomString(8)}
		err = h.LinkService.Add(link)
		if err != nil && !errors.Is(err, repository.ErrExistsError) {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	resp := "http://localhost:8080/" + link.Alias
	c.String(http.StatusCreated, resp)
}
