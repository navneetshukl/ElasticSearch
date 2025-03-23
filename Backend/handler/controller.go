package handler

import (
	"context"
	"elasticsearch/elasticsearch"
	"elasticsearch/models"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/bxcodec/faker/v4"
)

var ch chan models.User

var (
	size   int    = 2000
	INSERT string = "INSERT"
	RANDOM string = "RANDOMDATAGENERATE"
)

type ElasticSvcImpl struct {
	repo   elasticsearch.ElasticSearchRepo
	wg     *sync.WaitGroup
	errMap map[string][]string
	mut    *sync.Mutex
}

func NewElasticSvcUsecase(repo elasticsearch.ElasticSearchRepo) models.ElasticSrvUseCase {
	return &ElasticSvcImpl{
		repo:   repo,
		wg:     &sync.WaitGroup{},
		errMap: map[string][]string{},
		mut:    &sync.Mutex{},
	}
}

func (e *ElasticSvcImpl) generateRandomData(ch chan models.User) chan models.User {
	for i := 1; i <= size; i++ {
		defer e.wg.Done()
		var u models.User
		err := faker.FakeData(&u)
		if err != nil {
			log.Println("error in generating fake data ", err)
			e.mut.Lock()
			e.errMap[RANDOM] = append(e.errMap[RANDOM], "error in generating fake data")
			e.mut.Unlock()
		}
		ch <- u
	}
	return ch
}

func (e *ElasticSvcImpl) insertIntoElastic(ch chan models.User) {
	defer e.wg.Done()
	for val := range ch {
		err := e.repo.InsertUser(val)
		if err != nil {
			log.Println("error in inserting the user to elastic ", err)
			e.mut.Lock()
			e.errMap[INSERT] = append(e.errMap[INSERT], fmt.Sprintf("error in inserting %s userID to elastic ", val.ID))
			e.mut.Unlock()
		}
	}
}

func (e *ElasticSvcImpl) InsertToElastic(ctx context.Context) error {
	ch = make(chan models.User, size)
	for i := 1; i <= 5; i++ {
		go e.generateRandomData(ch)
	}
	for i := 1; i <= 5; i++ {
		e.wg.Add(1)
		go e.insertIntoElastic(ch)
	}
	e.wg.Wait()
	if len(e.errMap[RANDOM]) > 0 || len(e.errMap[INSERT]) > 0 {
		log.Println("error occured in generating data or inserting data")
		return errors.New("error in inserting")
	}

	return nil
}

func (e *ElasticSvcImpl) GetData(ctx context.Context, query string) ([]models.User, error) {
	user, err := e.repo.SearchUser(query)
	if err != nil {
		log.Println("error in getting the data of this query ", err)
		return user, err
	}

	if len(user) > 10 {
		user = user[:10]
	}
	return user, nil

}
