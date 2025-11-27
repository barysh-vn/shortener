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
	"go.uber.org/zap"
)

type LinkHandler struct {
	LinkService   *service.LinkService
	RandomService *service.RandomService
	URL           string
	DB            *sql.DB
	Logger        *zap.Logger
}

func (h *LinkHandler) HandleGet(c *gin.Context) {
	alias := c.Param("id")
	if alias == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	link, err := h.LinkService.GetLinkByAlias(c, alias)
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

	link, err := h.getLink(c, string(body))
	if err != nil {
		if errors.Is(err, repository.ErrExistsError) {
			c.String(http.StatusConflict, h.getURL(link))
			return
		}
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

	link, err := h.getLink(c, request.URL)
	if err != nil {
		if errors.Is(err, repository.ErrExistsError) {
			c.JSON(http.StatusConflict, api.ShortenResponse{
				Result: h.getURL(link),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, api.ShortenResponse{
		Result: h.getURL(link),
	})
}

func (h *LinkHandler) HandleBatchAPIShorten(c *gin.Context) {
	body, err := h.getBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var requests []api.ShortenBatchURLRequest
	if err = json.Unmarshal(body, &requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect request body"})
		return
	}

	var links []model.Link

	for _, req := range requests {
		link, err := h.LinkService.GetLinkByURL(c, req.URL)
		if err != nil {
			if errors.Is(err, repository.ErrNotFoundError) {
				link = &model.Link{
					URL:    req.URL,
					Alias:  h.RandomService.GetRandomString(8),
					UserID: c.GetString("user_id"),
				}
			} else {
				h.Logger.Info("Failed to get link", zap.String("url", req.URL), zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
				return
			}
		}
		links = append(links, *link)
	}

	var newLinks []model.Link
	for _, link := range links {
		if _, err = h.LinkService.GetLinkByAlias(c, link.Alias); errors.Is(err, repository.ErrNotFoundError) {
			newLinks = append(newLinks, link)
		}
	}

	if len(newLinks) > 0 {
		if err = h.LinkService.AddBatch(c, h.DB, newLinks); err != nil {
			h.Logger.Info("Failed to add batch links", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
			return
		}
	}

	var response []api.ShortenBatchURLResponse
	for i, req := range requests {
		response = append(response, api.ShortenBatchURLResponse{
			ID:  req.ID,
			URL: h.getURL(&links[i]),
		})
	}

	c.JSON(http.StatusCreated, response)
}

func (h *LinkHandler) HandlePingDB(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()
	if h.DB == nil {
		h.Logger.Info("Failed to ping db (no db connection)")
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	if err := h.DB.PingContext(ctx); err != nil {
		h.Logger.Info("Failed to ping db", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *LinkHandler) HandleUserURLs(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}

	links, err := h.LinkService.GetLinksByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	if len(*links) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"error": http.StatusText(http.StatusNoContent)})
		return
	}

	var response []api.ShortenUserURLsResponse
	for _, req := range *links {
		response = append(response, api.ShortenUserURLsResponse{
			Alias: h.getURL(&req),
			URL:   req.URL,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *LinkHandler) getBody(c *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		return []byte(""), errors.New("incorrect request body")
	}

	return body, nil
}

func (h *LinkHandler) getLink(c *gin.Context, url string) (*model.Link, error) {
	link, err := h.LinkService.GetLinkByURL(c, url)
	if err != nil {
		if errors.Is(err, repository.ErrNotFoundError) {
			link = &model.Link{URL: url, Alias: h.RandomService.GetRandomString(8), UserID: c.GetString("user_id")}
			err = h.LinkService.Add(c, *link)
			if err != nil && !errors.Is(err, repository.ErrExistsError) {
				return link, err
			}
		} else {
			return link, err
		}
	} else {
		err = repository.ErrExistsError
	}

	return link, err
}

func (h *LinkHandler) getURL(link *model.Link) string {
	result, err := url.JoinPath(h.URL, link.Alias)
	if err != nil {
		return ""
	}
	return result
}
