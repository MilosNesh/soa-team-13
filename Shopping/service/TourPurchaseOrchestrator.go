package service

import (
	"log"

	saga "github.com/MilosNesh/soa-team-13/common/saga/messaging"
	events "github.com/MilosNesh/soa-team-13/common/saga/purchase_tour"
)

type PurchaseCartOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewPurchaseCartOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*PurchaseCartOrchestrator, error) {
	o := &PurchaseCartOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	if err := o.replySubscriber.Subscribe(o.handle); err != nil {
		return nil, err
	}
	return o, nil
}

// Start: pokreni sagu sa CalculateTotal; prosledi stavke iz korpe i UserID
func (o *PurchaseCartOrchestrator) Start(purchaseID string, details events.PurchaseCartDetails) error {
	cmd := &events.PurchaseCartCommand{
		Type:     events.CalculateTotal,
		Purchase: details,
	}
	return o.commandPublisher.Publish(cmd)
}

func (o *PurchaseCartOrchestrator) handle(reply *events.PurchaseCartReply) {
	next := o.nextCommandType(reply.Type)
	if next == events.UnknownCommand {
		log.Printf("[Orchestrator] Saga finished/aborted for %s", reply.Purchase.PurchaseID)
		return
	}

	cmd := &events.PurchaseCartCommand{
		Purchase: reply.Purchase,
		Type:     next,
	}
	log.Printf("[Orchestrator] Next command '%v' for %s", cmd.Type, cmd.Purchase.PurchaseID)
	_ = o.commandPublisher.Publish(cmd)
}

func (o *PurchaseCartOrchestrator) nextCommandType(rt events.PurchaseCartReplyType) events.PurchaseCartCommandType {
	switch rt {
	// happy path
	case events.TotalCalculated:
		return events.CheckBalance
	case events.BalanceSufficient:
		return events.DeductBalance
	case events.BalanceDeducted:
		return events.CompletePurchase

	// failure grane
	case events.BalanceInsufficient:
		return events.AbortPurchase
	case events.BalanceDeductionFailed:
		return events.RefundBalance
	// refund done → end
	case events.BalanceRefunded:
		return events.AbortPurchase

	// terminal
	case events.PurchaseCompleted, events.PurchaseAborted:
		return events.UnknownCommand
	default:
		return events.UnknownCommand
	}
}
