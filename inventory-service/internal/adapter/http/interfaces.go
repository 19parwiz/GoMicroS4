package http

import (
	"github.com/19parwiz/inventory-service/internal/adapter/http/handler"
)

type ProductHandler interface {
	handler.ProductUseCase
}
