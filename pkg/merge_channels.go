package pkg

func MergeChannels(left, right <-chan int) <-chan int {
	result := make(chan int)

	mergeWorker := func() {
		defer close(result)

		// Drain both input channels by selecting
		// on both of them, until one of them closes.
		// then drain the remaining channel.
		//
		// Note: there is a different way of implementing
		// this merge channels mechanism, where once a
		// channel is detected as being closed, it
		// is replaced with the nil channel. Attempted
		// reads from a nil channel block forever, so it
		// effectively removes the corresponding case from
		// the select (provided we remember to avoid
		// selecting after all channels are closed). That
		// approach is arguably more elegant, or horrifying,
		// depending on how you feel about nil channels.
		// ref: https://medium.com/justforfunc/why-are-there-nil-channels-in-go-9877cc0b2308

		leftOpen, rightOpen := true, true
		var v int
		var ok bool
		for leftOpen && rightOpen {
			select {
			case v, ok = <-left:
				if ok {
					result <- v
				} else {
					leftOpen = false
				}
			case v, ok = <-right:
				if ok {
					result <- v
				} else {
					rightOpen = false
				}
			}
		}
		// One of the channels was closed, and has been drained.
		// So, try to drain the other one.

		// Could use for statements with range clauses here.
		// ref: https://go.dev/ref/spec#For_statements
		for leftOpen {
			select {
			case v, ok = <-left:
				if ok {
					result <- v
				} else {
					leftOpen = false
				}
			}
		}
		for rightOpen {
			select {
			case v, ok = <-right:
				if ok {
					result <- v
				} else {
					rightOpen = false
				}
			}
		}
	}
	go mergeWorker()
	return result
}
