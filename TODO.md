# TODO

## Timelapsin'
- [ ] Use fluent FFMPEG instead of a system call
- [ ] Figure out thumbnail creation OOM kill
- [ ] Paginate folder list section
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

## Other
- [ ] Create a proper systemd unit instead of relying on tmux
