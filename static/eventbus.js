var EventBus = {
  logging: true,
  remote: null
};

// This runs the EventBus websocket client
EventBus.run = function(source) {
  if (window["WebSocket"]) {
    conn = new WebSocket("ws://localhost:3000/ws");
    conn.onclose = function(evt) {
      EventBus.log("Connection closed.");
    }
    conn.onmessage = function(evt) {
      EventBus.log("Received event from EventBus");
      eventData = JSON.parse(evt.data);
      EventBus.log(eventData.data);
      source.trigger(eventData.name, eventData.data);
    }
  } else {
    EventBus.log("WebSockets not supported");
  }
}

// This sends the invent locally and to the remote event bus.
EventBus.send = function(source, eventName, data) {
  source.trigger(eventName, data);
  EventBus.log("Sending to remote event bus: " + EventBus.remote);
  EventBus.log(eventName, data);
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
