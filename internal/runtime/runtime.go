package runtime

import "fmt"

type Runtime interface {
	Run(p Opts) error
}

type Opts struct {
	env     Env
	root    string
	in, out map[string]chan Msg
}

type runtime struct{}

func (r runtime) Run(o Opts) error {
	mod, ok := o.env[o.root]
	if !ok {
		return fmt.Errorf("root module '%s' not found in env", o.root)
	}

	mod.Run(in, out)

	return nil
}

func (r runtime) spawnWorker() Worker {
	return Worker{}
}
