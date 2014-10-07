# Overview

The puropose of this project is to provide an event bus that can be used from within web browsers via JavaScript as well as from applications in other languages. Events are added to the bus using an HTTP end point and then are distributed to attached listeners via WebSockets.

# Status

This project is brand new. It is NOT production-ready.

# Design

## Frontend

The front end is HTML and JavaScript. It relies solely on events as the mechanism for dealing with input and output. For example, a form is filled in an the submit button is pressed. This results in an event. Events are sent to local listeners to update the UI as well as to a remote event bus for any additional processing. When a listener on the remote event bus handles the event it will fire a new event, directed towards the client's event bus (via websockets) where a listener handles the event and updates the UI. Browser UI events are translated from the UI event (i.e. button clicked) to an event that is domain-specific (for example "check-domain-availability").

## Backend

The backend is purely functional event handlers. A listener handles the event, perhaps firing off events to the bus, potentially handled by other listeners or the client for UI updates. A backend listener could be implemented in any language as long as it can make HTTP calls and has a WebSocket library.

## Creating Events

Events are always added to the bus using a synchronous HTTP call. A successfully added event returns a 200. Any other response code indicates that the event add failed and should be retried.

## Authentication

### Clients

Clients represent browsers or devices.

A client must identify itself when the websocket connection is established. This must occur before it sends its first event.

Events must include the client identifier in the event context using the key `identifier`.

### Services

Services are listeners that can handle events and execute work.

A service must include an authorization token on each request.

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

# Issues

* Currently the EventBus WebSocket service handles only unecrypted calls.
