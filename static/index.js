$(document).ready(function() {
  // This is the App object. It is wrapped for jQuery
  var App = $({});

  App.components = [];
  App.components.push($('#new-domain'));
  App.components.push($('#checking-domain'));
  App.components.push($('#domain-check-result'));
  App.components.push($('#domain-registration'));
  App.components.push($('#registering-domain'));
  App.components.push($('#domain-registered'));

  App.show = function(element) {
    $.each(App.components, function(index, component) {
      component.addClass("hidden");
    });
    element.removeClass("hidden");
  }

  // Connect the application to the EventBus
  EventBus.connect(App);

  // This is an app-specific event listener. It's job is to update the UI
  EventBus.listen(App, 'check-domain', function(evt, domainName) {
    $('.domain-name').html(domainName);
    App.show($('#checking-domain'));
  });

  // Currently this can only handle one result
  EventBus.listen(App, 'check-domain-completed', function(evt, result) {
    $('.domain-name').html(result.name);
    $('.domain-availability').html(result.availability);

    if (result.availability === 'available') {
      $('#domain-name-field').val(result.name);
      App.show($('#domain-registration'));
    } else {
      App.show($('#domain-check-result'));
    }
  });

  EventBus.listen(App, 'register-domain', function(evt, domainRegistration) {
    App.show($('#registering-domain'));
  });

  // EventBus event that occurs when the domain registration is completed
  EventBus.listen(App, 'register-domain-completed', function(evt, result) {
    EventBus.log("Received registration completed");
    EventBus.log(result);
    if (result.registered) {
      App.show($('#domain-registered'));
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
  $('.new-domain-search').on("click", function(evt) {
    $('#domain-name').val("");
    App.show($('#new-domain'));
  });

});
