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
			site:  site,
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
	return pq[i].site.Y < pq[j].site.Y
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

// This example creates a EventQueue with some events, adds and manipulates an event,
// and then removes the events in priority order.
/*
func main() {
	// Some events and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(EventQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.value, 5)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}
*/
