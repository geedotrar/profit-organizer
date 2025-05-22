package service

import (
	"context"
	"laba_service/internal/models"
	repository "laba_service/internal/repositories"
	"sort"
	"time"
)

type LabaService interface {
	GetAll(ctx context.Context) (map[string]map[string]float64, error)
	Create(ctx context.Context, labas []models.Laba) error
}

type labaService struct {
	repo repository.LabaRepository
}

func NewLabaService(repo repository.LabaRepository) LabaService {
	return &labaService{repo: repo}
}

func (s *labaService) GetAll(ctx context.Context) (map[string]map[string]float64, error) {
	rawData, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(rawData, func(i, j int) bool {
		return rawData[i].CreatedAt.After(rawData[j].CreatedAt)
	})

	result := make(map[string]map[string]float64)
	latest := make(map[string]map[string]time.Time)

	for _, item := range rawData {
		label := item.LabelRekonsiliasiFiskal
		periodeStr := item.Periode.Format("2006-01-02")

		if _, ok := result[label]; !ok {
			result[label] = make(map[string]float64)
			latest[label] = make(map[string]time.Time)
		}

		if prev, exists := latest[label][periodeStr]; !exists || item.CreatedAt.After(prev) {
			result[label][periodeStr] = item.Nilai
			latest[label][periodeStr] = item.CreatedAt
		}
	}

	return result, nil
}

func (s *labaService) Create(ctx context.Context, labas []models.Laba) error {
	return s.repo.Create(ctx, labas)
}
