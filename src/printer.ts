import { SerialPort } from "serialport";
const { execFileSync } = require('node:child_process');
import { ReadlineParser } from "serialport";
import { error, debug, log } from "./logger";
import { match } from "ts-pattern";
import { takePicture } from "./camera";
import fs from "node:fs/promises";
import process from 'node:process'
export const startPrinterSerialChannel = (): void => {
  const port = new SerialPort({
    path: "/dev/ttyACM0",
    baudRate: 115200,
  });

  // Set up a readline parser
  const parser = port.pipe(new ReadlineParser({ delimiter: "\n" }));

  // Open the serial port
  port.on("open", () => {
    log("Serial port opened");
  });

  // Read data from the serial port
  parser.on("data", (data) => {
    debug("Received:", data);

    match(data as unknown)
      .with("action:capture", async () => {
        log("Capturing...");
        const x = execFileSync('pidof',[ 'gphoto2'])
        process.kill(x, 'SIGUSR1')
      })
      .with("echo:print_stop", () => {})
      .with("echo:print_start", () => {})
      .otherwise(() => error(`Unknown message: ${data}`));
  });

  // Error handling
  port.on("error", (err) => {
    error("Error: ", err.message);
  });
};
