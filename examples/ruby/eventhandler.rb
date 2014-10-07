#!/usr/bin/env ruby

require 'websocket-eventmachine-client'
require 'json'
require 'date'

def connect(connect_delay = 0)
  if connect_delay > 20
    raise "Max connection attempts reached"
  end

  if connect_delay > 0
    puts "Connecting in #{connect_delay} seconds"
    sleep(connect_delay)
  end

  puts "Opening websocket connection"
  ws = WebSocket::EventMachine::Client.connect(:uri => 'ws://localhost:3000/ws')

  ws.onopen do
    puts "Connected"
    connect_delay = 0
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

      message = JSON.generate({name: 'check-domain-completed', data: results})
      ws.send(message)
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

      message = JSON.generate({name: 'register-domain-completed', data: results})
      ws.send(message)
    end
  end

  ws.onclose do |code, reason|
    puts "Disconnected with status code: #{code}"
    connect(connect_delay + 1)
  end
end

EM.run do
  connect
end
