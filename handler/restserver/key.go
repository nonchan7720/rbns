package restserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/api-rbac/domain/repository"
)

type apiKeyHandle interface {
	generate(*gin.Context)
	get(*gin.Context)
}

type apiKeyHandler struct {
	repo repository.ApiKey
}

func newApiKeyHander(repo repository.ApiKey) apiKeyHandle {
	return &apiKeyHandler{
		repo: repo,
	}
}

func (h *apiKeyHandler) generate(c *gin.Context) {
	token, err := h.repo.Generate()
	if err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func (h *apiKeyHandler) get(c *gin.Context) {
	token := h.repo.Get()
	c.JSON(http.StatusOK, gin.H{"token": token})
}
