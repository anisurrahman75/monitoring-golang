# Prometheus

## Overview
Prometheus is an open-source systems monitoring and alerting toolkit. It uses a pull-based model to scrape metrics from configured targets and stores them in a time-series database.

---

## Prometheus Architecture

### Components
1. **Client Libraries**:  
   These are libraries integrated into your application to expose custom metrics directly from your code. Prometheus scrapes these metrics using its pull-based model.

2. **Exporters**:  
   Exporters are software components that collect and expose metrics from systems, applications, or services in a format Prometheus can scrape. For example:
    - **Node Exporter**: Exposes hardware and OS-level metrics (e.g., CPU, memory, disk usage).
    - **Custom Exporters**: Can be built for specific applications.

3. **Service Discovery**:  
   Service Discovery dynamically identifies and monitors services without manual configuration. Examples include:
    - Tags in AWS EC2.
    - Labels and annotations in Kubernetes.

4. **Scraping**:  
   Prometheus periodically sends HTTP requests (scrapes) to configured targets to fetch metrics.

5. **Storage**:  
   Metrics are stored locally in Prometheus' custom time-series database.

6. **Dashboards**:  
   Prometheus provides HTTP APIs for querying and visualizing data. However, **Grafana** is recommended for creating advanced dashboards.

7. **Alerting**:  
   Alerts are defined using PromQL (Prometheus Query Language). When conditions are met, alerts are sent to the **Alertmanager**, which processes and routes notifications (e.g., emails, Slack messages).

8. **Recording Rules**:  
   Precompute frequently-used or expensive queries and store them as new time series to improve performance.

---


## Getting Started with Prometheus

### Docker Compose File
Create a `docker-compose.yml` file with the following content:

```yaml
services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    restart: unless-stopped
    volumes:
      - ./linux-ymls/prometheus.yml:/etc/prometheus/prometheus.yml
      - ./linux-ymls/rules.yml:/etc/prometheus/rules.yml
      - prometheus_data:/prometheus
    command:
      - --config.file=/etc/prometheus/prometheus.yml
      - --storage.tsdb.path=/prometheus
      - --web.console.libraries=/etc/prometheus/console_libraries
      - --web.console.templates=/etc/prometheus/consoles
      - --web.enable-lifecycle
    network_mode: host
  node_exporter:
    image: quay.io/prometheus/node-exporter:latest
    container_name: node_exporter
    command:
      - '--path.rootfs=/host'
    network_mode: host
    pid: host
    restart: unless-stopped
    volumes:
      - node_exporter_data:/host:ro

  alertmanager:
    image: prom/alertmanager:latest
    container_name: alertmanager
    restart: unless-stopped
    volumes:
      - ./linux-ymls/alertmanager.yml:/etc/alertmanager/alertmanager.yml
    command:
      - --config.file=/etc/alertmanager/alertmanager.yml
    network_mode: host

volumes:
  prometheus_data:
  node_exporter_data:
```

### Running the Stack
```bash
docker-compose up -d
```

- **Prometheus UI**: [http://localhost:9090](http://localhost:9090)
- **Node Exporter UI**: [http://localhost:9100](http://localhost:9100)
- **Alertmanager UI**: [http://localhost:9093](http://localhost:9093)

---

## Configuring Prometheus

### `prometheus.yml`
Define scrape configurations in `prometheus.yml`. Example:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets:
            - "localhost:9093"

rule_files:
  - "/etc/prometheus/rules.yml"

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'node_exporter'
    static_configs:
      - targets: ['localhost:9100']

  - job_name: 'go_app'
    static_configs:
      - targets: ['localhost:8000']  # Go application's metrics endpoint
```

---

## Alerting Setup

### `alertmanager.yml`
Configure Alertmanager to send notifications via email or other channels. Example:

```yaml
global:
  smtp_smarthost: 'localhost:25'  # Replace with your SMTP server and port
  smtp_from: 'anisur@appscode.com'

route:
  receiver: 'email-alert'

receivers:
  - name: 'email-alert'
    email_configs:
      - to: 'anisur@appscode.com'  # Replace with the recipient's email address

inhibit_rules:
  - source_match:
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['alertname', 'dev', 'instance']
```

### Defining Alerts (`rules.yml`)
Define alerting rules in `rules.yml`. Example:

```yaml
groups:
  - name: instance_down
    rules:
      - alert: InstanceDown
        expr: up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Instance {{ $labels.instance }} is down"
          description: "Prometheus target {{ $labels.instance }} is unreachable for more than 1 minute."
  - name: go_app_alerts
    rules:
      - alert: HighRequestLatency
        expr: hello_world_latency_seconds > 1
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "High request latency detected"
          description: "Requests to Go application are taking longer than 1 second."
      - alert: HighErrorRate
        expr: rate(hello_world_exceptions_total[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate in Go application"
          description: "More than 10% of requests are failing in the last 5 minutes."
```

Start Prometheus to apply changes:

```bash
docker-compose restart prometheus
```

## Metric Types

Prometheus supports four core metric types:

### 1. Counter
A cumulative metric that only increases or resets to zero on restart.  
**Use Case**: Counting events like HTTP requests.  
**Example**: Total number of HTTP requests served by a server.

```go
var httpRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests received",
}, []string{"status", "path", "method"})

func httpRequestCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		method := r.Method
		path := r.URL.Path
		status := strconv.Itoa(recorder.statusCode)
		next.ServeHTTP(recorder, r)
		httpRequestCounter.WithLabelValues(status, path, method).Inc()
	})
}
```

To test this, you can send a steady stream of requests using `wrk`:
```bash
$ wrk -t 1 -c 1 -d 300s --latency "http://localhost:8000"
```

Metrics output:
```text
# HELP http_requests_total Total number of HTTP requests received
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/",status="200"} 6
http_requests_total{method="GET",path="/favicon.ico",status="200"} 4
http_requests_total{method="GET",path="/metrics",status="200"} 4
```

---

### 2. Gauge
Represents a single numerical value that can increase or decrease.  
**Use Case**: Tracking values like memory usage, CPU load, or active connections.  
**Example**: Number of active HTTP connections.

```go
var activeRequestsGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "http_active_requests",
		Help: "Number of active connections to the service",
	},
)

func activeRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activeRequestsGauge.Inc()
		defer activeRequestsGauge.Dec()
		// Simulate processing time
		delay := time.Duration(rand.Intn(900)) * time.Millisecond
		time.Sleep(delay)
		next.ServeHTTP(w, r)
	})
}
```

Test with `wrk`:
```bash
$ wrk -t 10 -c 400 -d 5m --latency "http://localhost:8000"
```

Metrics output:
```text
# HELP http_active_requests Number of active users in the system
# TYPE http_active_requests gauge
http_active_requests 407
```

---

### 3. Histogram
Tracks the distribution of values, such as request durations or response sizes.  
**Use Case**: Analyzing latencies or request sizes.  
**Example**: Database query latencies.

```go
var dbQueryHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "db_query_duration_seconds",
		Help:    "Histogram of database query latencies in seconds.",
		Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10}, // Custom buckets
	},
	[]string{"query_type"},
)

func dbQuery(w http.ResponseWriter, r *http.Request) {
	executeDBQuery := func(queryType string) {
		start := time.Now()
		delay := time.Duration(rand.Intn(900)) * time.Millisecond
		time.Sleep(delay)
		dbQueryHistogram.WithLabelValues(queryType).Observe(time.Since(start).Seconds())
	}

	queryTypes := []string{"SELECT"} // "INSERT", "UPDATE", "DELETE"
	randomQuery := queryTypes[rand.Intn(len(queryTypes))]
	executeDBQuery(randomQuery)
	w.Write([]byte(fmt.Sprintf("Executed %s query\n", randomQuery)))
}
```

Test with `wrk`:
```bash
$ wrk -t 10 -c 400 -d 5m --latency "http://localhost:8000/db-query"
```

Metrics output:
```text
# HELP db_query_duration_seconds Histogram of database query latencies in seconds.
# TYPE db_query_duration_seconds histogram
db_query_duration_seconds_bucket{query_type="SELECT",le="0.1"} 2665
db_query_duration_seconds_bucket{query_type="SELECT",le="0.5"} 13438
db_query_duration_seconds_bucket{query_type="SELECT",le="1"} 24117
db_query_duration_seconds_bucket{query_type="SELECT",le="2.5"} 24117
db_query_duration_seconds_bucket{query_type="SELECT",le="5"} 24117
db_query_duration_seconds_bucket{query_type="SELECT",le="10"} 24117
db_query_duration_seconds_bucket{query_type="SELECT",le="+Inf"} 24117
db_query_duration_seconds_sum{query_type="SELECT"} 10854.228629013
db_query_duration_seconds_count{query_type="SELECT"} 24117
```

---

### 4. Summary
Similar to histograms but calculates quantiles over a sliding time window.  
**Use Case**: Precise latency monitoring.  
**Example**: HTTP request durations.

```go
var httpRequestSummary = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "http_request_duration_summary_seconds",
		Help:       "Summary of HTTP request durations in seconds.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)

func summary(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()
		httpRequestSummary.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}
```

Test with `wrk`:
```bash
$ wrk -t 10 -c 400 -d 5m --latency "http://localhost:8000/db-query"
```

Metrics output:
```text
# HELP http_request_duration_summary_seconds Summary of HTTP request durations in seconds.
# TYPE http_request_duration_summary_seconds summary
http_request_duration_summary_seconds{method="GET",path="/db-query",quantile="0.5"} NaN
http_request_duration_summary_seconds{method="GET",path="/db-query",quantile="0.9"} NaN
http_request_duration_summary_seconds{method="GET",path="/db-query",quantile="0.99"} NaN
http_request_duration_summary_seconds_sum{method="GET",path="/db-query"} 10854.562146928047
http_request_duration_summary_seconds_count{method="GET",path="/db-query"} 24117
http_request_duration_summary_seconds{method="GET",path="/favicon.ico",quantile="0.5"} NaN
http_request_duration_summary_seconds{method="GET",path="/favicon.ico",quantile="0.9"} NaN
http_request_duration_summary_seconds{method="GET",path="/favicon.ico",quantile="0.99"} NaN
http_request_duration_summary_seconds_sum{method="GET",path="/favicon.ico"} 1.2524e-05
http_request_duration_summary_seconds_count{method="GET",path="/favicon.ico"} 1
http_request_duration_summary_seconds{method="GET",path="/metrics",quantile="0.5"} 0.000513413
http_request_duration_summary_seconds{method="GET",path="/metrics",quantile="0.9"} 0.001450119
http_request_duration_summary_seconds{method="GET",path="/metrics",quantile="0.99"} 0.005328081
http_request_duration_summary_seconds_sum{method="GET",path="/metrics"} 0.20834717599999988
http_request_duration_summary_seconds_count{method="GET",path="/metrics"} 281
```


## Grafana with prometheus