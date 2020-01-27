# HUIC

An HTTP-UDP-Internet-Connection Protocol for downloading large files over large distances.

## About

HUIC is meant to be used in the special case where there are really large files to be downloaded over the internet.  It can download files 100's of times faster than HTTP.

This is an open-source alternative to [Aspera](https://www.ibm.com/cloud/high-speed-data-transfer) and a modern version of [Tsunami UDP](http://tsunami-udp.sourceforge.net/).

## Serving Files

```
docker run -d \
    -v /my/folder:/huic \
    -p 8080:8080 \
    -p 41000-41999:41000-41999 \
    millisecond/huic-server:latest
```

This will serve `/my/folder` over HTTP and HUIC.

### Securing Transfer

If you don't have another service (proxy, nginx, etc) terminating HTTPS in front of HUIC, you can start HUIC directly with a TLS certificate and key:

```
docker run -d \
    -v /my/folder:/huic \
    -p 8080:8080 \
    -p 41000-41999:41000-41999 \
    -e TLS_CERT=/my/secure/tls.cert \
    -e TLS_KEY=/my/secure/tls.key \
    millisecond/huic-server:latest
```

## Downloading Files

### Client

go get github.com/millisecond/huic/client

Or get the pre-compiled binaries for your platform on the [releases page](https://github.com/millisecond/huic/releases)

### Usage

```
huic get https://your.docker.domain/file
```

### Benchmarking

As performance is the main reason to use HUIC we have built-in tools to benchmark against HTTP.  This can also be used to determine benefit for your specific use-case/configuration.

```
huic benchmark https://your.docker.domain/file
```

NOTE: this will download the first 100MB of the file over both HTTP and HUIC using 200MB of bandwidth.
