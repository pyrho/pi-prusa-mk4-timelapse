# Timelapse Serial

2024-05: This is still a work in progress

## Goals

- Take a picture at each layer change with the extruder and bed at the same
  coordinates for each snap

## Actions

- `action:capture`
- `status:print_start`
- `status:print_stop`

## Building and Running

### Prerequisites
```shell
sudo apt install libudev-dev tmux libvips-dev libgphoto2
```

```shell
$> make build
$> sudo make install
$> tmux new-session -d -s timelapse-serial timelapse-serial -configPath /usr/local/etc/timelapse-serial.toml
```


## Config

A default config can be found in the following file: `configs/config.toml` 

### Find camera serial number
`$> lsusb -v`
look for your camera and the `iSerial` property

### Find the serial port for the printer

`$> ls -l /dev/serial/by-id/`

This should produce an output like:
```
total 0
lrwxrwxrwx 1 root root 13 May 22 18:05 usb-Prusa_Research__prusa3d.com__Original_Prusa_MK4_4914-27145608112151748-if00 -> ../../ttyACM0
```
(Thanks to [this kind stranger](https://stackoverflow.com/a/6768690/248978) for
the answer)

Which tells us that the printer's port is accessible via `/dev/ttyACM0`.
