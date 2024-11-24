package staterepo

import (
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/repository/models"
)

// type State struct {
// 	Connects []Connect `json:"connects" db:"connects"`
// 	// Contacts []Contact `json:"contacts" db:"contacts"`
// 	// Chats    []Chat    `json:"chats" db:"chats"`
// }

// type Connect struct {
// 	Conn      *websocket.Conn
// 	Pubsub    uuid.UUID `json:"pubsub" db:"pubsub"`
// 	CreatedAt int64     `json:"created_at" db:"created_at"`
// }

func stateToDomain(state models.State) domain.State {
	domainConnects := make([]domain.Connect, len(state.Connects))
	for i, connect := range state.Connects {
		domainConnects[i].Conn = connect.Conn
		domainConnects[i].Pubsub = connect.Pubsub
		domainConnects[i].CreatedAt = connect.CreatedAt
	}
	return domain.State{
		Connects: domainConnects,
	}
}

// func domainToState() {

// }

// type Contact struct {
// 	UserID  uint64
// 	Status  string
// 	Event   string
// 	Name    string
// 	Surname string
// 	Email   string
// }

// func domainToProvider(provider domain.Provider) models.Provider {
// 	return models.Provider{
// 		ID:     provider.ID(),
// 		Name:   provider.Name(),
// 		Origin: provider.Origin(),
// 	}
// }

// func providerToDomain(provider models.Provider) domain.Provider {
// 	return domain.NewProvider(domain.NewProviderData{
// 		ID:     provider.ID,
// 		Name:   provider.Name,
// 		Origin: provider.Origin,
// 	})
// }

// func domainToProduct(product domain.Product) models.Product {
// 	return models.Product{
// 		ID:         product.ID(),
// 		Name:       product.Name(),
// 		ProviderID: product.ProviderID(),
// 		Price:      product.Price(),
// 		Stock:      product.Stock(),
// 	}
// }

// func productToDomain(product models.Product) domain.Product {
// 	return domain.NewProduct(domain.NewProductData{
// 		ID:         product.ID,
// 		Name:       product.Name,
// 		ProviderID: product.ProviderID,
// 		Price:      product.Price,
// 		Stock:      product.Stock,
// 	})
// }
// func domainToOrderState(orderState domain.OrderState) models.OrderState {
// 	return models.OrderState{
// 		ID:   orderState.ID(),
// 		Name: orderState.Name(),
// 	}
// }

// func orderStateToDomain(orderState models.OrderState) domain.OrderState {
// 	return domain.NewOrderState(domain.NewOrderStateData{
// 		ID:   orderState.ID,
// 		Name: orderState.Name,
// 	})
// }

// func domainToUser(user domain.User) models.User {
// 	return models.User{
// 		ID:        user.ID(),
// 		Login:     user.Login(),
// 		Password:  user.Password(),
// 		Role:      user.Role(),
// 		Token:     user.Token(),
// 		CreatedAt: user.CreratedAt(),
// 		UpdatedAt: user.UpdatedAt(),
// 	}
// }

// func userToDomain(user models.User) domain.User {
// 	return domain.NewUser(domain.NewUserData{
// 		ID:        user.ID,
// 		Login:     user.Login,
// 		Password:  user.Password,
// 		Role:      user.Role,
// 		Token:     user.Token,
// 		CreatedAt: user.CreatedAt,
// 		UpdatedAt: user.UpdatedAt,
// 	})
// }

// func domainToOrder(order domain.Order) models.Order {
// 	return models.Order{
// 		ID:          order.ID(),
// 		UserID:      order.UserID(),
// 		StateID:     order.StateID(),
// 		TotalAmount: order.TotalAmount(),
// 		CreatedAt:   order.CreatedAt(),
// 	}
// }

// func orderToDomain(order models.Order) (domain.Order, error) {
// 	return domain.NewOrder(domain.NewOrderData{
// 		ID:          order.ID,
// 		UserID:      order.UserID,
// 		StateID:     order.StateID,
// 		TotalAmount: order.TotalAmount,
// 		CreatedAt:   order.CreatedAt,
// 	})
// }

// func domainToItem(item domain.Item) models.Item {
// 	return models.Item{
// 		ID:         item.ID(),
// 		ProductID:  item.ProductID(),
// 		Quantity:   item.Quantity(),
// 		TotalPrice: item.TotalPrice(),
// 		OrderID:    item.OrderID(),
// 	}
// }

// func itemToDomain(item models.Item) (domain.Item, error) {
// 	return domain.NewItem(domain.NewItemData{
// 		ID:         item.ID,
// 		ProductID:  item.ProductID,
// 		Quantity:   item.Quantity,
// 		TotalPrice: item.TotalPrice,
// 		OrderID:    item.OrderID,
// 	})
// }
