package handler

import (
	"log"

	saga "github.com/MilosNesh/soa-team-13/common/saga/messaging"
	events "github.com/MilosNesh/soa-team-13/common/saga/purchase_tour"
)

type ProfileService interface {
	HasSufficientBalance(userID string, amount float64) (bool, error)
	DeductBalance(userID string, amount float64) error
	RefundBalance(userID string, amount float64) error
}

type PurchaseStakeholdersHandler struct {
	profileService    ProfileService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewPurchaseStakeholdersHandler(
	profileService ProfileService,
	replyPublisher saga.Publisher,
	commandSubscriber saga.Subscriber,
) (*PurchaseStakeholdersHandler, error) {
	h := &PurchaseStakeholdersHandler{
		profileService:    profileService,
		replyPublisher:    replyPublisher,
		commandSubscriber: commandSubscriber,
	}
	if err := h.commandSubscriber.Subscribe(h.handle); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *PurchaseStakeholdersHandler) handle(cmd *events.PurchaseCartCommand) {
	reply := events.PurchaseCartReply{Purchase: cmd.Purchase}

	switch cmd.Type {
	case events.CheckBalance:
		log.Printf("[Stakeholders] Checking balance for user %s, amount: $%.2f", cmd.Purchase.UserID, cmd.Purchase.TotalPrice)
		ok, err := h.profileService.HasSufficientBalance(cmd.Purchase.UserID, cmd.Purchase.TotalPrice)
		if err != nil {
			log.Printf("[Stakeholders] Check balance error: %v", err)
			reply.Type = events.BalanceInsufficient
			reply.Message = err.Error()
		} else if ok {
			reply.Type = events.BalanceSufficient
			log.Printf("[Stakeholders] Balance sufficient for user %s", cmd.Purchase.UserID)
		} else {
			reply.Type = events.BalanceInsufficient
			reply.Message = "Insufficient balance"
			log.Printf("[Stakeholders] Insufficient balance for user %s", cmd.Purchase.UserID)
		}

	case events.DeductBalance:
		log.Printf("[Stakeholders] Deducting $%.2f from user %s", cmd.Purchase.TotalPrice, cmd.Purchase.UserID)
		if err := h.profileService.DeductBalance(cmd.Purchase.UserID, cmd.Purchase.TotalPrice); err != nil {
			log.Printf("[Stakeholders] Deduct balance error: %v", err)
			reply.Type = events.BalanceDeductionFailed
			reply.Message = err.Error()
		} else {
			reply.Type = events.BalanceDeducted
			log.Printf("[Stakeholders] Balance deducted successfully for user %s", cmd.Purchase.UserID)
		}

	case events.RefundBalance:
		log.Printf("[Stakeholders] Refunding $%.2f to user %s", cmd.Purchase.TotalPrice, cmd.Purchase.UserID)
		if err := h.profileService.RefundBalance(cmd.Purchase.UserID, cmd.Purchase.TotalPrice); err != nil {
			log.Printf("[Stakeholders] Refund balance error: %v", err)
			reply.Type = events.PurchaseAborted
			reply.Message = err.Error()
		} else {
			reply.Type = events.BalanceRefunded
			log.Printf("[Stakeholders] Balance refunded successfully for user %s", cmd.Purchase.UserID)
		}

	default:
		reply.Type = events.UnknownReply
		log.Printf("[Stakeholders] Unknown command type: %v", cmd.Type)
	}

	if reply.Type != events.UnknownReply {
		log.Printf("[Stakeholders] Publishing reply type: %v for purchase %s", reply.Type, cmd.Purchase.PurchaseID)
		_ = h.replyPublisher.Publish(&reply)
	}
}
