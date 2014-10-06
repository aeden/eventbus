$(document).ready(function() {
  // This is the App object. It is wrapped for jQuery
  var App = $({});

  // Run the EventBus application
  EventBus.run(App);

  // This is an app-specific event listener. It's job is to update the UI
  EventBus.listen(App, 'check-domain', function(evt, domainName) {
    $('#new-domain').addClass("hidden");

    $('.domain-name').html(domainName);
    $('#checking-domain').removeClass("hidden");
  });

  // Currently this can only handle one result
  EventBus.listen(App, 'check-domain-completed', function(evt, result) {
    $('#checking-domain').addClass("hidden");

    $('.domain-name').html(result.name);
    $('#domain-check-result .domain-availability').html(result.availability);

    if (result.availability === 'available') {
      $(".check-result-icon").addClass('glyphicon').addClass('glyphicon-ok');
      $(".check-result").addClass('alert-success');
      $('#domain-name-field').val(result.name);
      $('#domain-registration').removeClass("hidden");
    }

    $('#domain-check-result').removeClass("hidden");
  });

  EventBus.listen(App, 'register-domain', function(evt, domainRegistration) {
    $('#domain-check-result').addClass('hidden');
    $('#domain-registration').addClass('hidden');
    $('#registering-domain').removeClass('hidden');
  });

  EventBus.listen(App, 'register-domain-completed', function(evt, result) {
    EventBus.log("Received registration completed");
    EventBus.log(result);
    if (result.registered) {
      $('#domain-registered').removeClass('hidden');
      $('#registering-domain').addClass('hidden');
    }
  });


  // DOM event that occurs when the check form is submitted.
  $('#new-domain').on("submit", function(evt) {
    evt.preventDefault();
    EventBus.send(App, 'check-domain', [$("#domain-name").val()]);
  });

  // DOM event that occurs when the registration is submitted.
  $('#register-domain').on("submit", function(evt) {
    evt.preventDefault();
    EventBus.send(App, 'register-domain', $('#register-domain').serializeObject());
  });

  // Search again
  $('#search-again-link').on("click", function(evt) {
    $('#domain-name').val("");
    $('#new-domain').removeClass("hidden");

    $('#domain-registration').addClass("hidden");
    $('#domain-check-result').addClass("hidden");
  });

});
