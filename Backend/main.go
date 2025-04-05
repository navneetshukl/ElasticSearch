package main

import (
	"elasticsearch/elasticsearch"
	"elasticsearch/handler"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	client, err := elasticsearch.ConnectToElastic()
	if err != nil {
		log.Println("error in connecting to elastic search")
		return
	}
	elasticSrv := elasticsearch.NewElasticRepository(client)
	elsUsecase := handler.NewElasticSvcUsecase(elasticSrv)

	hand := handler.NewHandler(elsUsecase)

	//elsUsecase.InsertToElastic(context.Background())

	app.Get("/search", hand.GetQuery)

	app.Listen(":8080")
}
