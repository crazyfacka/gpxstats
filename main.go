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
				lastPoint := q.GetLast()

				curSpeed := point.SpeedBetween(&previousPoint, false)
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

func printSingleStats(gpxFile *gpx.GPX) {
	st := getStats(gpxFile)

	fmt.Println("== GPX File stats ==")
	fmt.Println("")

	fmt.Printf("Name: %s\n", gpxFile.Name)
	fmt.Printf("Description: %s\n", gpxFile.Description)
	fmt.Printf("Author: %s\n", gpxFile.AuthorName)

	fmt.Println("")

	movingTime := time.Date(0, 0, 0, 0, 0, int(gpxFile.MovingData().MovingTime), 0, time.UTC)
	stoppedTime := time.Date(0, 0, 0, 0, 0, int(gpxFile.MovingData().StoppedTime), 0, time.UTC)
	totalTime := time.Date(0, 0, 0, 0, 0, int(gpxFile.MovingData().MovingTime+gpxFile.MovingData().StoppedTime), 0, time.UTC)
	fmt.Printf("Moving time: %02d:%02d:%02d\n", movingTime.Hour(), movingTime.Minute(), movingTime.Second())
	fmt.Printf("Stopped time: %02d:%02d:%02d\n", stoppedTime.Hour(), stoppedTime.Minute(), stoppedTime.Second())
	fmt.Printf("Total time: %02d:%02d:%02d\n", totalTime.Hour(), totalTime.Minute(), totalTime.Second())

	fmt.Println("")

	fmt.Printf("Minimum elevation: %.2fm\n", gpxFile.ElevationBounds().MinElevation)
	fmt.Printf("Maximum elevation: %.2fm\n", gpxFile.ElevationBounds().MaxElevation)
	fmt.Printf("Max down gradient: %.2f (%f, %f, %.2fm) - BETA\n", st.maxDownSlope, st.maxDownSlopePoint.Latitude, st.maxDownSlopePoint.Longitude, st.maxDownSlopePoint.Elevation.Value())
	fmt.Printf("Max up gradient: %.2f (%f, %f, %.2fm) - BETA\n", st.maxUpSlope, st.maxUpSlopePoint.Latitude, st.maxUpSlopePoint.Longitude, st.maxUpSlopePoint.Elevation.Value())

	fmt.Println("")

	fmt.Printf("Total distance: %.2f km\n", gpxFile.MovingData().MovingDistance/1000)
	fmt.Printf("Maximum speed: %.2f km/h (%f, %f, %.2fm)\n", st.maxSpeed/10*3.6, st.maxSpeedPoint.Latitude, st.maxSpeedPoint.Longitude, st.maxSpeedPoint.Elevation.Value())
}

func printCombinedStats(gpxFiles []*gpx.GPX) {
	var st *stats
	var movingTime, stoppedTime, minElevation, maxElevation, minGradient, maxGradient, totalDistance, maxSpeed float64

	minElevation = math.MaxFloat64
	maxElevation = 0
	minGradient = -100
	maxGradient = 0
	maxSpeed = 0

	for _, gpxFile := range gpxFiles {
		st = getStats(gpxFile)

		movingTime += gpxFile.MovingData().MovingTime
		stoppedTime += gpxFile.MovingData().StoppedTime

		if minElevation > gpxFile.ElevationBounds().MinElevation {
			minElevation = gpxFile.ElevationBounds().MinElevation
		}

		if maxElevation < gpxFile.ElevationBounds().MaxElevation {
			maxElevation = gpxFile.ElevationBounds().MaxElevation
		}

		if minGradient > st.maxDownSlope {
			minGradient = st.maxDownSlope
		}

		if maxGradient < st.maxUpSlope {
			maxGradient = st.maxUpSlope
		}

		totalDistance += gpxFile.MovingData().MovingDistance

		if maxSpeed < st.maxSpeed {
			maxSpeed = st.maxSpeed
		}
	}

	fmt.Println("== GPX combined stats ==")
	fmt.Println("")

	movingTimeTime := time.Date(0, 0, 0, 0, 0, int(movingTime), 0, time.UTC)
	stoppedTimeTime := time.Date(0, 0, 0, 0, 0, int(stoppedTime), 0, time.UTC)
	totalTime := time.Date(0, 0, 0, 0, 0, int(movingTime+stoppedTime), 0, time.UTC)
	fmt.Printf("Moving time: %02d:%02d:%02d\n", movingTimeTime.Hour(), movingTimeTime.Minute(), movingTimeTime.Second())
	fmt.Printf("Stopped time: %02d:%02d:%02d\n", stoppedTimeTime.Hour(), stoppedTimeTime.Minute(), stoppedTimeTime.Second())
	fmt.Printf("Total time: %02d:%02d:%02d\n", totalTime.Hour(), totalTime.Minute(), totalTime.Second())

	fmt.Println("")

	fmt.Printf("Minimum elevation: %.2fm\n", minElevation)
	fmt.Printf("Maximum elevation: %.2fm\n", maxElevation)
	fmt.Printf("Max down gradient: %.2f - BETA\n", minGradient)
	fmt.Printf("Max up gradient: %.2f - BETA\n", maxGradient)

	fmt.Println("")

	fmt.Printf("Total distance: %.2f km\n", totalDistance/1000)
	fmt.Printf("Maximum speed: %.2f km/h\n", maxSpeed/10*3.6)
}

func main() {
	var file *os.File
	var gpxFiles []*gpx.GPX
	var data []byte
	var err error

	bytesRead := -1
	buffer := make([]byte, 100)

	gpxFilesToRead := os.Args[1:]
	if len(gpxFilesToRead) < 1 {
		fmt.Println("You must specify at least one GPX file to parse")
		os.Exit(-1)
	}

	gpxFiles = make([]*gpx.GPX, len(gpxFilesToRead))

	for i, curFileToRead := range gpxFilesToRead {
		file, err = os.Open(curFileToRead)
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err.Error())
			os.Exit(-1)
		}

		for bytesRead != 0 {
			bytesRead, err = file.Read(buffer)
			data = append(data, buffer[:bytesRead]...)
		}

		gpxFiles[i], err = gpx.ParseBytes(data)
		if err != nil {
			fmt.Printf("Error parsing GPX: %s\n", err.Error())
			os.Exit(-1)
		}

		file.Close()

		bytesRead = -1
		buffer = make([]byte, 100)
		data = data[:0]
	}

	/* Print stats */
	if len(gpxFiles) == 1 {
		printSingleStats(gpxFiles[0])
	} else {
		printCombinedStats(gpxFiles)
	}
}
