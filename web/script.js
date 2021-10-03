function connect(uuid) {
  let wsUrl;
  if (window.location.protocol == "https:") {
    wsUrl = "wss://"
  } else {
    wsUrl = "ws://"
  }
  wsUrl += `${window.location.host}/ws`

  let ws = new WebSocket(wsUrl);
  ws.addEventListener('open', function (event) {
    ws.send(uuid);
  });

  ws.addEventListener('close', function (event) {
    console.log("Reconnecting...");
    setTimeout(connect(uuid), 5000);
  });
  ws.addEventListener('message', function (event) {
    let data = event.data;
    let arr = data.split(",", 2);
    let srcUUID = arr[0];
    if (arr.length != 2) {
      console.log(`Invalid message: ${data}`);
      return;
    }

    if (window.confirm(data)) {
      ws.send(`${srcUUID},accept`);
    } else {
      ws.send(`${srcUUID},deny`);
    }
  });
  foo = ws;
}

var foo;

async function init() {
  let uuid = localStorage.getItem("uuid");
  if (uuid == null) {
    console.log("No UUID found");
    let response = await fetch('/uuid');
    if (response.ok) {
      let text = await response.text();
      if (text.length == 36) {
        uuid = text;
        localStorage.setItem("uuid", uuid);
        console.log(`Got UUID: ${uuid}`);
      } else {
        console.log(`Error invalid UUID: ${text}`);
      }
    } else {
      console.log(`Error getting UUID: ${response}`)
    }
  } else {
    console.log(`UUID already set: ${uuid}`)
  }
  if (uuid == null) {
    return;
  }

  document.getElementById("uuid").textContent = uuid;

  connect(uuid);
}

window.addEventListener("DOMContentLoaded", init);
