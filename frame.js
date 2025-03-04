function FrameBar() {
  // Add the required CSS rules
  const style = document.createElement("style");
  const css = [
    "#fcframebar { padding: 0; margin: 0; border-bottom: 2px solid #8A63D2; width: 100%; }",
    "#fcframebar .container { width:360px; margin: 0 auto; padding: 0; text-align: center; font-size: 1em; display: flex; flex-wrap: wrap; }",
    "#fcframebar .container .info { width: 250px; padding: 10px; box-sizing: border-box; vertical-align: middle; font-size: 0.8em; }",
    "#fcframebar button { display: inline-block; padding: 10px 20px; background-color: #855DCD; color: white; border-radius: 5px; border: none; font-size: 0.8em; height: 32px; margin: auto 5px; }",
    "#fcframebar button:hover { background-color: #7A59C9; }",
  ].join(" ");
  style.innerHTML = css;
  document.head.appendChild(style);

  // Create the header
  const header = document.createElement("header");
  const wplink =
    "https://warpcast.com/~/frames/launch?url=" +
    encodeURIComponent(window.location);
  header.innerHTML = `<div id="fcframebar"><div class="container">
    <div class="info">
    <svg width="32" height="32" viewBox="0 0 1000 1000" fill="none" xmlns="http://www.w3.org/2000/svg" style="vertical-align: middle; float:left">
              <rect width="1000" height="1000" rx="200" fill="#855DCD"/>
              <path d="M257.778 155.556H742.222V844.444H671.111V528.889H670.414C662.554 441.677 589.258 373.333 500 373.333C410.742 373.333 337.446 441.677 329.586 528.889H328.889V844.444H257.778V155.556Z" fill="white"/>
              <path d="M128.889 253.333L157.778 351.111H182.222V746.667C169.949 746.667 160 756.616 160 768.889V795.556H155.556C143.283 795.556 133.333 805.505 133.333 817.778V844.444H382.222V817.778C382.222 805.505 372.273 795.556 360 795.556H355.556V768.889C355.556 756.616 345.606 746.667 333.333 746.667H306.667V253.333H128.889Z" fill="white"/>
              <path d="M675.556 746.667C663.283 746.667 653.333 756.616 653.333 768.889V795.556H648.889C636.616 795.556 626.667 805.505 626.667 817.778V844.444H875.556V817.778C875.556 805.505 865.606 795.556 853.333 795.556H848.889V768.889C848.889 756.616 838.94 746.667 826.667 746.667V351.111H851.111L880 253.333H702.222V746.667H675.556Z" fill="white"/>
              </svg>
      <div style="vertical-align: top; text-align:left; padding-left:36px;"><b>Farcaster Frame</b><br>Best experienced in Farcaster</div>
    </div>
  <button onclick="window.location.href='${wplink}';">Open</button></div></div>`;
  document.body.insertBefore(header, document.body.firstChild);
}

async function inFrameContext() {
  const links = document.querySelectorAll("a[href^='http']");
  for (let link of links) {
    link.addEventListener("click", (event) => {
      event.preventDefault();
      frame.sdk.actions.openUrl(link.getAttribute("href"));
    });
  }
  await frame.sdk.actions.addFrame();
}

window.onload = async () => {
  try {
    await frame.sdk.actions.ready();
    const ctx = await frame.sdk.context;
    ctx ? await inFrameContext() : FrameBar();
  } catch (error) {
    console.error(error);
  }
};
