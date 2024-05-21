# Timelapse Serial

2024-05: This is still a work in progress

## Goals

- Take a picture at each layer change with the extruder and bed at the same
  coordinates for each snap

## Actions

- `action:capture`
- `status:print_start`
- `status:print_stop`

## Find camera serial number
`$> lsusb -v`
look for your camera and the `iSerial` property
