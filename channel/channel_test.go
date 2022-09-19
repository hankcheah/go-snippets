package channel

import (
    "testing"
    "reflect"
)

// TestCloneBridge creates a single duplicate, so it's like a bridge
// The output should match the input
func TestCloneBridgeOne(t *testing.T) {
    inCh := make(chan int, 2)
    outCh := <-Clone(inCh, 1)

    inCh <- 1
    inCh <- 2

    res := []int{<-outCh, <-outCh}
    expected := []int{1, 2}
    if !reflect.DeepEqual(res, expected) {
        t.Errorf("Expected output. Want: %v Got %v.", expected, res)
    }
}

// TestCloneBridge creates a single duplicate
// The difference between this and TestCloneBridgeOne is that in this
// test case, the values in the input channel is preloaded before it gets cloned.
func TestCloneBridgeTwo(t *testing.T) {
    inCh := make(chan int, 2)
    inCh <- 1
    inCh <- 2

    outCh := <-Clone(inCh, 1)

    res := []int{<-outCh, <-outCh}
    expected := []int{1, 2}
    if !reflect.DeepEqual(res, expected) {
        t.Errorf("Expected output. Want: %v Got %v.", expected, res)
    }
}

// TestCloneTwo tests a integer tee (one input channel cloned to two output channels)
func TestCloneTwoInt(t *testing.T) {
    inCh := make(chan int, 3)
    outCh := Clone(inCh, 2)

    // Fan out to two output channels
    outCh1 := <-outCh
    outCh2 := <-outCh

    inCh <- 1
    inCh <- 2
    inCh <- 3

    res := []int{<-outCh1, <-outCh1, <-outCh1, <-outCh2, <-outCh2, <-outCh2}
    expected := []int{1, 2, 3, 1, 2, 3}
    if !reflect.DeepEqual(res, expected) {
        t.Errorf("Expected output. Want: %v Got %v.", expected, res)
    }
}

// TestCloneTwo creates a a string tee using Clone and test its validity.
func TestCloneTwoString(t *testing.T) {
    inCh := make(chan string, 2)
    inChs := Clone(inCh, 2)

    outCh1 := <-inChs
    outCh2 := <-inChs

    // Send two values on the input channel exactly one time
    inCh <- "one"
    inCh <- "two"

    res := []string{<-outCh1, <-outCh1, <-outCh2, <-outCh2}
    expected := []string{"one", "two", "one", "two"}
    if !reflect.DeepEqual(res, expected) {
        t.Errorf("Expected output. Want: %v Got %v.", expected, res)
    }
}

// TestCloneFiftyInt clones the input channel fifty times.
func TestCloneNInt(t *testing.T) {
    count := 50

    inCh := make(chan int, 2)
    inChs := Clone(inCh, count)

    outChs := make([]<-chan int, count)

    for i := 0; i < count; i++ {
        outChs[i] = <-inChs
    }

    // Send two values on the input channel exactly once
    inCh <- 3
    inCh <- 2

    res := []int{}
    expected := []int{}

    for i := 0; i < count; i++ {
        res = append(res, <-outChs[i], <-outChs[i])
        expected = append(expected, 3, 2)
    }

    if !reflect.DeepEqual(res, expected) {
        t.Errorf("Expected output. Want: %v Got %v.", expected, res)
    }
}
