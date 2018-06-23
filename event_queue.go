package voronoi

// Priority queue adapted from the examples of the container/heap package.

import (
	"container/heap"
	"fmt"
	"sort"
)

// EventType represent the type of the event - either site or circle event.
type EventType int

const (
	EventSite   EventType = 0
	EventCircle EventType = 1
)

// Event represents a site or circle event.
type Event struct {
	X, Y      int       // X and Y of the site, or X and Y of the bottom point of the circle.
	index     int       // The index in the slice. Maintained by heap.Interface methods. Needed by Remove method.
	EventType EventType // The type of the event. Site = 0 and Circle = 1.
	Site      *Site     // Pointer to the related site. Only relevant for site events.
	Node      *Node     // The related arc node. Only relevant for circle events.
	Radius    int       // Radius of the circle.
}

// A EventQueue is a priority queue that implements heap.Interface and holds Events.
type EventQueue []*Event

// NewEventQueue creates a new queue and initializes it with events for the given list of sites.
func NewEventQueue(sites SiteSlice) EventQueue {
	sort.Sort(sites)

	eventQueue := make(EventQueue, len(sites))
	for i := 0; i < len(sites); i++ {
		site := &sites[i]
		eventQueue[i] = &Event{
			EventType: EventSite,
			Site:      site,
			X:         site.X,
			Y:         site.Y,
			index:     i,
		}
	}
	heap.Init(&eventQueue)
	return eventQueue
}

func (pq EventQueue) String() string {
	s := ""
	for i, event := range pq {
		prefix := ""
		if event.EventType == EventCircle {
			prefix = "C"
		} else {
			prefix = "S"
		}

		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("{%d#%s %d,%d}", event.index, prefix, event.X, event.Y)
	}
	return "{" + s + "}"
}

// Len returns the number of events in the queue.
func (pq EventQueue) Len() int { return len(pq) }

// Less compares two events and is needed as implementation of the Sort interface.
func (pq EventQueue) Less(i, j int) bool {
	// We want Pop to give us the event with highest 'y' position.
	return pq[i].Y < pq[j].Y || (pq[i].Y == pq[j].Y && pq[i].X < pq[j].X)
}

// Swap swaps two events, updating their index in the slice.
func (pq EventQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push appends an item to the queue and reorders items if necessary.
func (pq *EventQueue) Push(x interface{}) {
	n := len(*pq)
	event := x.(*Event)
	event.index = n
	*pq = append(*pq, event)
	heap.Fix(pq, n)
}

// Pop removes the last element from the queue and sets its index to -1.
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
