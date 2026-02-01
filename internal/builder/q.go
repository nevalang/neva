package builder

import "github.com/nevalang/neva/internal/compiler/ast/core"

type queue []core.ModuleRef

func (q *queue) enqueue(deps map[string]core.ModuleRef) {
	for _, dep := range deps {
		*q = append(*q, dep)
	}
}

func (q *queue) dequeue() core.ModuleRef {
	tmp := *q
	last := (tmp)[len(tmp)-1]
	*q = (tmp)[:len(tmp)-1]
	return last
}

func (q *queue) empty() bool {
	return len(*q) == 0
}

func newQueue(deps map[string]core.ModuleRef) *queue {
	q := make(queue, 0, len(deps))
	q.enqueue(deps)
	return &q
}
