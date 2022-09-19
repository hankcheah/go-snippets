package channel

// Clone creates a number of duplicate channels based on the input channel.
// It does this by forwarding the messages the input channel receives onto the output channels that it creates.
// The number of clones this function creates depends on the <size> parameter the caller provides.
//
// It returns a channel of channels, more specifically, a read-only channel of read-only channels.
// To use one of its output channels, just pop a channel from the returned channel.
//
// [Example]
//
//  msgs := make(chan int, 5)
//
//  clonedChs := channel.Clone(msgs, 2) // create two cloned channels
// 
//  outCh1 := <-clonedChs    // pop a cloned channel
//  outCh2 := <-clonedChs    // pop a cloned channel
//
//  oriMsgCh <- 123         // send 123 on the input channel
//
//  fmt.Println(<-outCh1)   // receive 123 on output channel 1
//  fmt.Println(<-outCh2)   // receive 123 on output channel 2
//

func Clone[T any](inCh chan T, size int) <-chan <-chan T {
    // The channel of channels to return at the end of this function call
    ret := make(chan (<-chan T), size)

    // This slice keeps track of all the output channels this function will be creating below.
    outChs := make([]chan T, size)

    // Create channels, keep track of them in the slice and send them on the return channel
    for i := 0; i < size; i++ {
        // The buffer size of the newly created channel is the same as the input channel
        outChs[i] = make(chan T, cap(inCh))
        ret <- outChs[i]
    }

    // Start a goroutine to manage receiving message from the input channels and sending out to the output channels
    // Close the output channels if the input channel has been closed.
    go func () {
        for {
            msg, more := <-inCh
            if more {
                for _, ch := range outChs {
                    ch <- msg
                }
            } else {
                for _, ch := range outChs {
                    close(ch)
                }
                return
            }
        }
    }()

    return ret
}
