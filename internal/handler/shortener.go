package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/barysh-vn/shortener/internal/model"
	"github.com/barysh-vn/shortener/internal/model/api"
	"github.com/barysh-vn/shortener/internal/repository"
	"github.com/barysh-vn/shortener/internal/service"
	"github.com/gin-gonic/gin"
)

type LinkHandler struct {
	LinkService   *service.LinkService
	RandomService *service.RandomService
	URL           string
	DB            *sql.DB
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
	body, err := h.getBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.getLink(string(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusCreated, h.getURL(link))
}

func (h *LinkHandler) HandleAPIShorten(c *gin.Context) {
	body, err := h.getBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request api.ShortenRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect request body"})
		return
	}

	link, err := h.getLink(request.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, api.ShortenResponse{
		Result: h.getURL(link),
	})
}

func (h *LinkHandler) HandlePingDB(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := h.DB.PingContext(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *LinkHandler) getBody(c *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		return []byte(""), errors.New("incorrect request body")
	}

	return body, nil
}

func (h *LinkHandler) getLink(url string) (*model.Link, error) {
	link, err := h.LinkService.GetLinkByURL(url)
	if err != nil {
		if errors.Is(err, repository.ErrNotFoundError) {
			link = &model.Link{URL: url, Alias: h.RandomService.GetRandomString(8)}
			err = h.LinkService.Add(*link)
			if err != nil && !errors.Is(err, repository.ErrExistsError) {
				return link, err
			}
		} else {
			return link, err
		}
	}

	return link, nil
}

func (h *LinkHandler) getURL(link *model.Link) string {
	result, err := url.JoinPath(h.URL, link.Alias)
	if err != nil {
		return ""
	}
	return result
}
