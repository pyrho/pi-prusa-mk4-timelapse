const gphoto2 = require("gphoto2");
const GPhoto = new gphoto2.GPhoto2();
const { debugLog } = require("./debug");

module.exports.takePicture = (() => {
  let cameraInstance = null;

  const getCamera = () =>
    new Promise((resolve, reject) => {
      if (cameraInstance === null) {
        new GPhoto().list((cameras) => {
          const [camera] = cameras;
          if (camera === undefined)
            return reject(new Error("No cameras found"));
          debugLog(`Found ${camera.model}, ready.`);
          cameraInstance = camera;
          return resolve(cameraInstance);
        });
      } else {
        return resolve(cameraInstance);
      }
    });

  // Get the camera at load time to speed up the first picture
  getCamera().then(() => debugLog("Camera initialized at load time"));

  const takePicture = (retried = false) =>
    new Promise((resolve, reject) => {
      getCamera().then((camera) => {
        debugLog("Taking picture...");
        camera.takePicture({ download: true }, function (err, data) {
          if (err) {
            // In the "normal" case, this will happen
            // if the camera is disconnected
            // Setting `cameraInstance` to null will force a reinstantiation
            // of the camera
            cameraInstance = null;
            // Retry once
            if (!retried) return takePicture(true);
            return reject(err);
          } else {
            return resolve(data);
          }
        });
      });
    });
  return takePicture;
})();
