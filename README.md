Here is the tutorial formatted for a README:

---

# Prometheus Setup and Architecture

## Resources:
- [YouTube: Nana - Architecture Setup](https://www.youtube.com/watch?v=ZtYMuxAj7EU)
- [GitHub: Prometheus Book (Chapters 13, 14)](https://github.com/appscode/books/tree/master/Prometheus)

## Prometheus Architecture:

### Client Libraries:
Client Libraries are libraries you can integrate into your applications to expose custom metrics directly from your code. These metrics are collected by Prometheus using the pull-based model.

### Exporters:
An exporter is a software component that collects and exposes metrics from a system, application, or service in a format that Prometheus can scrape and store. It runs beside your main application and obtains metrics.

### Service Discovery:
Service Discovery in Prometheus is the mechanism by which it dynamically identifies and monitors services, systems, or endpoints without requiring manual configuration. Instead of hardcoding the targets (e.g., servers or applications) to scrape, Prometheus can use service discovery to automatically update its list of scrape targets as systems are added, removed, or change.
Example: Tags in EC2, and Labels and Annotations in Kubernetes.

### Scraping:
Service discovery and relabelling give us a list of targets to be monitored. Now Prometheus needs to fetch the metrics. Prometheus does this by sending an HTTP request called a scrape.

### Storage:
Prometheus stores data locally in a custom database.

### Dashboards:
Prometheus has a number of HTTP APIs that allow you to both request raw data and evaluate PromQL queries. These can be used to produce graphs and dashboards. It is recommended that Grafana be used for dashboards. Grafana supports talking to multiple Prometheus servers, even within a single dashboard panel.

### Recording Rules and Alerts:
Recording rules are used to precompute and store frequently-used or expensive queries as new time series. This improves the performance of dashboards and alerts by reducing the computational load on Prometheus.

Alerting rules define conditions for generating alerts based on metrics. When an alerting rule's condition is met, Prometheus creates an "active alert," which is sent to an external system (e.g., PagerDuty, Slack) via the Alertmanager.

### Alert Management:
The Alertmanager receives alerts from Prometheus servers and turns them into notifications.

## Getting Started with Prometheus

### Running Prometheus:
1. Download: [Prometheus Download Page](https://prometheus.io/download/)
    ```bash
    $ tar -xzf prometheus-*.linux-amd64.tar.gz
    $ cd prometheus-*.linux-amd64/
    $ ./prometheus
    ```

2. UI: [http://localhost:9090/](http://localhost:9090/)

3. Expression Browser:
   The expression browser is useful for running ad hoc queries, developing PromQL expressions, and debugging both PromQL and the data inside Prometheus.
    - Example queries:
        - `process_resident_memory_bytes`
        - `rate(prometheus_tsdb_head_samples_appended_total[1m])`

### Running the Node Exporter:
The Node exporter exposes kernel- and machine-level metrics on Unix systems, such as Linux. It provides all the standard metrics such as CPU, memory, disk space, disk I/O, and network bandwidth.

1. Download: [Node Exporter Download Page](https://prometheus.io/download/#node_exporter)
    ```bash
    $ tar -xzf node_exporter-*.linux-amd64.tar.gz
    $ cd node_exporter-*.linux-amd64/
    $ ./node_exporter
    ```

2. UI: [http://localhost:9100/](http://localhost:9100/)

### Update Prometheus.yml File:
```yaml
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "node"
    static_configs:
      - targets:
          - localhost:9100
```

### Example Query:
```bash
rate(node_network_receive_bytes_total[1m])
```

## Alerting Setup:

### First: Adding Alerting Rules to Prometheus
Define the logic of what constitutes an alert.

### Second: The Alertmanager Converts Firing Alerts into Notifications
1. **Alertmanager Configuration:**

   Add the rules file and alerting address:
    ```yaml
    alerting:
      alertmanagers:
      - static_configs:
          - targets:
              - localhost:9093
    
    rule_files:
      # - "first_rules.yml"
      # - "second_rules.yml"
      - rules.yml
    ```

2. **Create a new file `rules.yml`:**
    ```yaml
    groups:
    - name: example
      rules:
        - alert: InstanceDown
          expr: up == 0
          for: 1m
    ```

### Install Alertmanager:
1. Download: [Alertmanager Download Page](https://prometheus.io/download/)
    ```bash
    $ tar -xzf alertmanager-*.linux-amd64.tar.gz
    $ cd alertmanager-*.linux-amd64/
    ```

2. Configure Alertmanager with your email:
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

3. Start Alertmanager:
    ```bash
    ./alertmanager
    ```

---

This README covers the Prometheus setup and architecture, along with instructions for configuring Prometheus, Node Exporter, Alertmanager, and setting up alerting.

## Metric types

The Prometheus client libraries offer four core metric types:
- Counter: Value can only increase or be reset to zero on restart. (cumulative)
- Gauge: Gauge is a Metric that represents a single numerical value that can arbitrarily go up and down.
- Histogram: 
- Summary: 