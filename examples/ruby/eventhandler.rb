#!/usr/bin/env ruby

require 'websocket-eventmachine-client'
require 'json'

EM.run do

  ws = WebSocket::EventMachine::Client.connect(:uri => 'ws://localhost:3000/ws')

  ws.onopen do
    puts "Connected"
  end

  ws.onmessage do |msg, type|
    puts "Received message: #{msg}"
    event = JSON.parse(msg)

    case event['name']
    when 'check.domain'
      results = event['data'].map do |domain_name|
        {name: domain_name, availability: 'available'}
      end
      message = JSON.generate({name: 'check.domain.result', data: results})
      ws.send message
    end
  end

  ws.onclose do |code, reason|
    puts "Disconnected with status code: #{code}"
  end

end
