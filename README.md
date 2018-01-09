# Spring Health Exporter

A Prometheus exporter for monitoring Spring Boot applications health status.

## Usage

Run:

	docker run -p 9117:9117 contentwisetv/spring-health-exporter
   
Get metrics via:
   
    curl http://localhost:9117/probe?target=host:port

## Prometheus configuration

The approach is similar to the one used for the Black Box Exporter, so the same configuration examples can fit: https://github.com/prometheus/blackbox_exporter#prometheus-configuration.

## Metrics

* `1`: service is up
* `0`: service is down
* `-1`: some error occurred while trying to retrieve health status (e.g.: service is unreachable or malformed response)

## Docker Hub

This image is published to [Docker Hub](https://hub.docker.com/r/contentwisetv/spring-health-exporter/) via automated build.

## License

Author: Marco Miglierina <marco.miglierina@contentwise.tv>

Licensed under the MIT License.

## Credits

This exporter was created starting from https://github.com/tolleiv/spring-ms-exporter