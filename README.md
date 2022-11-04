# G5 MsPacMan Backend Service

## Running the application
First, clone this repo to your workspace. The important files are only `.env` and `docker-compose.yml`.
```
docker compose up -d
```
Then, the application will serve on port `8000`.

## Example for accessing Socket.IO
You can emit the message to the sentiment analysis service using `message` event.
```
var msg = "PacMan noob"
socket.emit('message', {"message": msg})
```
You can also listen on `result` event for receiving a sentiment result. (optional)
```
{
    "result": "negative",
    "confidence": 0.81
}
```
Whenever the `/poll` API is called, the result will also broadcast through `message` event.
```
socket.on('message', (data) => {
    // data = {'cheer': 0.75, 'jeer': 0.25}
    // do some stuff
})
```

## Example for accessing REST API
You can poll the sentiment ratio by `/poll` with `GET` method.
```
{
    "cheer": 0.75,
    "jeer": 0.25
}
```
Whenever you call this API, the Socket.IO will broadcast to the its connected client and the sentiment will be reset. \
You can call this API silently without reset the sentiment result by using `?silent=true`.
