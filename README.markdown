# Overview

## Frontend

The front end is HTML and JavaScript. It relies solely on events as the mechanism for dealing with input and output. For example, a form is filled in an the submit button is pressed. This results in an event. Events are sent to local listeners to update the UI as well as to a remote event bus for any additional processing. When a listener on the remote event bus handles the event it will fire a new event, directed towards the client's event bus (via websockets) where a listener handles the event and updates the UI. 

## Backend

The backend is purely functional event handlers. A listener handles the event, perhaps firing off events to either the local bus (for other listeners) or the client bus (for UI updates).

# Running

To run locally, simply run the `eventbus` application.

Alternatively you can use foreman with `foreman start`

## CORS

If you are running the service under a different server (such as Apache or nginx) then you will need to make sure that the correct CORS headers are present when HTML files are requested.

The static HTML server response headers must include the following CORS headers:

```
Access-Control-Allow-Origin: eventbus-url
Access-Control-Allow-Headers: Content-Type
```

Where eventbus-url is replaced with the URL for the EventBus server.

# Security

* Currently the EventBus WebSocket service handles only unecrypted calls.
