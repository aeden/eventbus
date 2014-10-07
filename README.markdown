# Overview

## Frontend

The front end is HTML and JavaScript. It relies solely on events as the mechanism for dealing with input and output. For example, a form is filled in an the submit button is pressed. This results in an event. Events are sent to local listeners to update the UI as well as to a remote event bus for any additional processing. When a listener on the remote event bus handles the event it will fire a new event, directed towards the client's event bus (via websockets) where a listener handles the event and updates the UI. 

## Backend

The backend is purely functional event handlers. A listener handles the event, perhaps firing off events to the bus, potentially handled by other listeners or the client for UI updates.

## Creating Events

Events are always added to the bus using a synchronous HTTP call. A successfully added event returns a 200. Any other response code indicates that the event add failed and should be retried.

## Authentication

### Clients

Clients represent browsers or devices.

### Services

Services are listeners that can handle events and execute work.

### Services

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

# Issues

* Currently the EventBus WebSocket service handles only unecrypted calls.
* There is no authentication or authorization context.
* All messages are routed to all clients.
