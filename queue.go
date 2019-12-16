package main

import "github.com/tkrajina/gpxgo/gpx"

type gpxQueuePoint struct {
	point    gpx.GPXPoint
	next     *gpxQueuePoint
	previous *gpxQueuePoint
}

// Queue main struct
type Queue struct {
	head *gpxQueuePoint
	tail *gpxQueuePoint

	length   int
	elements int
}

// Push a new element to the queue
func (q *Queue) Push(p gpx.GPXPoint) {
	np := &gpxQueuePoint{
		point: p,
		next:  q.head,
	}

	if q.elements != 0 {
		q.head.previous = np
	} else {
		q.tail = np
	}
	q.head = np

	q.elements++

	if q.elements > q.length {
		q.tail = q.tail.previous
		q.tail.next = nil
		q.elements--
	}
}

// GetFirst returns the last added point
func (q *Queue) GetFirst() gpx.GPXPoint {
	return q.head.point
}

// GetLast returns the tail of the queue
func (q *Queue) GetLast() gpx.GPXPoint {
	return q.tail.point
}

// GetArray returns the queue as an array
func (q *Queue) GetArray() []gpx.GPXPoint {
	data := make([]gpx.GPXPoint, q.elements)

	p := q.head
	for i := 0; p != nil; i++ {
		data[i] = p.point
		p = p.next
	}

	return data
}

// NewQueue creates a new GPX Points queue
func NewQueue(len int) *Queue {
	return &Queue{
		length:   len,
		elements: 0,
	}
}
