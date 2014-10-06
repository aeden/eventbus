var EventBus = {};

EventBus.url = function(url) {
  EventBus.remote = url
}

// This sends the invent locally and to the remote event bus.
EventBus.send = function(source, eventName, data) {
  source.trigger(eventName, data);
  console.log("Sending to remote event bus: " + EventBus.remote);
  console.log(eventName, data);
  jQuery.ajax({
    method: 'post',
    url: EventBus.remote,
    dataType: 'json',
    data: JSON.stringify({
      name: eventName,
      data: data
    })
  });
}

// This registers a listener.
EventBus.listen = function(source, events, data, handler) {
  source.on(events, null, data, handler); 
}
