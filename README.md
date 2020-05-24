# MJPEG http handler

```sh
go get github.com/sg3des/mjpeg
```

## USAGE

```go
// initialize
stream := mjpeg.NewStream()

// serve
http.HandleFunc("/...", stream.ServeHTTP)

// set and update JPEG frames
for {
    ...
    stream.UpdateFrame(imageBytes)
}

```
