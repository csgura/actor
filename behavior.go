package actor

import (
	"log/slog"
)

// var (
// 	plog = slog.New(log.DebugLevel, "[ACTOR]")
// )

type Behavior []ReceiveFunc

func NewBehavior() Behavior {
	return make(Behavior, 0)
}

func (b *Behavior) Become(receive ReceiveFunc) {
	b.clear()
	b.push(receive)
}

func (b *Behavior) BecomeStacked(receive ReceiveFunc) {
	b.push(receive)
}

func (b *Behavior) UnbecomeStacked() {
	b.pop()
}

func (b *Behavior) Receive(context Context) {
	behavior, ok := b.peek()
	if ok {
		behavior(context)
	} else {
		slog.Error("empty behavior called", "pid", context.Self())
	}
}

func (b *Behavior) clear() {
	if len(*b) == 0 {
		return
	}

	for i := range *b {
		(*b)[i] = nil
	}
	*b = (*b)[:0]
}

func (b *Behavior) peek() (v ReceiveFunc, ok bool) {
	l := b.len()
	if l > 0 {
		ok = true
		v = (*b)[l-1]
	}
	return
}

func (b *Behavior) push(v ReceiveFunc) {
	*b = append(*b, v)
}

func (b *Behavior) pop() (v ReceiveFunc, ok bool) {
	l := b.len()
	if l > 0 {
		l--
		ok = true
		v = (*b)[l]
		(*b)[l] = nil
		*b = (*b)[:l]
	}
	return
}

func (b *Behavior) len() int {
	return len(*b)
}
