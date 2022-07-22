package backgroundTask

type Background struct {
	deferFn func()
	errorfn func(err error)
}

func (b *Background) SetDeferFn(deferFn func()) *Background {
	b.deferFn = deferFn
	return b
}
func (b *Background) SetErrorFn(errorFn func(err error)) *Background {
	b.errorfn = errorFn
	return b
}
func (b *Background) Run(fn func() error) {
	go func() {
		defer func() {
			if b.deferFn != nil {
				b.deferFn()
			}
		}()
		err := fn()
		if err != nil {
			b.errorfn(err)
		}
	}()
}
func New() *Background {
	return &Background{}
}
