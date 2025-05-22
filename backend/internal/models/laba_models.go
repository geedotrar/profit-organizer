package models

import (
	"time"
)

type Laba struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LabelRekonsiliasiFiskal string    `gorm:"size:100;not null;index:idx_label" json:"label_rekonsiliasi_fiskal"`
	Periode                 time.Time `gorm:"type:timestamptz;not null;index:idx_periode" json:"periode"`
	Nilai                   float64   `gorm:"type:numeric(30,5);not null;default:0;index:idx_nilai" json:"nilai"`
	CreatedAt               time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Laba) TableName() string {
	return "branch_laba_sebelum_pajak_penghasilan_tax"
}
