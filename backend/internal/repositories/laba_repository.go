package repository

import (
	"context"
	"laba_service/config"
	"laba_service/internal/models"
)

type LabaRepository interface {
	GetAll(ctx context.Context) ([]models.Laba, error)
	Create(ctx context.Context, labas []models.Laba) error
}

type labaRepository struct {
	db config.GormPostgres
}

func NewLabaRepository(db config.GormPostgres) LabaRepository {
	return &labaRepository{db: db}
}

func (r *labaRepository) GetAll(ctx context.Context) ([]models.Laba, error) {
	conn := r.db.GetConnection()
	var labas []models.Laba
	err := conn.WithContext(ctx).Find(&labas).Error
	return labas, err
}

func (r *labaRepository) Create(ctx context.Context, labas []models.Laba) error {
	conn := r.db.GetConnection()
	return conn.WithContext(ctx).Create(&labas).Error
}
