$(document).ready(function() {
  // This happens once the document ready is executed
  console.log("Running app");

  // This is the App object. It is wrapped for jQuery
  var App = $({});

  // This is an app-specific event listener. It's job is to update the UI
  EventBus.listen(App, 'check.domain', function(evt, domainName) {
    console.log(domainName);
    $('#new-domain').hide();
    
    $('#checking-domain .domain-name').html(domainName);
    $('#checking-domain').show();
  });


  // This is a normal DOM event.
  $('#new-domain').on("submit", function(evt) {
    evt.preventDefault();
    EventBus.send(App, 'check.domain', [$("#domain-name").val()]);
  });
});
