package models

import "time"

// ExchangeRate 汇率模型
// 用于存储货币兑换汇率信息
type ExchangeRate struct {
	ID           uint      `gorm:"primarykey" json:"_id"` // 主键
	FromCurrency string    `json:"fromCurrency" binding:"required"` // 源货币代码，如 USD
	ToCurrency   string    `json:"toCurrency" binding:"required"` // 目标货币代码，如 CNY
	Rate         float64   `json:"rate" binding:"required"` // 汇率值
	Date         time.Time `json:"date"` // 汇率日期
}
