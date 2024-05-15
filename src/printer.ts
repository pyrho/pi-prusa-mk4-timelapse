import { SerialPort } from "serialport";
import { ReadlineParser } from "serialport";
import { error, debug, log } from "./logger";
import { match } from "ts-pattern";
import { takePicture } from "./camera";
import fs from "node:fs/promises";

type Messages = "echo:print_start" | "echo:capture" | "echo:print_stop";

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
      .with("echo:capture", async () => {
        log("Capturing...");
        await takePicture().then((d) =>
          fs.writeFile(`/tmp/${+new Date()}.jpg`, d),
        );
        log("Captured!");
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
