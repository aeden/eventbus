$(document).ready(function() {
  // This is the App object. It is wrapped for jQuery
  var App = $({});

  // Run the EventBus application
  EventBus.run(App);

  // This is an app-specific event listener. It's job is to update the UI
  EventBus.listen(App, 'check.domain', function(evt, domainName) {
    $('#new-domain').addClass("hidden");

    $('#checking-domain .domain-name').html(domainName);
    $('#checking-domain').removeClass("hidden");
  });

  // Currently this can only handle one result
  EventBus.listen(App, 'check.domain.result', function(evt, result) {
    $('#checking-domain').addClass("hidden");

    $('#domain-check-result .domain-name').html(result.name);
    $('#domain-check-result .domain-availability').html(result.availability);

    if (result.availability === 'available') {
      $(".check-result-icon").addClass('glyphicon').addClass('glyphicon-ok');
      $(".check-result").addClass('alert-success');
      $('#domain-registration').removeClass("hidden");
    }

    $('#domain-check-result').removeClass("hidden");
  });


  // This is a normal DOM event.
  $('#new-domain').on("submit", function(evt) {
    evt.preventDefault();
    EventBus.send(App, 'check.domain', [$("#domain-name").val()]);
  });

  $('#register-domain').on("submit", function(evt) {
    evt.preventDefault();
    EventBus.send(App, 'register.domain', $('#register-domain').serializeJSON());
  });

  $('#search-again-link').on("click", function(evt) {
    $('#domain-name').val("");
    $('#new-domain').removeClass("hidden");

    $('#domain-registration').addClass("hidden");
    $('#domain-check-result').addClass("hidden");
  });

});
