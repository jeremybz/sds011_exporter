# sds011_exporter
Prometheus exporter for the widely-available [sds011](http://inovafitness.com/en/a/chanpinzhongxin/95.html) air particle sensor, based on [https://github.com/ryszard/sds011](https://github.com/ryszard/sds011)

# Compiling
``` go get golang.org/x/sys/unix
go get github.com/ryszard/sds011/go/sds011
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp

env GOOS=linux GOARCH=arm GOARM=5 go build ```

# Raspberry Pi Hardware Setup
Connect the sensor to your Rpi using its USB adaptor

## Node Exporter (optional)
``` wget https://github.com/prometheus/node_exporter/releases/download/v1.0.0-rc.0/node_exporter-1.0.0-rc.0.linux-armv5.tar.gz
tar xfz node_exporter-1.0.0-rc.0.linux-armv5.tar.gz
node_exporter --collector.disable-defaults --collector.cpu --collector.cpufreq --collector.filefd --collector.hwmon --collector.ipvs --collector.loadavg --collector.meminfo --collector.netdev --collector.netstat --collector.stat --collector.time --collector.timex --collector.uname & ```

## Prometheus (optional)
Prometheus will serve metrics locally, or can send them to a remote endpoint in the cloud, such as [Victoria Metrics](https://github.com/VictoriaMetrics/VictoriaMetrics/wiki/Single-server-VictoriaMetrics).  See [prom.yml](../blob/master/prom.yml)
``` wget `curl -s https://api.github.com/repos/prometheus/prometheus/releases/latest | grep browser_download_url | perl -nle '/browser_download_url.*"(https.*linux-armv7.*.tar.gz)"/ and print $1'`
tar xfz prometheus-*.linux-armv5.tar.gz
./prometheus-2.16.0.linux-armv5/prometheus --config.file=/home/pi/prom.yml & ```
## Usage
``` Usage of ./sds011_exporter:
  -alsologtostderr
    	log to standard error as well as files
  -cycle-minutes string
    	Length of time to cycle sensor off (1-30).  0 will disable cycling, and the sensor will stream data every second.  SDS011 have an expected working life of 8000 hours, so a cycle time of 1-2 minutes is recommended (default "2")
  -http-port string
    	port to listen on for HTTP requests. (default ":9227")
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -port-path string
    	serial port path (default "/dev/ttyUSB0")
  -stderrthreshold value
    	logs at or above this threshold go to stderr
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
```

## Dashboard
Simple grafana [dashboard](https://grafana.com/grafana/dashboards/11866)

# License
[Apache 2.0](https://www.tldrlegal.com/l/apache2), please see the file [LICENSE](../blob/master/LICENSE).
