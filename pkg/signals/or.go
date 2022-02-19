package signals

func Or(channels ...<-chan struct{}) <-chan struct{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			or(done, channels)
		}
	}()

	return done
}

func or(done <-chan struct{}, channels []<-chan struct{}) {
	select {
	case <-channels[0]:
	case <-channels[1]:
	case <-channels[2]:
	case <-Or(append(channels[3:], done)...):
	}
}
