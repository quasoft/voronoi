package goalgorithms

// Priority queue adapted from the examples of the container/heap package.

import (
	"container/heap"
	"sort"
)

type EventType int

const (
	EventSite   EventType = 0
	EventCircle EventType = 1
)

// Event represents a site or circle event.
type Event struct {
	Site Site
	// The index is needed by update and is maintained by the heap.Interface methods.
	index     int // The index of the event in the heap.
	EventType EventType
}

// A EventQueue is a priority queue that implements heap.Interface and holds Events.
type EventQueue []*Event

// Creates a new event queue and initializes it with events for the given list of sites.
func NewEventQueue(sites SiteSlice) EventQueue {
	sort.Sort(sites)

	eventQueue := make(EventQueue, len(sites))
	i := 0
	for _, site := range sites {
		eventQueue[i] = &Event{
			Site:  site,
			index: i,
		}
		i++
	}
	heap.Init(&eventQueue)
	return eventQueue
}

func (pq EventQueue) Len() int { return len(pq) }

func (pq EventQueue) Less(i, j int) bool {
	// We want Pop to give us the event with highest 'y' position.
	return pq[i].Site.Y < pq[j].Site.Y
}

func (pq EventQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *EventQueue) Push(x interface{}) {
	n := len(*pq)
	event := x.(*Event)
	event.index = n
	*pq = append(*pq, event)
}

func (pq *EventQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	event := old[n-1]
	event.index = -1 // for safety
	*pq = old[0 : n-1]
	return event
}

// Remove removes the element with the specified index from the queue.
func (pq *EventQueue) Remove(event *Event) {
	heap.Remove(pq, event.index)
}
