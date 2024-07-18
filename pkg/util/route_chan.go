package util

import "context"

// RouteChanContext monitors [context.Context] and reciever channel.
//
// It also loops until interrupted by ctx or reciever is closed.
func RouteChanContext[T any](ctx context.Context, reciever <-chan T) func(yield func(val T) bool) {
	return func(yield func(val T) bool) {
		for {
			select {
			case <-ctx.Done():
				return

			case val, exist := <-reciever:
				if !exist {
					return
				} else if !yield(val) {
					return
				}
			}
		}
	}
}
