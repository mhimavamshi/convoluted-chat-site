# Convoluted Chat Site

## Idea
Let's make a chat site that doesn't use websockets, but also doesn't use long/short polling.
Alright. How? 
What if we have a SSE endpoint to get messages and HTTP API to send messages? 
Kinda like Pub/Sub.
Well, that's convoluted. ...And gives rise to unecessary problems.
So let's make unecessary solutions.

## Architecture/Flow
htmx for front-end
each browser is a consumer over SSE
each message REST call is a producer of messages
message broker between them: I will just make a very dumb simple TCP topic-based pub-sub server
(might iterate and improve this, but need to get the initial one going)
subjects are chat groups
API call to have username, chat group REST API


Backend flow:

Browsers (Subscribers) <--- SSE (chat group in params) --- Server 
Server subscribers on behalf of browser <--- Pub/Sub (subject = chat group)

Browsers (Producer) ---> Server --- REST API (chat group, username, message) 
Server produces on behalf of browser ----> Pub/Sub (subject = chat group)


that TCP server takes a connection and a command:
- PUB: just parses the text CHATROOM MESSAGE and we iterate over the mapped connections to that chatroom and send a message
- SUB: just parses CHATROOM and creates a connection mapped to that chatroom
- EXIT: closes the connection that sent it (must be SUB connection)

we will handle the failure case of SUB connection being lost (for now just remove from hashmap) and add PING/PONG 

client library just does a light abstraction over these commands