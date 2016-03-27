

var mobSockets = []
var Mobs = []

function loadMobs(mobs,mapID) {

  if(mobSockets) {
    mobSockets.forEach(function(e,i){
      e.close()
    });

    Mobs.forEach(function(e,i){
      $("mob-".e.memID).remove()
    });
  }

  var mobSockets = []

  var wsuri = "ws://127.0.0.1:9020/getMob";

  mobs.forEach(function(e,f){
    mobSockets[f] = new WebSocket(wsuri);
    var mobID = e

    mobSockets[f].onopen = function() {
      console.log("mob connected to " + wsuri);
      this.send("map-id:"+mapID)
      this.send("mob-id:"+mobID)
    }

    mobSockets[f].onclose = function(e) {
      console.log("mob connection closed (" + e.code + ")");
    }

    mobSockets[f].onerror = function (error) {
      console.log('mob Error: ' + error); //log errors
    };

    mobSockets[f].onmessage = function(e) {
      if(e.data.indexOf(":")) {
        var res = e.data, index = e.data.indexOf(":");
        cmd = res.substring(0,index)
        data = JSON.parse(res.substring(index+1,res.lenght))

        switch(cmd) {
          case "info":
            Mobs[mobID] = new Mob(mobID, data)
            break;
          case "move":
            Mobs[mobID].move(data)
            break;
        }

      }
    }
  });

}

var Mob = function(id, data) {
  this.memID = id
  this.id = data.id
  this.hp = data.hp
  this.positionX = data.positionX
  this.positionY = data.positionY
  this.walkSpeed = data.walkSpeed

  $("body").append("<div class='mob-"+this.id+"' id='mob-"+this.memID+"'></div>");

  var mob = $("#mob-"+this.memID)
  .css("transition","top "+this.walkSpeed+"s, left "+this.walkSpeed+"s")

  var cell = $("#map-grid").find("#cell-"+this.positionX+"-"+this.positionY+"");

  mob.css("top", cell.offset().top+"px")
  .css("left", cell.offset().left+"px");

  this.move = function(data) {
    var walkSpeed = this.walkSpeed;

    data.forEach(function(elem,i) {
      var cell = $("#map-grid").find("#cell-"+elem[0]+"-"+elem[1]+"");

      mob.css("top", cell.offset().top+"px")
      .css("left", cell.offset().left+"px");

      setTimeout(function(){
        //console.log("Poring-"+id+" moving to x:"+elem[0]+" y:"+elem[1])
        mob.css("top", cell.offset().top+"px")
        .css("left", cell.offset().left+"px");
      },i*(walkSpeed*800))
    });
  }
}
