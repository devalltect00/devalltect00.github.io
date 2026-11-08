// Deployment runtime composes reviewed clients for the explicitly gated local operator command.
package main

import (
	"context"
	"strings"
	"sync/atomic"
	"time"
)

// deploymentRuntime owns separately reviewed policy and a one-shot execution budget, never public request authority.
type deploymentRuntime struct {
	policy       deploymentPolicy
	issuerOrigin string
	spent        atomic.Bool
}

// deploymentClients groups concrete fixed-origin adapters for one run and owns their idle socket cleanup.
type deploymentClients struct {
	issuer   *workflowOIDC
	receiver *receiverClient
	creator  *pagesCreator
}

// close releases sockets but never resets issuer, sender or server-side admission budgets.
func (c deploymentClients) close() {
	if c.issuer != nil {
		c.issuer.close()
	}
	if c.receiver != nil {
		c.receiver.close()
	}
	if c.creator != nil {
		c.creator.Close()
	}
}

// runtimeClients constructs adapters only after local verification; its credentials never leave their assigned transports.
func runtimeClients(p deploymentPolicy, origin string, settings deploymentSettings, plan receiverPlan, get func(string) string) (deploymentClients, error) {
	var clients deploymentClients
	issuer, e := newWorkflowOIDC(p, origin, get)
	if e != nil {
		return clients, errReceiver
	}
	clients.issuer = issuer
	receiver, e := newReceiverClient(settings.Receiver)
	if e != nil {
		clients.close()
		return deploymentClients{}, errReceiver
	}
	clients.receiver = receiver
	repository := strings.Split(settings.Repository, "/")
	creator, e := newPagesCreator(repository[0], repository[1], get("GITHUB_TOKEN"), plan)
	if e != nil {
		clients.close()
		return deploymentClients{}, errReceiver
	}
	clients.creator = creator
	return clients, nil
}

// run supplies only production constructors; no loopback override or factory is exposed through CLI/workflow input.
func (d *deploymentRuntime) run(ctx context.Context, raw, candidate, uploaded []byte, get func(string) string) (senderOutcome, error) {
	return d.runWith(ctx, raw, candidate, uploaded, get, runtimeClients)
}

// runWith permits a package-local fixture factory; injected adapters are trusted code, not a production extension point.
// It freezes local bytes, validates context/content before credentials, then joins issuance and the protected sender once.
func (d *deploymentRuntime) runWith(ctx context.Context, raw, candidate, uploaded []byte, get func(string) string, build func(deploymentPolicy, string, deploymentSettings, receiverPlan, func(string) string) (deploymentClients, error)) (senderOutcome, error) {
	var out senderOutcome
	if d == nil || !d.spent.CompareAndSwap(false, true) || build == nil || ctx.Err() != nil {
		return out, errReceiver
	}
	ctx, cancel := context.WithTimeout(ctx, 75*time.Second)
	defer cancel()
	settings, e := validateDeploymentSettings(d.policy, get)
	if e != nil {
		return out, errReceiver
	}
	request, e := ParseDeployRequest(raw)
	if e != nil || request.Stage.HeadSHA != settings.Receiver.HeadSHA || settings.Receiver.RunID == request.RunID || settings.Receiver.RunID == request.Stage.RunID || len(candidate) > MaxArchive || len(uploaded) > MaxDeployWrapper {
		return out, errReceiver
	}
	candidate = append([]byte(nil), candidate...)
	uploaded = append([]byte(nil), uploaded...)
	if verifyDeploymentInput(ctx, request, candidate, uploaded) != nil || ctx.Err() != nil {
		return out, errReceiver
	}
	clients, e := build(d.policy, d.issuerOrigin, settings, planFor(request, settings.Receiver.RunID), get)
	if e != nil {
		clients.close()
		return out, errReceiver
	}
	defer clients.close()
	if clients.issuer == nil || clients.receiver == nil || clients.creator == nil {
		return out, errReceiver
	}
	token, e := clients.issuer.acquire(ctx)
	if e != nil {
		return out, errReceiver
	}
	sender := protectedSender{receiver: clients.receiver, creator: clients.creator}
	return sender.send(ctx, request, candidate, uploaded, token, clients.issuer.acquire)
}
