Speed Wi-Fi Home Exporter
=========================

The prometheus exporter for [Speed Wi-Fi HOME L02](https://www.au.com/mobile/product/data/hws33/).


## Usage

### run on local

``` shell
$ go get github.com/shidax-tech/speed-wifi-home-exporter
$ speed-wifi-home-exporter -listen=localhost:9999
```

And, access to [http://localhost:9999/metrics](http://localhost:9999/metrics)

### run on docker

``` shell
$ git clone https://github.com/shidax-tech/speed-wifi-home-exporter && cd speed-wifi-home-exporter
$ docker build -t speed-wifi-home-exporter:latest .
$ docker run -p 9999:9999 speed-wifi-home-exporter
```
