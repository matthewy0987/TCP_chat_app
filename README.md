# TCP_chat_app
A simple go server utilizing go routines and channels with a simple terminal based client with a version in elixir and go

### To use:
1. run the server
`go run server.go`

2. run the client
for the go client:
`go run chat_client.go`
for the elixir client:
`elixir chat_client.exs`

Run the server and clients on seperate terminals. To stop use cntrl + c

### Functionality:
#### client functions:
`/NICK` or `/N`
- registers a client under a nickname
- to send messages, a client must be registered

`/LIST` or `/L`
- lists out all currently registered users connected to the server

`/MSG <recipient>` or `/M <recipient>`
- Messages a client that is registered under the username
- `<recipients>` can be multiple users is done as such `/M user1,user2,user3 message to be sent`
- note: no spaces between users
- `/M * message to broadcast` the "*" broadcasts to all other clients

### Design Choice:
To handle concurrency and deal with multiple users, 3 initial processes are used.
The first is the main infinite loop that will handle TCP connections and create a go routine for each connection to listen for messages. The Second is a go routine that will handle the logic of the received messages from the previously created go routines. The third handles the sending of messages back to the clients.
The messages are passed along through two channels; one is for receiving and the other for sending.
Connections are kept track of in a Map.

