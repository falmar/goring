var mapSocket

function loadMap(mapID) {

  if(mapSocket) {
    mapSocket.close()
  }

  var wsuri = "ws://127.0.0.1:9020/getMap";

  mapSocket = new WebSocket(wsuri);

  mapSocket.onopen = function() {
    console.log("connected to " + wsuri);
    mapSocket.send("map-id:"+mapID)
  }

  mapSocket.onclose = function(e) {
    console.log("connection closed (" + e.code + ")");
  }

  mapSocket.onerror = function (error) {
    console.log('Error Logged: ' + error); //log errors
  };

  mapSocket.onmessage = function(e) {

    if(e.data.indexOf(":")) {
      var res = e.data, index = e.data.indexOf(":");
      cmd = res.substring(0,index)
      data = JSON.parse(res.substring(index+1,res.lenght))

      switch(cmd) {
        case "info":
          $("#map-name").html(data.name)
          buildMapGrid(data.size)
          break;
      }

    }

    console.log("message received: " + e.data);
  }
}

function buildMapGrid(size) {
  
}
