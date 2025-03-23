package elasticsearch

import (
	"context"
	"elasticsearch/models"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func (r *ElasticRepositoryImpl) createIndexIfNotExists() error {
	mapping := `{
        "settings": {
            "analysis": {
                "analyzer": {
                    "substring_analyzer": {
                        "tokenizer": "ngram_tokenizer",
                        "filter": ["lowercase"]
                    }
                },
                "tokenizer": {
                    "ngram_tokenizer": {
                        "type": "ngram",
                        "min_gram": 2,
                        "max_gram": 10,
                        "token_chars": ["letter", "digit"]
                    }
                }
            }
        },
        "mappings": {
            "properties": {
                "name": {
                    "type": "text",
                    "analyzer": "substring_analyzer"
                },
                "email": {
                    "type": "text",
                    "analyzer": "substring_analyzer"
                },
                "phone": {
                    "type": "text",
                    "analyzer": "substring_analyzer"
                }
            }
        }
    }`

	ctx := context.Background()
	res, err := r.client.Indices.Create(
		indexName,
		r.client.Indices.Create.WithBody(strings.NewReader(mapping)),
		r.client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		// Ignore if index already exists
		if res != nil && strings.Contains(res.String(), "resource_already_exists_exception") {
			return nil
		}
		log.Printf("Error creating index: %v", err)
		return err
	}
	defer res.Body.Close()

	return nil
}

func (r *ElasticRepositoryImpl) InsertUser(user models.User) error {
	// Create index with mapping if it doesn't exist
	if err := r.createIndexIfNotExists(); err != nil {
		log.Printf("Error ensuring index exists: %v", err)
		return err
	}

	body, err := json.Marshal(user)
	if err != nil {
		log.Println("error in marshalling to json ", err)
		return err
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: user.ID,
		Body:       strings.NewReader(string(body)),
		Refresh:    "true",
	}

	ctx := context.Background()
	res, err := req.Do(ctx, r.client)
	if err != nil {
		log.Println("error indexing the document ", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("error response from elasticsearch: %s", res.String())
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return nil
}

func (r *ElasticRepositoryImpl) SearchUser(query string) ([]models.User, error) {
	var results []models.User

	queryBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":    query,
				"fields":   []string{"name", "email", "phone"},
				"analyzer": "substring_analyzer", // Use our custom analyzer
				"type":     "best_fields",        // Matches any field with substring
			},
		},
	}

	body, err := json.Marshal(queryBody)
	if err != nil {
		log.Println("error in converting to json ", err)
		return results, err
	}

	ctx := context.Background()
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(indexName),
		r.client.Search.WithBody(strings.NewReader(string(body))),
	)

	if err != nil {
		log.Println("error searching documents ", err)
		return results, err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("error response from elasticsearch: %s", res.String())
		return results, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var response struct {
		Hits struct {
			Hits []struct {
				Source models.User `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println("error in decoding the response ", err)
		return results, err
	}

	for _, hit := range response.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}
