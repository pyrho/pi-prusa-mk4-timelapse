window.addEventListener("focus", () => {
  if (document.querySelector("#live-feed")?.src == null) return;

  console.log(document.querySelector("#live-feed")?.src == null);

  const backup = document.querySelector("#live-feed")?.src;
  document.querySelector("#live-feed").src = "";
  document.querySelector("#live-feed").src = backup;
  console.log("refreshed");
});
