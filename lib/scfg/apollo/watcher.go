package apollo

import (
	"context"
	"mylib/scfg/kcfg"

	"github.com/philchia/agollo/v4"
)

type watcher struct {
	a        *apollo
	out      chan []*kcfg.KeyValue
	ctx      context.Context
	cancelFn func()
}

func newWatcher(a *apollo) kcfg.Watcher {
	outChan := make(chan []*kcfg.KeyValue)
	ctx, cancel := context.WithCancel(context.Background())
	w := &watcher{
		a:   a,
		out: outChan,
		ctx: ctx,
		cancelFn: func() {
			a.cliProxy.Client.OnUpdate(nil)
			cancel()
		},
	}
	a.cliProxy.Client.OnUpdate(w.onNamespaceChange)
	return w
}

func (w *watcher) onNamespaceChange(event *agollo.ChangeEvent) {
	kv, err := w.a.load(GetNSMapValue(event.Namespace))
	if err != nil {
		return
	}
	if kv != nil {
		w.out <- []*kcfg.KeyValue{kv}
	}
}

func (w *watcher) WaitRead() ([]*kcfg.KeyValue, error) {
	select {
	case kv := <-w.out:
		return kv, nil
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

func (w *watcher) Stop() error {
	if w.cancelFn != nil {
		w.cancelFn()
	}
	close(w.out)
	return nil
}
