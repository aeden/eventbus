# Overview

The purpose of this project is to provide a unified event bus that can be used to trigger actions in browser-based JavaScript applications as well as from applications in other languages that are running on remote servers. Events are added to the bus using an HTTP end point and then are distributed to attached listeners via WebSockets.

![Communications Flow](http://cl.ly/image/2i3Q2z0e2U3e/communications.png)

This project is largely about the design of applications using this architecture, not just this implementation. Alternative implementations could be created with other components, such as existing open source or commercial message queues and web servers.

Thus, some of the key considerations are:

* What does an end-to-end event-driven web application look like?
* Is it a good way to develop web applications?
* What are the benefits and drawbacks when compared to existing synchronous and synchronous/asynchronous solutions?

This system's design borrows heavily from:

* The event-driven design of JavaScript.
* Concepts for micro-services.
* Existing queue systems such as ZeroMQ and Kafka.

# Status

This project is brand new. It is NOT production-ready.

# Design

## Definitions

### Client

A client is a web browser, native client or application. Clients must always identify themselves after connecting via their WebSocket. Clients must also always include their identity in their event context when creating new events using the `identifier` key.

### Services

A service is application code that processes events and executes business logic. Services must always authenticate themselves after connecting via their WebSocket. Services must also always pass their authorization token whenever they create new events.

## Frontend

The sample front end is HTML and JavaScript, however native clients could also be used.

The system relies solely on events as the mechanism for dealing with input and output. For example, a form is filled in an the submit button is pressed. This results in an event. Events are sent to local listeners to update the UI as well as to a remote event bus for any additional processing. When a listener on the remote event bus handles the event it will fire a new event, directed towards the client's event bus (via WebSockets) where a listener handles the event and updates the UI. Browser UI events are translated from the UI event (i.e. button clicked) to an event that is domain-specific (for example "check-domain-availability").

## Backend

The backend is purely functional event handlers. A listener handles the event, perhaps firing off events to the bus, potentially handled by other listeners or the client for UI updates. A backend listener could be implemented in any language as long as it can make HTTP calls and has a WebSocket library.

The backend consists of the following:

* A static file server for serving up the JS/HTML/CSS files
* An HTTP endpoint for pushing events into the event bus
* A WebSocket endpoint for establishing WebSocket connections
* Event routing logic for determining who can receive an event
* A stateful identity store for tracking client message source
* An authentication layer for services

## The EventBus Queue

The event bus is a publish/subscribe queue. Events are published to the bus using a synchronous HTTP call. A successfully added event returns a 200. Any other response code indicates that the event add failed and should be retried.

A successfully added event is persisted in an event store. The current implementation is in-memory and disappears when the application terminates.

Subscribers always receive events via WebSockets. Clients only receive events that are addressed to them. Services receive all events.

There is currently no guarantee that an event will be delivered.

# Running

Install Foreman.

Set up a .env file as follows:

```
HTTP_FILE_SERVER_HOST=locahost
HTTP_FILE_SERVER_PORT=3000
HTTP_EVENTBUS_SERVER_HOST=localhost
HTTP_EVENTBUS_SERVER_PORT=3001
```

Then run `foreman start`.

## CORS

If you are running the service under a different server (such as Apache or nginx) then you will need to make sure that the correct CORS headers are present when HTML files are requested.

The static HTML server response headers must include the following CORS headers:

```
Access-Control-Allow-Origin: eventbus-url
Access-Control-Allow-Headers: Content-Type
```

Where eventbus-url is replaced with the URL for the EventBus server.

# Example

There is an example application in the `static` directory. To see it in action, do the following:

* Run the EventBus server as described above.
* In a separate console, run the example Ruby script in `examples/ruby` using the command `foreman start` from within that directory.
* Open your browser to [http://localhost:3000](http://localhost:3000)

Your web browser must support WebSockets.

The sample Ruby script uses Event Machine.

# FAQ

In which I explain certain design decisions.

### Why not use WebSockets for both publishing events in and distributing events out?

The HTTP protocol is ideally suited for the synchronous nature of publishing an event to the queue. The system can take advantage of the protocol's built-in mechanisms for authorization, errors, etc. On the other hand, WebSockets are great for streaming data once a connection is established. Server-sent events (SSE) are another option: http://www.html5rocks.com/en/tutorials/eventsource/basics/ - right now IE has no support for SSE but does have support for WebSockets.

# Issues

* Currently the EventBus WebSocket service handles only unecrypted calls.
* The event system does not guarantee at-least-once delivery to connected services.
* Events are persistent but subscribers have no way of getting at old events.
