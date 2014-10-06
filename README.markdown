# Overview

## Frontend

The front end is HTML and JavaScript. It relies solely on events as the mechanism for dealing with input and output. For example, a form is filled in an the submit button is pressed. This results in an event. This event is handled by multiple listeners: one listener updates the UI, the other forwards the event to a remote event bus. When a listener on the remote event bus handles the event it will fire a new event, directed towards the client's event bus (via websockets for a browser) where a listener handles the event and again updates the UI. 

## Backend

The backend is purely functional event handlers. A listener handles the event, perhaps firing off events to either the local bus (for other listeners) or the client bus (for UI updates).

# Running

To run locally, simply run the `eventbus` application.

## CORS

If you are running the service under a different server (such as Apache or nginx) then you will need to make sure that the correct CORS headers are present when HTML files are requested.
