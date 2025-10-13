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
	URL           string
}

func (h *LinkHandler) HandleGet(c *gin.Context) {
	alias := c.Param("id")
	if alias == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	link, err := h.LinkService.GetLinkByAlias(alias)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, link.URL)
}

func (h *LinkHandler) HandlePost(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect request body"})
		return
	}

	url := string(body)
	link, err := h.LinkService.GetLinkByURL(url)
	if err != nil && errors.Is(err, repository.ErrNotFoundError) {
		link = &model.Link{URL: url, Alias: h.RandomService.GetRandomString(8)}
		err = h.LinkService.Add(*link)
		if err != nil && !errors.Is(err, repository.ErrExistsError) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	resp := h.URL + "/" + link.Alias
	c.String(http.StatusCreated, resp)
}
