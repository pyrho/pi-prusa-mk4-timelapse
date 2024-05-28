function resetStream() {
  if (document.querySelector("#live-feed")?.src == null) return;

  console.log(document.querySelector("#live-feed")?.src == null);

  const backup = document.querySelector("#live-feed")?.src;
  document.querySelector("#live-feed").src = "";
  document.querySelector("#live-feed").src = backup;
  console.log("refreshed");

}
window.addEventListener("focus", resetStream);
window.addEventListener("touchend", resetStream);
