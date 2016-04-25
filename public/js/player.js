var playersIDs = [];

function loadPlayers(players){

  $("#player-"+players[0]).remove();

  $("body").append("<div class='player' id='player-"+players[0]+"'></div>");

  var cell = $("#map-grid").find("#cell-"+players[1]+"-"+players[2]+"");

  player = $("#player-"+players[0]);

  player.css("top", cell.offset().top-50+"px")
  .css("left", cell.offset().left+"px");

}
