import gphoto2, { Camera } from "gphoto2";
import { debug } from "./logger";

const GP2 = gphoto2.GPhoto2;

export const takePicture = (() => {
  let cameraInstance: Camera | null = null;

  const getCamera = (o?: { reset: boolean }) =>
    new Promise<Camera>((resolve, reject) => {
      if (Boolean(o?.reset) || cameraInstance === null) {
        new GP2().list((cameras) => {
          const [camera] = cameras;
          if (camera === undefined)
            return reject(new Error("No cameras found"));
          debug(`Found ${camera.model}, ready.`);
          cameraInstance = camera;
          return resolve(cameraInstance);
        });
      } else {
        return resolve(cameraInstance);
      }
    });

  // Get the camera at load time to speed up the first picture
  getCamera().then(() => debug("Camera initialized at load time"));

  const takePicture = (retried = false): Promise<Buffer> =>
    new Promise<Buffer>((resolve, reject) =>
      getCamera({ reset: retried }).then((camera) => {
        debug("Taking picture...");
        return camera.takePicture({ download: true }, (err, data) => {
          if (err !== null && err !== undefined) {
            // In the "normal" case, this will happen
            // if the camera is disconnected
            // Retry once, retrying will also reset the camera instance
            if (!retried) return takePicture(true);
            return reject(err);
          } else {
            return resolve(data);
          }
        });
      }),
    );
  return takePicture;
})();
