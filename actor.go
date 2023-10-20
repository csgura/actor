package actor

import "github.com/asynkron/protoactor-go/actor"

type Actor interface {
	Receive(c Context)
}

type actorWrapper struct {
	actor Actor
}

func (r *actorWrapper) Receive(c actor.Context) {
	r.actor.Receive(ContextWrapper{c})
}

func (r *actorWrapper) Unwrap() interface{} {
	return r.actor
}

// The ReceiveFunc type is an adapter to allow the use of ordinary functions as actors to process messages
type ReceiveFunc func(c Context)

// Receive calls f(c)
func (f ReceiveFunc) Receive(c Context) {
	f(c)
}

type ReceiverFunc func(c ReceiverContext, envelope *MessageEnvelope)

type ReceiverMiddleware func(next ReceiverFunc) ReceiverFunc

type Props struct {
	*actor.Props
}

func (props *Props) WithReceiverMiddleware(middleware ...ReceiverMiddleware) *Props {

	if len(middleware) == 0 {
		return props
	}

	m := make([]actor.ReceiverMiddleware, len(middleware))

	for i := range middleware {
		if middleware[i] != nil {
			m[i] = func(next actor.ReceiverFunc) actor.ReceiverFunc {

				cn := func(c ReceiverContext, envelope *MessageEnvelope) {
					next(c.protoReceiverContext(), envelope)
				}

				ret := middleware[i](cn)

				return func(c actor.ReceiverContext, envelope *actor.MessageEnvelope) {

					if fc, ok := c.(actor.Context); ok {
						ret(ContextWrapper{fc}, envelope)
					} else {
						ret(ReceiverContextWrapper{c}, envelope)
					}
				}
			}
		}
	}

	return &Props{props.Props.Configure(actor.WithReceiverMiddleware(m...))}
}

type ActorFunc = ReceiveFunc
