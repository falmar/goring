
$(document).ready(function(){

  $("#map_form").submit(function(e) {
    e.preventDefault();

    $.ajax({
      url: "/getMap",
      dataType: 'json'
      success: function(data) {
        alert(data.id)
      }
    })

  });

})
