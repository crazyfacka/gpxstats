package tools

import (
	"fmt"
	"math"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

func getTimeFormattedForStats(t time.Time) string {
	timeStamp := t.Format("15:04:05")
	if t.Day() != 30 { // FIXME for when we have more than 30 days combined of resting or riding
		timeStamp = fmt.Sprintf("%dd %s", t.Day(), timeStamp)
	}

	return timeStamp
}

// PrintSingleStats prints the stats of a single GPX file
func PrintSingleStats(gpxFile *gpx.GPX, fname string) {
	st := getStats(gpxFile)

	fmt.Println("== GPX File stats ==")
	fmt.Println("")

	fmt.Printf("Filename: %s\n", fname)
	fmt.Printf("Name: %s\n", gpxFile.Name)
	fmt.Printf("Description: %s\n", gpxFile.Description)
	fmt.Printf("Author: %s\n", gpxFile.AuthorName)

	fmt.Println("")

	if !st.invalidTime {
		movingTime := time.Date(0, 0, 0, 0, 0, int(gpxFile.MovingData().MovingTime), 0, time.UTC)
		stoppedTime := time.Date(0, 0, 0, 0, 0, int(gpxFile.MovingData().StoppedTime), 0, time.UTC)
		totalTime := time.Date(0, 0, 0, 0, 0, int(gpxFile.MovingData().MovingTime+gpxFile.MovingData().StoppedTime), 0, time.UTC)
		fmt.Printf("Moving time: %s\n", getTimeFormattedForStats(movingTime))
		fmt.Printf("Stopped time: %s\n", getTimeFormattedForStats(stoppedTime))
		fmt.Printf("Total time: %s\n", getTimeFormattedForStats(totalTime))
	} else {
		fmt.Println("Time is invalid for this GPX file")
	}

	fmt.Println("")

	fmt.Printf("Minimum elevation: %.2fm\n", gpxFile.ElevationBounds().MinElevation)
	fmt.Printf("Maximum elevation: %.2fm\n", gpxFile.ElevationBounds().MaxElevation)
	fmt.Printf("Max down gradient: %.2f%% (%f, %f, %.2fm) - BETA\n", st.maxDownSlope, st.maxDownSlopePoint.Latitude, st.maxDownSlopePoint.Longitude, st.maxDownSlopePoint.Elevation.Value())
	fmt.Printf("Max up gradient: %.2f%% (%f, %f, %.2fm) - BETA\n", st.maxUpSlope, st.maxUpSlopePoint.Latitude, st.maxUpSlopePoint.Longitude, st.maxUpSlopePoint.Elevation.Value())

	fmt.Println("")

	fmt.Printf("Total distance: %.2f km\n", gpxFile.MovingData().MovingDistance/1000)
	fmt.Printf("Maximum speed: %.2f km/h (%f, %f, %.2fm)\n", st.maxSpeed, st.maxSpeedPoint.Latitude, st.maxSpeedPoint.Longitude, st.maxSpeedPoint.Elevation.Value())
}

// PrintCombinedStats prints the stats of a combination of GPX files
func PrintCombinedStats(gpxFiles []*gpx.GPX, fnames []string) {
	var st *stats
	var movingTime, stoppedTime, minElevation, maxElevation, minGradient, maxGradient, totalDistance, maxSpeed, maxStretch float64
	var minElevationFile, maxElevationFile, minGradientFile, maxGradientFile, maxSpeedFile, maxStretchFile string
	var ignoredTimeFiles []string

	minElevation = math.MaxFloat64
	maxElevation = 0
	minGradient = -100
	maxGradient = 0
	maxSpeed = 0
	maxStretch = 0

	for i, gpxFile := range gpxFiles {
		st = getStats(gpxFile)

		if !st.invalidTime {
			movingTime += gpxFile.MovingData().MovingTime
			stoppedTime += gpxFile.MovingData().StoppedTime
		} else {
			ignoredTimeFiles = append(ignoredTimeFiles, fnames[i])
		}

		if minElevation > gpxFile.ElevationBounds().MinElevation {
			minElevation = gpxFile.ElevationBounds().MinElevation
			minElevationFile = fnames[i]
		}

		if maxElevation < gpxFile.ElevationBounds().MaxElevation {
			maxElevation = gpxFile.ElevationBounds().MaxElevation
			maxElevationFile = fnames[i]
		}

		if minGradient > st.maxDownSlope {
			minGradient = st.maxDownSlope
			minGradientFile = fnames[i]
		}

		if maxGradient < st.maxUpSlope {
			maxGradient = st.maxUpSlope
			maxGradientFile = fnames[i]
		}

		totalDistance += gpxFile.MovingData().MovingDistance

		if maxSpeed < st.maxSpeed {
			maxSpeed = st.maxSpeed
			maxSpeedFile = fnames[i]
		}

		if gpxFile.MovingData().MovingDistance > maxStretch {
			maxStretch = gpxFile.MovingData().MovingDistance
			maxStretchFile = fnames[i]
		}
	}

	fmt.Println("== GPX combined stats ==")
	fmt.Println("")

	movingTimeTime := time.Date(0, 0, 0, 0, 0, int(movingTime), 0, time.UTC)
	stoppedTimeTime := time.Date(0, 0, 0, 0, 0, int(stoppedTime), 0, time.UTC)
	totalTime := time.Date(0, 0, 0, 0, 0, int(movingTime+stoppedTime), 0, time.UTC)
	fmt.Printf("Moving time: %s\n", getTimeFormattedForStats(movingTimeTime))
	fmt.Printf("Stopped time: %s\n", getTimeFormattedForStats(stoppedTimeTime))
	fmt.Printf("Total time: %s\n", getTimeFormattedForStats(totalTime))

	if len(ignoredTimeFiles) > 0 {
		fmt.Println("")
		fmt.Printf("Files with erroneous timestamps (ignored in the counting): %+v\n", ignoredTimeFiles)
	}

	fmt.Println("")

	fmt.Printf("Minimum elevation: %.2fm (%s)\n", minElevation, minElevationFile)
	fmt.Printf("Maximum elevation: %.2fm (%s)\n", maxElevation, maxElevationFile)
	fmt.Printf("Max down gradient: %.2f%% (%s) - BETA\n", minGradient, minGradientFile)
	fmt.Printf("Max up gradient: %.2f%% (%s) - BETA\n", maxGradient, maxGradientFile)

	fmt.Println("")

	fmt.Printf("Total distance: %.2f km\n", totalDistance/1000)
	fmt.Printf("Max stretch: %.2f km (%s)\n", maxStretch/1000, maxStretchFile)
	fmt.Printf("Maximum speed: %.2f km/h (%s)\n", maxSpeed, maxSpeedFile)
}
