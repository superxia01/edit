package handler

import (
	"github.com/keenchase/edit-business/internal/service"
	"github.com/gin-gonic/gin"
)

// StatsHandler 统计数据处理器
type StatsHandler struct {
	statsService *service.StatsService
}

// NewStatsHandler 创建统计数据处理器实例
func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// GetStats 获取统计数据
// @Summary 获取统计数据
// @Description 获取笔记和博主的统计数据
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/v1/stats [get]
func (h *StatsHandler) GetStats(c *gin.Context) {
	stats, err := h.statsService.GetStats()
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessResponse(c, stats)
}
