package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	version string
}

func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{version: version}
}

// Health godoc
// @Summary      Chequeo de salud
// @Description  Devuelve el estado de la API y su versión
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"version":   h.version,
		"timestamp": time.Now().UTC(),
		"message":   "API is ready",
	})
}
