package syncTask

func Call(fn func() error) error {
	err := fn()
	if err != nil {
		return err
	}
	return nil
}
