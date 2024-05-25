# TODO

## Timelapsin'
- [ ] Use fluent FFMPEG instead of a system call
- [ ] Display number of thumbnails for folder
- [x] Figure out thumbnail creation OOM kill
- [x] Paginate folder list section
- [x] Web server to show status and access to timelapses
    - use htmx and go's net/http package
- [x] Spool up camera
- [x] Retry when the printer disconnects
- [x] Listen when the camera connects (to reconnect to it midway)

## Raspberry pi camera setup
This should be running in another `cmd/` as it's wholly unrelated to timelapses.

- [ ] Implement a simple webpage to display the camera
- [ ] Implement the snapshot thing to upload to prusaconnect
- [ ] Move the camera code to the go server
  https://medium.com/go4vl/build-a-wifi-camera-using-the-raspberry-pi-zero-w-a-camera-module-and-go-1d5fadfa7d76
  Or 
  https://gocv.io/writing-code/more-examples/

## PrusaLink
- [ ] Fetch info from PrusaLink

### Links
This is to get the PNG thumbnail
```
shell
curl 'http://mk4.lan/thumb/l/usb/PART~23F.BGC' \
  -H 'Accept: image/*' \
  -H 'Accept-Language: en-US,en;q=0.9' \
  -H 'Authorization: Digest username="maker", realm="Printer API", nonce="b544617a0001754a", uri="/thumb/l/usb/PART~23F.BGC", response="ae238d92cea2a0684a8f4cd28b6fe49c"' \
  -H 'Connection: keep-alive' \
  -H 'DNT: 1' \
  -H 'If-None-Match: "713555695"' \
  -H 'Referer: http://mk4.lan/' \
  -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36' \
  --insecure
```

```shell
$> http -vvv mk4.lan/api/v1/status X-Api-Key:YOUR_PRUSALINK_PW Accept:application/json
```



## Other
- [ ] Create a proper systemd unit instead of relying on tmux
