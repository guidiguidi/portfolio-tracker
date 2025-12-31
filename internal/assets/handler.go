package assets

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
	log  *slog.Logger
}

func NewHandler(repo Repository, log *slog.Logger) *Handler {
	return &Handler{repo: repo, log: log.With(slog.String("component", "assets_handler"))}
}

type createAssetRequest struct {
	Symbol string `json:"symbol" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

// POST /assets
func (h *Handler) CreateAsset(c *gin.Context) {
	var req createAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid request format", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	a := Asset{
		Symbol: req.Symbol,
		Name:   req.Name,
	}

	created, err := h.repo.Create(c.Request.Context(), a)
	if err != nil {
		h.log.Error("failed to create asset", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.log.Info("asset created", slog.Int64("id", created.ID))
	c.JSON(http.StatusCreated, created)
}

// GET /assets/:id
func (h *Handler) GetAsset(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	a, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == ErrNotFound {
			h.log.Warn("asset not found", slog.Int64("id", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
			return
		}
		h.log.Error("failed to get asset", "error", err, slog.Int64("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, a)
}

// GET /assets
func (h *Handler) ListAssets(c *gin.Context) {
	assets, err := h.repo.List(c.Request.Context())
	if err != nil {
		h.log.Error("failed to list assets", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, assets)
}
