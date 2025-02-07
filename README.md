# Practice run

Simple chat server. After connecting over websockets, 
the following commands are available:

- `/create #<room>`: Create a new room
- `/join #<room>`: Join a room
- `/leave #<room>`: Leave a room
- `/msg #<room> <message>`: Send a message to a room

## Development

Start the server on port 8080:
```shell
  make start
``` 

Run the tests:
```shell
  make test
```

Generate mocks:
```shell
  make generate
```
