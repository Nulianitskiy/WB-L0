package mycache

import (
	"WB-L0/internal/model"
	"github.com/patrickmn/go-cache"
)

func GetCacheOrders(c *cache.Cache) (orders []model.Order) {
	items := c.Items()
	for _, item := range items {
		orders = append(orders, item.Object.(model.Order))
	}
	return orders
}

func GetCacheOrdersById(c *cache.Cache, id string) (orders []model.Order) {
	items := c.Items()
	for _, item := range items {
		order := item.Object.(model.Order)
		if order.OrderUID == id {
			orders = append(orders, item.Object.(model.Order))
		}
	}
	return orders
}
