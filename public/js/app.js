
$(document).ready(function(){

  var sock = null;
  var wsuri = "ws://127.0.0.1:9020/getMap";

  sock = new WebSocket(wsuri);

  sock.onopen = function() {
    console.log("connected to " + wsuri);
    sock.send("map-id:prontera")
  }

  sock.onclose = function(e) {
    console.log("connection closed (" + e.code + ")");
  }

  sock.onerror = function (error) {
    console.log('Error Logged: ' + error); //log errors
  };

  sock.onmessage = function(e) {
    console.log("message received: " + e.data);
  }

})
