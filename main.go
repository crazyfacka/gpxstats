package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

type stats struct {
	maxSpeed      float64
	maxSpeedPoint gpx.Point

	maxUpSlope        float64
	maxUpSlopePoint   gpx.Point
	maxDownSlope      float64
	maxDownSlopePoint gpx.Point
}

/*
 * From: https://gist.github.com/cdipaolo/d3f8db3848278b49db68
 */

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func getDistance(lat1, lon1, lat2, lon2 float64) float64 {
	var la1, lo1, la2, lo2, r float64

	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in meters

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

/* === */

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

				previousPoint := q.GetFirst()
				curSpeed := point.SpeedBetween(&previousPoint, false)
				if curSpeed > st.maxSpeed {
					st.maxSpeed = curSpeed
					st.maxSpeedPoint = point.Point
				}

				distance := getDistance(point.GetLatitude(), point.GetLongitude(), previousPoint.GetLatitude(), previousPoint.GetLongitude())
				elevationDiff := point.Elevation.Value() - previousPoint.Elevation.Value()

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

func main() {
	var file *os.File
	var gpxFile *gpx.GPX
	var data []byte
	var err error

	bytesRead := -1
	buffer := make([]byte, 100)

	file, err = os.Open("test.gpx")
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
		os.Exit(-1)
	}

	defer file.Close()

	for bytesRead != 0 {
		bytesRead, err = file.Read(buffer)
		data = append(data, buffer[:bytesRead]...)
	}

	gpxFile, err = gpx.ParseBytes(data)
	if err != nil {
		fmt.Printf("Error parsing GPX: %s\n", err.Error())
		os.Exit(-1)
	}

	/* Print stats */

	st := getStats(gpxFile)

	fmt.Println("== GPX File stats ==")
	fmt.Println("")

	fmt.Printf("Name: %s\n", gpxFile.Name)
	fmt.Printf("Description: %s\n", gpxFile.Description)
	fmt.Printf("Author: %s\n", gpxFile.AuthorName)

	fmt.Println("")

	movingTime := time.Date(0, 0, 0, 0, 0, int(gpxFile.MovingData().MovingTime), 0, time.UTC)
	stoppedTime := time.Date(0, 0, 0, 0, 0, int(gpxFile.MovingData().StoppedTime), 0, time.UTC)
	fmt.Printf("Moving time: %02d:%02d:%02d\n", movingTime.Hour(), movingTime.Minute(), movingTime.Second())
	fmt.Printf("Stopped time: %02d:%02d:%02d\n", stoppedTime.Hour(), stoppedTime.Minute(), stoppedTime.Second())

	fmt.Println("")

	fmt.Printf("Minimum elevation: %.2fm\n", gpxFile.ElevationBounds().MinElevation)
	fmt.Printf("Maximum elevation: %.2fm\n", gpxFile.ElevationBounds().MaxElevation)
	fmt.Printf("Max up gradient: %.2f (%f, %f, %.2fm)\n", st.maxUpSlope, st.maxUpSlopePoint.Latitude, st.maxUpSlopePoint.Longitude, st.maxUpSlopePoint.Elevation.Value())
	fmt.Printf("Max down gradient: %.2f (%f, %f, %.2fm)\n", st.maxDownSlope, st.maxDownSlopePoint.Latitude, st.maxDownSlopePoint.Longitude, st.maxDownSlopePoint.Elevation.Value())

	fmt.Println("")

	fmt.Printf("Maximum speed: %.2f km/h (%f, %f, %.2fm)\n", st.maxSpeed/10*3.6, st.maxSpeedPoint.Latitude, st.maxSpeedPoint.Longitude, st.maxSpeedPoint.Elevation.Value())
}
