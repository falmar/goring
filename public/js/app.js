var WSHost = ''

$(document).ready(function () {
  WSHost = $('#ws_host').val()

  // $('#map_form').submit(function (e) {
  //   e.preventDefault()
  //
  //   var mapID = $(this).find('#map-id').val()
  //
  //   if (mapID) {
  //     loadMap(mapID)
  //   }
  // })

  loadMap('prontera')
})
