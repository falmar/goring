function loadPlayers(players){

  $("body").append("<div class='player' id='player-"+players[0]+"'></div>");

  var cell = $("#map-grid").find("#cell-"+players[1]+"-"+players[2]+"");

  player = $("#player-"+players[0]);

  player.css("top", cell.offset().top+"px")
  .css("left", cell.offset().left+"px");

}
