package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"time"
	"net/http"
	"github.com/ryszard/sds011/go/sds011"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":9227", "The address to listen on for HTTP requests.")
var (
	portPath = flag.String("port_path", "/dev/ttyUSB0", "serial port path")
)

var (
	particleCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "particle_total",
			Help: "The total number of particles",
		},
		[]string{"micron_size"},
	)
)

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr,
			`sds011 reads data from the SDS011 sensor and sends them to stdout as CSV.
The columns are: an RFC3339 timestamp, the PM2.5 level, the PM10 level.`)
		fmt.Fprintf(os.Stderr, "\n\nUsage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	prometheus.MustRegister( particleCount )
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

	go func() {
		//log.Fprintf(os.Stdout, "Listening")
		log.Fatal(http.ListenAndServe(*addr, nil))
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
