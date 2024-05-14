import { SerialPort } from "serialport";
import { ReadlineParser } from "@serialport/parser-readline";

export function startSerialRead() {
  // Configure the serial port
  const port = new SerialPort({
    path: "/dev/ttyACM0",
    baudRate: 115200,
  });

  // Set up a readline parser
  const parser = port.pipe(new ReadlineParser({ delimiter: "\n" }));

  // Open the serial port
  port.on("open", () => {
    console.log("Serial port opened");
  });

  // Read data from the serial port
  parser.on("data", (data) => {
    console.log("Received:", data);
  });

  // Error handling
  port.on("error", (err) => {
    console.error("Error: ", err.message);
  });
}
