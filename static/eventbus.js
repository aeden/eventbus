var EventBus = {
  logging: true,
  websocket_url: null,
  eventbus_url: null
};

// This runs the EventBus websocket client
EventBus.run = function(source) {
  if (window["WebSocket"]) {
    conn = new WebSocket(websocket_url);
    conn.onclose = function(evt) {
      EventBus.log("Connection closed.");
    }
    conn.onmessage = function(evt) {
      EventBus.log("Received event from web socket");
      EventBus.log(evt);
      eventData = JSON.parse(evt.data);
      source.trigger(eventData.name, eventData.data);
    }
  } else {
    EventBus.log("WebSockets not supported");
  }
}

// This sends the event locally and to the remote event bus.
EventBus.send = function(source, eventName, data) {
  source.trigger(eventName, data);
  EventBus.log("Sending " + eventName + " event to remote event bus: " + EventBus.eventbus_url);
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

EventBus.log = function(objects) {
  if (EventBus.logging) {
    console.log(objects);
  }
}
