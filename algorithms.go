package main

import (
	"math"

	"github.com/tkrajina/gpxgo/gpx"
)

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

/*
 * To be replaced with geometric median
 */

func getMean(points []gpx.GPXPoint) float64 {
	var total float64
	total = 0.0

	for _, v := range points {
		total += v.Elevation.Value()
	}

	return math.Round(total / float64(len(points)))
}

func getElevationDiff(point gpx.GPXPoint, points []gpx.GPXPoint) float64 {
	leftMean := getMean(append([]gpx.GPXPoint{point}, points[0:len(points)-1]...))
	rightMean := getMean(points)

	return leftMean - rightMean
}

/* === */
