# HUIC Design

## HTTP Extensions

### Request Header: `Accept-Encoding`

A `huic` client makes its initial HTTP request with an additional header:

```
Accept-Encoding: huic
```

This notifies the server that the client can support download huic files and causes the server to start a `huic session`.

Follow-on requests (and some responses) with HUIC-specific methods use some HUIC-specific headers for their data.

`HUIC-SessionID`: a GUID representing this download session, used later to post status.

`HUIC-SessionKey`: a symmetric encryption key to protect the bytes sent over the UDP connections.

`HUIC-FileSizeBytes`: total size of the file to download.

`HUIC-UDPPortRange`: valid ports for the client to use for UDP data connections.

### HTTP Methods

`GET`: the standard HTTP get is sent with `Accept-Encoding: huic` which starts a HUIC download session.

During download the original URL will also be requested additional times with custom HTTP methods to communicate download state to the client.

`HUICSTATUS`

`HUICSTOP`

## Testing

```
docker build . -f Dockerfile -t huic-server \
    && docker run \
        -e DEBUG=true \
        -p 8080:8080 \
        -p 41000-41999:41000-41999 \
        huic-server
```

## References

[WAN Latency vs Window Size](https://www.sd-wan-experts.com/blog/why-your-maximum-thru-put-is-less-than-your-bandwidth/)

[Cloud Protocol Wars](https://www.theregister.co.uk/2015/10/01/aspera/)

[Free Software Plea](https://news.ycombinator.com/item?id=21898072)