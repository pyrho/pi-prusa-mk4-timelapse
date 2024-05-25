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

### G-Code

#### Start G-Code
```gcode
;TIMELAPSE RELATED
M118 A1 status:print_start
```

#### End G-Code
```gcode
; TIMELAPSE START

; get next to the switch
G0 X{first_layer_print_min[0] + ((first_layer_print_max[0]-first_layer_print_min[0])/2)} Y{first_layer_print_max[1] + 35} F18000

;Wait for all moves to finish
M400

;Ask the Pi to take the pic
M118 A1 action:capture

; ... the rest

; TIMELAPSE
M118 A1 status:print_stop
```

#### After layer change 
Print head hover behind print
```gcode

; TIMELAPSE RELATED START
{ if layer_num >= 1 }
;Retract -0.8mm at 2100mm/min or 35mm/sec
;G1 E1 F2100 

; get next to the switch (hover)
G0 X{first_layer_print_min[0] + ((first_layer_print_max[0]-first_layer_print_min[0])/2)} Y{first_layer_print_max[1] + 35} F18000

;Wait for all moves to finish
M400

;Ask the Pi to take the pic
M118 A1 action:capture

;Wait for 300ms, to let the camera take the picture
G4 P300

;Use relative positioning mode
;G91
 ;Return Z to current layer
;G1 Z1 F2000
; Go back to absolute position mode for all axes
;G90

;Prime and lift
;G1 E1 F2100
{endif}
; TIMELAPSE RELATED END
```

Print head freeze
```gcode
just remove the (hover) line
```

Print next to print freeze (the +35 is to account for the fan duct, mine is
bigger than stock)
```gcode
G0 X{first_layer_print_min[0]} Y{first_layer_print_max[1] + 35} F18000
```

Print head parked
```gcode
G0 X0 Y170 F18000
```
