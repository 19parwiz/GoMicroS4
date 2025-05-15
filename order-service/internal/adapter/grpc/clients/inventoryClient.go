package clients

import (
	"context"
	"github.com/19parwiz/order-service/internal/domain"
	proto "github.com/19parwiz/order-service/protos/gen/golang"
	"time"
)

type InventoryClient struct {
	client proto.InventoryServiceClient
}

func NewInventoryClient(client proto.InventoryServiceClient) *InventoryClient {
	return &InventoryClient{client: client}
}

func (c *InventoryClient) GetProduct(ctx context.Context, productID uint64) (domain.Product, error) {
	req := &proto.GetProductRequest{
		ProductId: productID,
	}
	resp, err := c.client.GetProduct(ctx, req)
	if err != nil {
		return domain.Product{}, err
	}

	return domain.Product{
		ID:        resp.ProductId,
		Name:      resp.Name,
		Category:  resp.Category,
		Price:     resp.Price,
		Stock:     resp.Stock,
		CreatedAt: parseTime(resp.CreatedAt),
		UpdatedAt: parseTime(resp.UpdatedAt),
	}, nil
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr) // Add error handling if needed
	return t
}
