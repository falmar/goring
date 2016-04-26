

var mobSockets = [];
var Mobs = [];
var MobIDs = [];

function loadMobs(mobs,mapID) {

  if(mobSockets) {
    mobSockets.forEach(function(e,i){
      e.close();
    });

    MobIDs.forEach(function(e,i){
      $("#mob-"+e).remove();
    });
  }

  mobSockets = [];
  Mobs = [];
  MobIDs = [];

  var wsuri = "ws://127.0.0.1:9020/getMob";

  mobs.forEach(function(e,f){
    mobSockets[f] = new WebSocket(wsuri);
    var mobID = e;

    mobSockets[f].onopen = function() {
      console.log("mob connected to " + wsuri);
      this.send("map-id:"+mapID);
      this.send("mob-id:"+mobID);
    };

    mobSockets[f].onclose = function(e) {
      console.log("mob connection closed (" + e.code + ")");
    };

    mobSockets[f].onerror = function (error) {
      console.log('mob Error: ' + error); //log errors
    };

    mobSockets[f].onmessage = function(e) {
      if(e.data.indexOf(":")) {
        var res = e.data, index = e.data.indexOf(":");
        cmd = res.substring(0,index);
        data = JSON.parse(res.substring(index+1,res.lenght));

        switch(cmd) {
          case "info":
            Mobs[mobID] = new Mob(mobID, data);
            MobIDs.push(mobID);
            break;
          case "move":
            Mobs[mobID].move(data);
            break;
          case "die":
            Mobs[mobID].die();
            break;
          case "respawn":
            Mobs[mobID].respawn(data);
            break;
        }

      }
    };
  });

}

var Mob = function(id, data) {
  this.memID = id;
  this.id = data.id;
  this.hp = data.hp;
  this.positionX = data.positionX;
  this.positionY = data.positionY;
  this.walkSpeed = data.walkSpeed;

  this.move = function(data) {
    cell = $("#map-grid").find("#cell-"+data[1]+"-"+data[0]+"");
    mob.css("top", cell.offset().top+"px")
    .css("left", cell.offset().left+"px");
  };

  this.die = function(){
    mob.addClass('hide');
  };

  this.respawn = function(data){
    this.move([data.positionX,data.positionY]);
    this.hp = data.hp;
    mob.removeClass('hide');
  };

  // start

  $("body").append("<div class='hide mob-"+this.id+"' id='mob-"+this.memID+"'></div>");

  var mob = $("#mob-"+this.memID)
  .css("transition","top "+this.walkSpeed+"ms, left "+this.walkSpeed+"ms");

  if (!data.dead) {
    this.move([this.positionX,this.positionY]);
    mob.removeClass('hide');
  }

};
