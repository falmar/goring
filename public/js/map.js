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
          buildMapGrid(data.name,data.size)
          loadMobs(data.mobs,data.id)
          break;
      }

    }

    console.log("message received: " + e.data);
  }
}

function buildMapGrid(name,size) {
  var table = $("<span><table class='text-center'><thead><tr><th class='text-center' colspan='"+size[1]+"'>Map: "+name+"</th></tr></thead><tbody></tbody></table></span>");

  for (f=0; f<size[0]; f++) {
    var tr = $("<tr></tr>");
    for (z = 0; z<size[1]; z++) {
      var td = $("<td style='background-color: #FFF;'>&nbsp;</td>");
      td.attr("id","cell-"+f+"-"+z);
      tr.append(td);
    }
    table.find("tbody").append(tr);
  }

  $("#map-grid").html(table.html())
}
