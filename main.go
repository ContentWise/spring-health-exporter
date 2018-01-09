package main

import (
	"flag"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strings"
)

var (
	addr = flag.String("listen-address", ":9117", "The address to listen on for HTTP requests.")
)

type ServiceHealth struct {
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>SpringBoot Health Exporter</title></head>
            <body>
            <h1>SpringBoot Health Exporter</h1>
            <p><a href="/probe">Run a probe</a></p>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})
	flag.Parse()
	http.HandleFunc("/probe", probeHandler)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", 400)
		return
	}
	healthGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "service_health",
		Help: "Displays whether or not the service is healthy",
	})

	registry := prometheus.NewRegistry()
	registry.MustRegister(healthGauge)
	metric := 1.0

	bytes, err := getJson("http://" + target + "/health")
	if err != nil {
		metric = -1.0
	} else {
		health := ServiceHealth{}
		err = json.Unmarshal([]byte(bytes), &health)
		if err != nil {
			metric = -1.0
		} else if strings.ToUpper(health.Status) != "UP" {
			metric = 0.0
		}
	}

	healthGauge.Set(metric)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func getJson(target string) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
