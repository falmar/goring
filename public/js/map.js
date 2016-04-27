var mapSocket;

function loadMap(mapID) {

  if(mapSocket) {
    mapSocket.close();
  }

  var wsuri = "ws://127.0.0.1:9020/getMap";

  mapSocket = new WebSocket(wsuri);

  mapSocket.onopen = function() {
    console.log("connected to " + wsuri);
    mapSocket.send("map-id:"+mapID);
  };

  mapSocket.onclose = function(e) {
    console.log("connection closed (" + e.code + ")");
  };

  mapSocket.onerror = function (error) {
    console.log('Error Logged: ' + error); //log errors
  };

  mapSocket.onmessage = function(e) {

    if(e.data.indexOf(":")) {
      var res = e.data, index = e.data.indexOf(":");
      cmd = res.substring(0,index);
      data = JSON.parse(res.substring(index+1,res.lenght));

      switch(cmd) {
        case "info":
          buildMapGrid(data.size);
          loadMobs(data.mobs,data.id);
          loadPlayers(data.players,data.id);
          break;
      }

    }

    console.log("message received: " + e.data);
  };
}

function buildMapGrid2(name,size) {
  var table = $("<span><table class='text-center'><thead><tr><th class='text-center' colspan='"+size[0]+"'>Map: "+name+"</th></tr></thead><tbody></tbody></table></span>");

  for (f=1; f<=size[1]; f++) {
    var tr = $("<tr></tr>");
    for (z = 1; z<=size[0]; z++) {
      var td = $("<td>&nbsp;</td>");
      td.attr("id","cell-"+f+"-"+z);
      tr.append(td);
    }
    table.find("tbody").append(tr);
  }

  $("#map-grid").html(table.html());
}

function buildMapGrid(size) {
  var grid = $("<span></span>");

  for (f=1; f<=size[1]; f++) {
    var tr = $("<div class='tr shrink'></div>");
    for (z = 1; z<=size[0]; z++) {
      var td = $("<div class='td'>&nbsp;</div>");
      td.attr("id","cell-"+f+"-"+z);
      tr.append(td);
    }
    grid.append(tr);
  }

  $("#map-grid").css('width',(size[0] * 40)+5 +"px");
  $("#map-grid").html(grid.html());
}
