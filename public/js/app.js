
$(document).ready(function(){

  $("#map_form").submit(function(e){
    e.preventDefault();
    var mapID = $(this).find("#map-id").val();
    if(mapID) {
      loadMap(mapID);
    }
  });

  loadMap('prontera');
});
