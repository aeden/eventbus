$(document).ready(function() {
  // This is the App object. It is wrapped for jQuery
  var App = $({});

  // Run the EventBus application
  EventBus.run(App);

  // This is an app-specific event listener. It's job is to update the UI
  EventBus.listen(App, 'check.domain', function(evt, domainName) {
    $('#new-domain').hide();

    $('#checking-domain .domain-name').html(domainName);
    $('#checking-domain').show();
  });

  // Currently this can only handle one result
  EventBus.listen(App, 'check.domain.result', function(evt, result) {
    $('#checking-domain').hide();

    $('#domain-check-result .domain-name').html(result.name);
    $('#domain-check-result .domain-availability').html(result.availability);

    if (result.availability === 'available') {
      $('#domain-registration').show();

      $('#search-again-link').on("click", function(evt) {
        $('#new-domain').show();

        $('#domain-registration').hide();
        $('#domain-check-result').hide();
      });
    }

    $('#domain-check-result').show();
  });


  // This is a normal DOM event.
  $('#new-domain').on("submit", function(evt) {
    evt.preventDefault();
    EventBus.send(App, 'check.domain', [$("#domain-name").val()]);
  }); 

});
