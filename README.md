# gpxstats
Stats output from GPX files

```
.\gpxstats.exe .\20191020_frecha_mizarela.gpx
== GPX File stats ==

Filename: .\20191020_frecha_mizarela.gpx
Name:
Description:
Author:

Moving time: 04:49:02
Stopped time: 05:21:54
Total time: 10:10:56

Minimum elevation: -9.97m
Maximum elevation: 1074.87m
Max down gradient: -100.00% (0.000000, 0.000000, 0.00m) - BETA
Max up gradient: 62.67% (41.086293, -8.291996, 34.73m) - BETA

Total distance: 183.97 km
Maximum speed: 121.00 km/h (41.183153, -8.645468, 92.89m)
```

```
.\gpxstats.exe .\20191020_frecha_mizarela.gpx .\20191110_gilhofrei.gpx
== GPX combined stats ==

Moving time: 08:07:38
Stopped time: 07:31:12
Total time: 15:38:50

Minimum elevation: -9.97m (.\20191020_frecha_mizarela.gpx)
Maximum elevation: 1074.87m (.\20191020_frecha_mizarela.gpx)
Max down gradient: -100.00% () - BETA
Max up gradient: 62.67% (.\20191020_frecha_mizarela.gpx) - BETA

Total distance: 353.50 km
Max stretch: 183.97 km (.\20191020_frecha_mizarela.gpx)
Maximum speed: 131.00 km/h (.\20191110_gilhofrei.gpx)
```

## Build

```
$ go build
```

This will fetch all the necessary dependencies. Currently the only external dependency is:

* github.com/tkrajina/gpxgo


## Run

As seen above on the description