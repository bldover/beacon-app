package ticketmaster

// The ticketmaster public API has a rate limit of 5 requests per second.
// To handle exceeding the rate limit during the parallel repeated calls to fetch
// pages of event data, we use a queue of pending requests that will automatically
// retry after a delay when we get a rate exceeded response.
//
// This obviously could have been solved using sequential calls with a built-in delay,
// but that sounded less fun :)

type Request interface {
	Execute() (interface{}, error)
}

var requests chan Request

func StartRequestPool() {
	requests = make(chan Request)
	go processRequests()
}

func processRequests() error {
    for req := range requests {
		req.Execute()
	}
	return nil
}

//func
