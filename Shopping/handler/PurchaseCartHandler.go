package handler

import (
	"Shopping/model"
	"Shopping/service"
	"log"

	saga "github.com/MilosNesh/soa-team-13/common/saga/messaging"
	events "github.com/MilosNesh/soa-team-13/common/saga/purchase_tour"
)

// Apstrahuj servis(e) koje koristi handler
type CartService interface {
	ClearCart(cartID string) error
	CreatePurchaseTokens(userID string, items []events.CartItem) ([]model.TourPurchaseToken, error)
}

type PurchaseCartHandler struct {
	cartService       *service.ShoppingCartService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewPurchaseCartHandler(
	cartService *service.ShoppingCartService,
	replyPublisher saga.Publisher,
	commandSubscriber saga.Subscriber,
) (*PurchaseCartHandler, error) {
	h := &PurchaseCartHandler{
		cartService:       cartService,
		replyPublisher:    replyPublisher,
		commandSubscriber: commandSubscriber,
	}
	if err := h.commandSubscriber.Subscribe(h.handle); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *PurchaseCartHandler) handle(cmd *events.PurchaseCartCommand) {
	reply := events.PurchaseCartReply{Purchase: cmd.Purchase}

	switch cmd.Type {

	case events.CalculateTotal:
		log.Printf("[Shopping-Handler] Calculating total for purchase %s", cmd.Purchase.PurchaseID)
		total := 0.0
		for _, it := range cmd.Purchase.Items {
			total += it.Price
			log.Printf("[Shopping-Handler]   Item %s: $%.2f", it.ItemID, it.Price)
		}
		reply.Purchase.TotalPrice = total
		reply.Type = events.TotalCalculated
		log.Printf("[Shopping-Handler] Total calculated: $%.2f for purchase %s", total, cmd.Purchase.PurchaseID)

	case events.CompletePurchase:
		log.Printf("[Shopping-Handler] Completing purchase %s - clearing cart %s", cmd.Purchase.PurchaseID, cmd.Purchase.CartID)
		if _, err := h.cartService.CreatePurchaseTokens(cmd.Purchase.UserID, cmd.Purchase.CartID, cmd.Purchase.Items); err != nil {
			log.Printf("[Shopping-Handler] Create purchase tokens failed for cart %s: %v", cmd.Purchase.CartID, err)
			reply.Type = events.PurchaseAborted
			reply.Message = err.Error()
			break
		}
		reply.Type = events.PurchaseCompleted
		log.Printf("[Shopping-Handler] Purchase completed successfully for %s", cmd.Purchase.PurchaseID)

	default:
		reply.Type = events.UnknownReply
		log.Printf("[Shopping-Handler] Unknown command type: %v", cmd.Type)
	}

	if reply.Type != events.UnknownReply {
		log.Printf("[Shopping-Handler] Publishing reply type: %v for purchase %s", reply.Type, cmd.Purchase.PurchaseID)
		_ = h.replyPublisher.Publish(&reply)
	}
}
