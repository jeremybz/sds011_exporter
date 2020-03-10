package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"time"
	"strconv"
	"net/http"
	"github.com/ryszard/sds011/go/sds011"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


var (
	port = flag.String("http-port", ":9227", "port to listen on for HTTP requests.")
	cycleMinutes = flag.String("cycle-minutes", "2", "Length of time to cycle sensor off (1-30).  0 will disable cycling, and the sensor will stream data every second.  SDS011 have an expected working life of 8000 hours, so a cycle time of 1-2 minutes is recommended")
	portPath = flag.String("port-path", "/dev/ttyUSB0", "serial port path")
	particleCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "particle_total",
			Help: "The total number of particles",
		},
		[]string{"micron_size"},
	)
	sensorInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sensor_info",
			Help: "Information about the sensor: device ID and Firmware",
		},
		[]string{"device_id","firmware_version"},
	)
)

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `sds011 reads data from the SDS011 sensor and sends them to stdout as CSV.
The columns are: an RFC3339 timestamp, the PM2.5 level, the PM10 level.`)
		fmt.Fprintf(os.Stderr, "\n\nUsage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	prometheus.MustRegister( particleCount )
	prometheus.MustRegister( sensorInfo )
}


func main() {
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	sensor, err := sds011.New(*portPath)
	if err != nil {
		log.Fatal(err)
	}
	defer sensor.Close()

	sensor.Awake()
	sensor.MakeActive()

	cymins, err := strconv.Atoi(*cycleMinutes)
	if err != nil {
		log.Printf("FATAL ERROR: Invalid cycle-minutes setting: %v", err)
		os.Exit(1)
	}
	sensor.SetCycle( uint8( cymins ) )

	firmware, err := sensor.Firmware()
	if err != nil {
		log.Printf("FATAL ERROR: Unable to get sensor firmware: %v", err)
		os.Exit(1)
	}

	deviceid, err := sensor.DeviceID()
	if err != nil {
		log.Printf("FATAL ERROR: Unable to get sensor device id: %v", err)
		os.Exit(1)
	}

	sensorInfo.WithLabelValues( firmware, deviceid ).Set( 1 )

	go func() {
		log.Printf("Listening on port %s", *port)
		log.Fatal(http.ListenAndServe(*port, nil))
	}()

	for {
		point, err := sensor.Get()
		if err != nil {
			log.Printf("ERROR: sensor.Get: %v", err)
			continue
		}
		fmt.Fprintf(os.Stdout, "%v,%v,%v\n", point.Timestamp.Format(time.RFC3339), point.PM25, point.PM10)
		particleCount.WithLabelValues("2.5").Add( point.PM25 )
		particleCount.WithLabelValues("10").Add( point.PM10 )
	}
}
