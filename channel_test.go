package main

import (
    "testing"
    "reflect"
)

// TestFanoutBridge creates a one-to-one fanout or a bridge
// The output should match the input
func TestFanoutBridge(t *testing.T) {
    inCh := make(chan int, 2)
    outCh := <-Fanout(inCh, 1)

    inCh <- 1
    inCh <- 2

    res := []int{<-outCh, <-outCh}
    expected := []int{1, 2}
    if !reflect.DeepEqual(res, expected) {
        t.Errorf("Expected output. Want: %v Got %v.", expected, res)
    }
}

// TestFanoutTwo creates a integer tee (one-to-two fanout) test it
func TestFanoutTwoInt(t *testing.T) {
    inCh := make(chan int, 3)
    inChs := Fanout(inCh, 2)

    // Fan out to two output channels
    outCh1 := <-inChs
    outCh2 := <-inChs

    inCh <- 1
    inCh <- 2
    inCh <- 3

    res := []int{<-outCh1, <-outCh1, <-outCh1, <-outCh2, <-outCh2, <-outCh2}
    expected := []int{1, 2, 3, 1, 2, 3}
    if !reflect.DeepEqual(res, expected) {
        t.Errorf("Expected output. Want: %v Got %v.", expected, res)
    }
}

// TestFanoutTwo creates a string tee (one-to-two fanout) and test it
func TestFanoutTwoString(t *testing.T) {
    inCh := make(chan string, 2)
    inChs := Fanout(inCh, 2)

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

// TestFanoutNInt tests a large-ish fanout setup.
func TestFanoutNInt(t *testing.T) {
    fanoutSize := 50

    inCh := make(chan int, 2)
    inChs := Fanout(inCh, fanoutSize)

    outChs := make([]<-chan int, fanoutSize)

    for i := 0; i < fanoutSize; i++ {
        outChs[i] = <-inChs
    }

    // Send two values on the input channel exactly once
    inCh <- 3
    inCh <- 2

    res := []int{}
    expected := []int{}

    for i := 0; i < fanoutSize; i++ {
        res = append(res, <-outChs[i], <-outChs[i])
        expected = append(expected, 3, 2)
    }

    if !reflect.DeepEqual(res, expected) {
        t.Errorf("Expected output. Want: %v Got %v.", expected, res)
    }
}
