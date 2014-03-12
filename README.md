cli_chat
========

cli_chat is a simple CLI based client and server for tcp based chat connection. Both
server and client written in Go. You can use telnet too.

## Install

Clone this repository.

~~~text
$ cd /path/to/cloned/repo
$ make build
~~~

## Run the server

~~~text
./bin/server
~~~

## Use the client to connect to server

~~~text
./bin/client
~~~

## Use telnet to connect to server

~~~text
telnet 0.0.0.0 8888
~~~

## TODO

* Creates a nice client interface. Using [termbox](https://github.com/nsf/termbox-go) maybe?
* Type `Message` should contains more attribute, like `action`. There will be message
  formatter on the client code that check `Message.action`. A good use case would be `action`
  to retrieve all connected clients, to quit, and to kick another user.
* Allows multiple rooms so that upon initial connection user can join open room after his/her nick
  is registered.

## License

MIT License - see [LICENSE](LICENSE) file.
