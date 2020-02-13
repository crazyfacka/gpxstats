FROM golang:1.13 AS build

WORKDIR /go/src/github.com/crazyfacka/gpxstats/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
WORKDIR /app
COPY --from=build /go/src/github.com/crazyfacka/gpxstats/app .
ENTRYPOINT [ "./app" ]
