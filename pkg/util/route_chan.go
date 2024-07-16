package util

import "context"

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
