# kafka-websocket-proxy

Two-way proxy between the browser and a Kafka instance.

## Running

```bash
make run
```

or

```bash
make run-docker
```

## API

```javascript
const ws = WebSocket("/v1/broker/?topic=mytopic1,mytopic2,mytopic3")

ws.onmessage = (event) => {
    console.log("Message: " + event.data)
}

ws.send({
    topic: "mytopic2",
    data: "hello!",
})

ws.send({
    topic: "mytopic,mytopic3",
    data: "hello!",
})

```
