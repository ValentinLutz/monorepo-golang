package orders

import (
	"app/internal"
	"github.com/rs/zerolog"
)

type Service struct {
	logger              *zerolog.Logger
	orderRepository     *OrderRepository
	orderItemRepository *OrderItemRepository
	config              *internal.Config
}

func (s Service) name() {

}
