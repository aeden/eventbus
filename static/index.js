$(document).ready(function() {
  // This is the App object. It is wrapped for jQuery
  var App = $({});

  // These are components on the page
  App.components = [
    $('#domain-new'),
    $('#checking-domain'),
    $('#domain-check-result'),
    $('#domain-registration'),
    $('#registering-domain'),
    $('#domain-registered')
  ];

  // show one component, hide the others
  App.show = function(element) {
    $.each(App.components, function(index, component) {
      component.addClass("hidden");
    });
    element.removeClass("hidden");
  }

  /*------------------------------*
   * Start the EventBus setup     *
   *------------------------------*/

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
    evt.preventDefault();
    $('#domain-name').val("");
    App.show($('#domain-new'));
  });

});
