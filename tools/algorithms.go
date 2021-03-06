package tools

import (
	"math"

	"github.com/tkrajina/gpxgo/gpx"
)

func getMean(values []float64) float64 {
	var total float64
	total = 0.0

	for _, v := range values {
		total += v
	}

	return math.Round(total / float64(len(values)))
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

func getElevationDiff(point gpx.GPXPoint, points []gpx.GPXPoint) float64 {
	var elevationP float64
	var elevationPoints []float64

	elevationP = point.Elevation.Value()
	elevationPoints = make([]float64, len(points))
	for i, p := range points {
		elevationPoints[i] = p.Elevation.Value()
	}

	leftMean := getMean(append([]float64{elevationP}, elevationPoints[0:len(points)-1]...))
	rightMean := getMean(elevationPoints)

	return leftMean - rightMean
}

func getSpeed(points []gpx.GPXPoint) float64 {
	var previousPoint gpx.GPXPoint
	var speedPoints []float64
	var invalidPoint []bool
	var invalidPointCounter int

	speedPoints = make([]float64, len(points)-1)
	invalidPoint = make([]bool, len(points)-1)

	for i, p := range points {
		if i == 0 {
			previousPoint = p
			continue
		}

		distance := getDistance(previousPoint.GetLatitude(), previousPoint.GetLongitude(), p.GetLatitude(), p.GetLongitude())
		timeDiff := previousPoint.Timestamp.Sub(p.Timestamp)

		if distance <= 0 || timeDiff <= 0 {
			invalidPoint[i-1] = true
			invalidPointCounter++
			continue
		}

		invalidPoint[i-1] = false
		speedPoints[i-1] = distance / timeDiff.Seconds() * 3.6
	}

	if invalidPointCounter < len(speedPoints) && invalidPointCounter > 0 {
		newSpeedPoints := make([]float64, len(points)-1-invalidPointCounter)
		for i, j := 0, 0; i < len(speedPoints); i++ {
			if invalidPoint[i] {
				continue
			}
			newSpeedPoints[j] = speedPoints[i]
			j++
		}
	}

	return getMean(speedPoints)
}
