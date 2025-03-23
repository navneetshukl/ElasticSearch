package elasticsearch

import (
	"elasticsearch/models"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	indexName string = "USERS"
)

type ElasticRepositoryImpl struct {
	client *elasticsearch.Client
}

func ConnectToElastic() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("error in connecting to elastic search ", err)
		return nil, err
	}

	return client, nil
}

type ElasticSearchRepo interface {
	InsertUser(user models.User) error
	SearchUser(query string) ([]models.User, error)
}

func NewElasticRepository(cl *elasticsearch.Client) *ElasticRepositoryImpl {
	return &ElasticRepositoryImpl{
		client: cl,
	}
}
