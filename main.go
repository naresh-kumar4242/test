package main

import (
	"fmt"
	"sync"
	"time"
)

// Stack for storing the stuff
type Stack struct {
	data []interface{} // cints, strings
	lock sync.Mutex    // to avoid race condition 
}

// Push ... adds something to the top of the stack
func (s *Stack) Push(item interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock() // 
	s.data = append(s.data, item)
}

// Pop ... takes something off the top
func (s *Stack) Pop() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(s.data) == 0 {
		return nil
	}
	// Grab the last item and shorten the slice
	item := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return item
}

func main() {
	stack := &Stack{}

	// goroutine to add stuff to the stack
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Println("Pushing", i)
			stack.Push(i) // An integer!
			stack.Push(fmt.Sprintf("String %d", i))
			time.Sleep(500 * time.Millisecond) 
		}
		fmt.Println("Done pushing items") 
	}()

	// These channels will hold the popped values
	intChan := make(chan int)
	stringChan := make(chan string)

	// This goroutine will try to pop integers
	go func() {
		for {
			val := stack.Pop()
			if val == nil {
				fmt.Println("Int popper: Nothing")
				time.Sleep(100 * time.Millisecond) 
				continue
			}
			if num, ok := val.(int); ok {
				fmt.Println("Int popper: Found an int!", num)
				intChan <- num
			} else {
				fmt.Println("Int popper:  not an int")
			}
		}
	}()

	// to pop strings
	go func() {
		for {
			val := stack.Pop()
			if val == nil {
				fmt.Println("String popper: Stack's empty")
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if str, ok := val.(string); ok {
				fmt.Println("String popper: Found a string!", str)
				stringChan <- str
			} else {
				fmt.Println("String popper: not a string...")
			}
		}
	}()

	// main logic go herre
	for {
		select {
		case num := <-intChan:
			fmt.Println("Got an int:", num)
			time.Sleep(time.Second) 
		case str := <-stringChan:
			fmt.Println("Got a string:", str)
			time.Sleep(time.Second) 
		}
	}
}
