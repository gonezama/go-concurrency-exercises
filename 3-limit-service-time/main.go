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
import ("time"
		"sync")

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mu sync.Mutex

}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	u.mu.Lock()
	if u.TimeUsed >= 10 {
		u.mu.Unlock()
		return false
	}
	u.mu.Unlock()

	remaining := 10 - u.TimeUsed
	if remaining <= 0 {
		return false // already used up quota
	}
	
	c1 := make(chan struct{}) // a channel for structs don't use memmory and by convention is just for signaling without information
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		process()
        close(c1)
		}()
		
	for { 
		select {
			case <-c1:
				return true
			case <-ticker.C:
				u.mu.Lock()
				u.TimeUsed++
				if u.TimeUsed >= 10 {
					u.mu.Unlock()
					return false // Quota exceeded mid-process
				}				
				u.mu.Unlock()
    }
}}

func main() {
	RunMockServer()
}
