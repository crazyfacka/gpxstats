# gpxstats
Stats output from GPX files

```
$ docker run -it --rm -v ./gpxdata:/app/data crazyfacka/gpxstats data/20200209_mondim_n305.gpx
== GPX File stats ==

Filename: data/20200209_mondim_n305.gpx
Name:
Description:
Author:

Moving time: 03:00:34
Stopped time: 00:48:47
Total time: 03:49:21

Minimum elevation: 30.89m
Maximum elevation: 940.77m
Max down gradient: -100.00% (0.000000, 0.000000, 0.00m) - BETA
Max up gradient: 40.83% (41.335953, -7.863357, 680.25m) - BETA

Total distance: 239.52 km
Maximum speed: 145.00 km/h (41.375006, -8.503440, 85.68m)
```

```
$ docker run -it --rm -v ./gpxdata:/app/data crazyfacka/gpxstats data/20191020_frecha_mizarela.gpx data/20191110_gilhofrei.gpx
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

You can use docker for that.

```
$ docker build -t <image_name> .
```

## Run

As seen above on the description. You have a Docker image available @ https://hub.docker.com/r/crazyfacka/gpxstats