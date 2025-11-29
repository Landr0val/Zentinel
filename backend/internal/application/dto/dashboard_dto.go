package dto

import "zentinel/internal/domain/enums"

type DashboardStatsResponse struct {
	TotalTransactions int64                       `json:"total_transactions"`
	TotalVolume       float64                     `json:"total_volume"`
	FlaggedCount      int64                       `json:"flagged_count"`
	FlaggedVolume     float64                     `json:"flagged_volume"`
	AlertsByStatus    map[enums.AlertStatus]int64 `json:"alerts_by_status"`
}
