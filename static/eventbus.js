var EventBus = {
  logging: true,
  websocket_url: null,
  eventbus_url: null,
  context: {}
};

// This runs the EventBus websocket client
EventBus.connect = function(source, retry_count) {
  if(typeof(retry_count) === 'undefined') { retry_count = 0 }
  EventBus.log("Retry count: " + retry_count);

  if (window["WebSocket"]) {
    conn = new WebSocket(EventBus.websocket_url);

    conn.onopen = function(evt) {
      EventBus.log("Connection opened");
      conn.send(JSON.stringify({action: "identify"}));
      retry_count = 0;
    }

    conn.onclose = function(evt) {
      EventBus.log("Connection closed");
      setTimeout(function() { EventBus.connect(source, retry_count + 1) }, (retry_count + 1) * 1000);
    }

    conn.onmessage = function(evt) {
      EventBus.log("Received event from web socket");
      EventBus.log(evt);
      eventData = JSON.parse(evt.data);
      if (eventData['action']) {
        EventBus.log("Web socket message was an action");
        EventBus.context['identifier'] = eventData.token;
      } else {
        EventBus.log("Web socket message was an event");
        source.trigger(eventData.name, eventData.data);
      }
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
    url: EventBus.eventbus_url,
    dataType: 'json',
    data: JSON.stringify({
      name: eventName,
      context: EventBus.context,
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
