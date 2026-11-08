// Protected sender joins reviewed adapters without a CLI, environment loader or hosted workflow.
package main

import (
	"context"
	"sync/atomic"
	"time"
)

// senderOutcome retains only an actual validated acknowledgement for explicit report-only recovery.
// It must never be logged; neither a receipt nor an acknowledgement certifies public completion.
type senderOutcome struct {
	Acknowledgement []byte
	Report          receiverReport
}

// protectedSender owns a workflow-local provider creator and fixed CMS receiver transport.
// The process spent bit supplements, but cannot replace, durable server admission.
type protectedSender struct {
	receiver *receiverClient
	creator  *pagesCreator
	spent    atomic.Bool
}

// send performs complete verification, original-token authorization/admission, one create, then reporting.
// freshReportToken is trusted runtime code; it is called only after a validated provider response.
// All errors fail closed; returned acknowledgement supports report-only recovery, never another create.
func (s *protectedSender) send(ctx context.Context, r DeployRequest, candidate, uploaded []byte, original string, freshReportToken func(context.Context) (string, error)) (senderOutcome, error) {
	var out senderOutcome
	if s == nil || !s.spent.CompareAndSwap(false, true) || s.receiver == nil || s.creator == nil || freshReportToken == nil || s.creator.plan != planFor(r, s.receiver.config.RunID) {
		return out, errReceiver
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	decision, e := s.receiver.authorizeVerified(ctx, r, candidate, uploaded, original)
	if e != nil {
		return out, errReceiver
	}
	exchange := receiverExchange{client: s.receiver}
	permit, e := exchange.beginVerified(ctx, r, candidate, uploaded, original, decision)
	if e != nil || !permit.validFor(r, s.receiver.config.RunID, time.Now()) || ctx.Err() != nil {
		return out, errReceiver
	}
	writeCtx, writeCancel := context.WithDeadline(ctx, permit.ExpiresAt)
	ack, e := s.creator.create(writeCtx, permit.receiverPlan, original)
	writeCancel()
	if e != nil {
		return out, errReceiver
	}
	out.Acknowledgement = append([]byte(nil), ack...)
	token, e := freshReportToken(ctx)
	if e != nil || ctx.Err() != nil {
		return out, errReceiver
	}
	out.Report, e = exchange.reportVerified(ctx, r, candidate, uploaded, token, ack)
	if e != nil {
		return out, errReceiver
	}
	return out, nil
}
