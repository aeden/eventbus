var EventBus = {};

EventBus.url = function(url) {
  EventBus.remote = url
}

EventBus.run = function() {
  if (window["WebSocket"]) {
    conn = new WebSocket("ws://localhost:3000/ws");
    conn.onclose = function(evt) {
      console.log("Connection closed.");
    }
    conn.onmessage = function(evt) {
      console.log("Received event from EventBus");
      console.log(evt.data);
    }
  } else {
    console.log("Your browser does not support WebSockets.");
  }
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
