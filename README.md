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
message broker between them: NATS
subjects are chat groups
API call to have username, chat group REST API


Backend flow:

Browsers (Subscribers) <--- SSE (chat group in params) --- Server 
Server subscribers on behalf of browser <--- NATS (subject = chat group)

Browsers (Producer) ---> Server --- REST API (chat group, username, message) 
Server produces on behalf of browser ----> NATS (subject = chat group)