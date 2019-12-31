package tools

import (
	"math"

	"github.com/tkrajina/gpxgo/gpx"
)

type stats struct {
	invalidTime bool

	maxSpeed      float64
	maxSpeedPoint gpx.Point

	maxUpSlope        float64
	maxUpSlopePoint   gpx.Point
	maxDownSlope      float64
	maxDownSlopePoint gpx.Point
}

func getStats(gpxFile *gpx.GPX) *stats {
	var st *stats
	var q *Queue

	st = &stats{
		maxSpeed:     0,
		maxUpSlope:   0,
		maxDownSlope: -100,
	}

	q = NewQueue(3)

	for _, track := range gpxFile.Tracks {
		for _, segment := range track.Segments {
			for i, point := range segment.Points {
				if i == 0 {
					q.Push(point)
					continue
				}

				lastPoint := q.GetLast()

				if point.Timestamp.Unix() < 0 {
					st.invalidTime = true
				}

				curSpeed := getSpeed(q.GetArray())
				if curSpeed > st.maxSpeed {
					st.maxSpeed = curSpeed
					st.maxSpeedPoint = point.Point
				}

				distance := getDistance(point.GetLatitude(), point.GetLongitude(), lastPoint.GetLatitude(), lastPoint.GetLongitude())
				elevationDiff := getElevationDiff(point, q.GetArray())

				if distance > 0 && elevationDiff != 0 && math.Abs(elevationDiff) < distance {
					if elevationDiff < 0 { // Down slope
						slope := 100 * elevationDiff / distance
						if st.maxDownSlope > slope {
							st.maxDownSlope = slope
							st.maxDownSlopePoint = point.Point
						}
					} else { // Up slope
						slope := 100 * elevationDiff / distance
						if slope > st.maxUpSlope {
							st.maxUpSlope = slope
							st.maxUpSlopePoint = point.Point
						}
					}
				}

				q.Push(point)
			}
		}
	}

	return st
}
