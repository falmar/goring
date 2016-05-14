
var ws_ip = "";

$(document).ready(function(){

  ws_ip = $("#ws_ip").val();

  $("#map_form").submit(function(e){
    e.preventDefault();
    var mapID = $(this).find("#map-id").val();
    if(mapID) {
      loadMap(mapID);
    }
  });

  loadMap('prontera');
});
