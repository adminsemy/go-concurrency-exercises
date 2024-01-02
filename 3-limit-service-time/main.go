//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	done      chan struct{}
	mx        sync.Mutex
}

func (u *User) Done() {
	u.mx.Lock()
	defer u.mx.Unlock()
	if u.done != nil {
		return
	}
	u.done = make(chan struct{})
	go func() {
		time.Sleep(10 * time.Second)
		close(u.done)
	}()
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	chProcess := make(chan struct{})
	u.Done()
	go func() {
		process()
		close(chProcess)
	}()
	for {
		select {
		case <-chProcess:
			return true
		case <-u.done:
			if u.IsPremium {
				return true
			}
			return false
		}
	}
}

func main() {
	RunMockServer()
}
