package pkgmanager

import src "github.com/nevalang/neva/pkg/sourcecode"

type queue []src.ModuleRef

func (q *queue) enqueue(deps map[string]src.ModuleRef) {
	for _, dep := range deps {
		*q = append(*q, dep)
	}
}

func (q *queue) dequeue() src.ModuleRef {
	tmp := *q
	last := (tmp)[len(tmp)-1]
	*q = (tmp)[:len(tmp)-1]
	return last
}

func (q *queue) empty() bool {
	return len(*q) == 0
}

func newQueue(deps map[string]src.ModuleRef) *queue {
	q := make(queue, len(deps))
	q.enqueue(deps)
	return &q
}
