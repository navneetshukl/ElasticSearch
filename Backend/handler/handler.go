package handler

import (
	"context"
	"elasticsearch/models"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	ElasticSrv models.ElasticSrvUseCase
}

func NewHandler(el models.ElasticSrvUseCase) *Handler {
	return &Handler{
		ElasticSrv: el,
	}
}

func (h *Handler) GetQuery(ctx *fiber.Ctx) error {
	q := ctx.Query("search")
	if len(q) == 0 {
		return ctx.JSON(fiber.Map{
			"error":   true,
			"message": "please provide the valid query",
			"data":    []interface{}{},
		})

	}

	data, err := h.ElasticSrv.GetData(context.Background(), q)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"error":   true,
			"message": "something went wrong",
			"data":    []interface{}{},
		})

	}

	return ctx.JSON(fiber.Map{
		"error":   false,
		"message": "data fetched successfully",
		"data":    data,
	})
}
