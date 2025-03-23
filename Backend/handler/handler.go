package handler

import (
	"elasticsearch/elasticsearch"
	"elasticsearch/models"
)

type ElasticSvcImpl struct {
	repo elasticsearch.ElasticSearchRepo
}

func NewElasticSvcUsecase(repo elasticsearch.ElasticSearchRepo)models.ElasticSrvUseCase{
	return &ElasticSvcImpl{
		repo: repo,
	}
}