# go-snippets
This repository contains some boilerplate code and helper functions for some of the more commonly repeated Go patterns.

## *channel* package

func  **Clone**[T  any](inCh  chan  T, size  int) <-chan  <-chan  T

Clone is a generic function that duplicates a channel. The number of duplicates it creates depends on the size parameter. The type and the buffer size of the cloned channels are the same as the source channel.

**[Example]**

Say this is our current pipeline. intStream is consumed by a function call consumeOne. This is working fine, but what if we have another function consumeTwo that also needs to access the **entirety** of intStream? 
<code>
intStream := make(chan int, 10)
consumeOne(intStream)	// This is fine
consumeTwo(intStream)   // This is not
</code>

To solve this, we can create two clones of the channel, one for each of the consumers of the input channel.

<code>
intStream := make(chan int, 10)
intStreamChs := Clone(intStream, 2)
consumeOne(<-intStreamChs)
consumeTwo(<-intStreamChs)
</code>
