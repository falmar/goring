var players = [];
var playersIDs = [];
var playersSockets = [];

function loadPlayers(p,mapID){

  if(playersSockets) {
    playersSockets.forEach(function(e,i){
      e.close();
    });

    playersIDs.forEach(function(e,i){
      $("#player-"+e).remove();
    });
  }

  playersSockets = [];
  playersIDs = [];
  players = [];

  var wsuri = "ws://127.0.0.1:9020/getPlayer";

  p.forEach(function(e,i){
    playersSockets[i] = new WebSocket(wsuri);
    var playerID = e;

    playersSockets[i].onopen = function() {
      console.log("player connected to " + wsuri);
      this.send("map-id:"+mapID);
      this.send("player-id:"+playerID);
    };

    playersSockets[i].onclose = function(e) {
      console.log("player connection closed (" + e.code + ")");
    };

    playersSockets[i].onerror = function (error) {
      console.log('player Error: ' + error); //log errors
    };

    playersSockets[i].onmessage = function(e) {
      if(e.data.indexOf(":")) {
        var res = e.data, index = e.data.indexOf(":");
        cmd = res.substring(0,index);
        data = JSON.parse(res.substring(index+1,res.lenght));

        switch(cmd) {
          case "info":
            players[playerID] = new Player(playerID, data);
            playersIDs.push(playerID);
            break;
          case "p_move":
            players[playerID].move(data);
            break;
          case "p_dmg":
            players[playerID].dmg(data);
            break;
          case "p_die":
            players[playerID].die();
            break;
          case "p_respawn":
            players[playerID].respawn(data);
            break;
        }

      }
    };

  });
}

var playerStatusIdle = 1;
var playerStatusCombat = 2;
var playerStatusDead = 3;

var Player = function(id,data){
  this.memID = id;
  this.id = data.id;
  this.hp = data.hp;
  this.maxHP = data.maxHP;
  this.positionX = data.positionX;
  this.positionY = data.positionY;
  this.dead = data.dead;
  this.status = playerStatusIdle;

  // start



  $("body").append("<div class='hide player' id='player-"+this.memID+"'>"+
  '<meter value="'+(this.hp * 100) /this.maxHP+'" min="0" low="25" high="60" optimum="100" max="100"></meter>'+
  "</div>");

  var player = $("#player-"+this.memID)
  .css("transition","top "+this.walkSpeed+"ms, left "+this.walkSpeed+"ms");

  this.move = function(data) {
    var cell = $("#map-grid").find("#cell-"+data[1]+"-"+data[0]+"");
    player.css("top", cell.offset().top-50+"px")
    .css("left", cell.offset().left+"px");
  };

  this.dmg = function(dmg){
    this.hp += dmg * -1;
    player.find("meter").attr("value",(this.hp * 100) / this.maxHP);
    var self = this;

    if(this.status == playerStatusIdle) {
      this.status = playerStatusCombat;
      player.addClass('dmg');
      player.css("top", player.offset().top+12+"px");
      setTimeout(function(){
        if (self.status == playerStatusCombat) {
          self.status = playerStatusIdle;
          player.css("top", player.offset().top-12+"px");
        }
        player.removeClass('dmg');
      },600);
    }
  };

  this.die = function(){
    this.hp = 0;
    $(player).find("meter").attr("value",0);
    this.status = playerStatusDead;
    player.css("top", player.offset().top+31+"px");
    player.addClass('dead');
  };

  this.respawn = function(data){
    this.move([data.positionX,data.positionY]);
    this.hp = data.hp;
    $(player).find("meter").attr("value",(this.hp * 100) / this.maxHP);
    this.status = playerStatusIdle;
    player.removeClass('dead');
  };

  if(!this.dead) {
    player.removeClass('hide');
    this.move([this.positionX,this.positionY]);
  }

};
