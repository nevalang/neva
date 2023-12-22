package builder

import src "github.com/nevalang/neva/internal/compiler/sourcecode"

type Queue []src.ModuleRef

func (q *Queue) Enqueue(deps map[string]src.ModuleRef) {
	for _, dep := range deps {
		*q = append(*q, dep)
	}
}

func (q *Queue) Dequeue() src.ModuleRef {
	tmp := *q
	last := (tmp)[len(tmp)-1]
	*q = (tmp)[:len(tmp)-1]
	return last
}

func (q *Queue) Empty() bool {
	return len(*q) == 0
}

func NewQueue(deps map[string]src.ModuleRef) *Queue {
	q := make(Queue, len(deps))
	q.Enqueue(deps)
	return &q
}
