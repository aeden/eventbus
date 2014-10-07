#!/usr/bin/env ruby

require 'websocket-eventmachine-client'
require 'json'
require 'date'
require 'net/http'

EVENT_BUS_SERVER_URL = URI('http://localhost:3001')
WEBSOCKET_URI = 'ws://localhost:3000/ws'

def authorization
  'token xxx'
end

# Create a new event in the EventBus.
def deliver(message)
  Net::HTTP.start(EVENT_BUS_SERVER_URL.host, EVENT_BUS_SERVER_URL.port) do |http|
    request = Net::HTTP::Post.new EVENT_BUS_SERVER_URL
    request.body = message
    request.content_type = "application/json"
    request["Authorization"] = authorization
    response = http.request request # Net::HTTPResponse object
    puts response.inspect

    # TODO: If the response is not 200 then there needs to be retry logic
  end
end

# Connect to the websocket server and respond to inbound websocket messages.
def connect(retry_count = 0)
  if retry_count > 20
    raise "Max connection attempts reached"
  end

  if retry_count > 0
    puts "Connecting in #{retry_count} seconds"
    sleep(retry_count)
  end

  puts "Opening websocket connection"
  ws = WebSocket::EventMachine::Client.connect(:uri => WEBSOCKET_URI)

  ws.onopen do
    puts "Connected"
    retry_count = 0
    ws.send(JSON.generate({action: 'authenticate', credentials: authorization}))
  end

  ws.onmessage do |msg, type|
    puts "Received message: #{msg}"
    event = JSON.parse(msg)

    case event['name']
    when 'check-domain'
      # Normally this is the part where we'd check the domain name at the specific
      # registry to see if it is available. Sleep for 1 second to simulate.
      sleep(1)

      results = event['data'].map do |domain_name|
        {name: domain_name, availability: 'available'}
      end

      message = JSON.generate({name: 'check-domain-completed', data: results, context: event['context']})
      deliver(message)
    when 'register-domain'
      received_data = event['data']

      # Normally this is the part where we'd register the domain name at the specific
      # registry. Sleep for 2 seconds to simulate.
      sleep(2)

      results = {
        name: received_data['name'],
        registered: true,
        expiration: (Date.today + 365).rfc3339
      }

      message = JSON.generate({name: 'register-domain-completed', data: results, context: event['context']})
      deliver(message)
    end
  end

  ws.onclose do |code, reason|
    puts "Disconnected with status code: #{code}"
    connect(retry_count + 1)
  end
end

# Run EventMachine
EM.run do
  connect
end
