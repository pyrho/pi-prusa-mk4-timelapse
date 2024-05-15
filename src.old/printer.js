const { SerialPort } = require("serialport");
const { debugLog } = require("./debug");
const { ReadlineParser } = require("@serialport/parser-readline");

function startPrinterListen() {
  // Configure the serial port
  const port = new SerialPort({
    path: "/dev/ttyACM0",
    baudRate: 115200,
  });

  // Set up a readline parser
  const parser = port.pipe(new ReadlineParser({ delimiter: "\n" }));

  // Open the serial port
  port.on("open", () => {
    debugLog("Serial port opened");
  });

  // Read data from the serial port
  parser.on("data", (data) => {
    debugLog("Received:", data);
    messageHandler(data);
  });

  // Error handling
  port.on("error", (err) => {
    debugLog("Error: ", err.message);
  });
}

function messageHandler(message) {
  if (message === "action:new_print_start") {
      // Create a new folder to store the images
    return;
  } else if (message === "action:capture_img") {
      // Store the pics in the current print's folder
      // (if it exists, deal with crashes and restarts?)
  } else {
    return;
  }
}

module.exports = { startPrinterListen };
