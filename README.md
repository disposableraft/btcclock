
In the following example I'm going to build a Bitcoin stock ticker using Go, Docker and Minikube.

The stock ticker is a simple TCP server that listens on port 3000 and updates clients every two seconds.

```go
func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}
```

In order for the server to connect to multiple clients at the same time, `handleConn` is a go routine. Otherwise, the currently connected client would block all other clients.

```go
func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		quote, err := quote.Get("BTC-USD")
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.WriteString(conn, fmt.Sprintf("%v: $%v\r", time.Now().Format("15:04:05"), quote.RegularMarketPrice))
		if err != nil {
			return
		}
		time.Sleep(2 * time.Second)
	}
}
```

This connection handler gets a stock quote from Yahoo, sends it to the TCP server along with a formated time and then sleeps for two seconds before starting over.

The handler only executes while a client is connected. When a client disconnects, `io.WriteString` returns an error, the handler returns and the connection is closed by `defer conn.Close()`.

By including a carriage return, `\r`, after the stock quote, we tell the terminal's cursor to move back to the beginning of the line so that the next stock quote overwrites the current one.

```dockerfile
FROM alpine:3.14
RUN apk add --no-cache go
WORKDIR /app
COPY . .
RUN go build main.go
CMD ["go", "run", "btcclock"]
EXPOSE 3000
```

We ask Docker to do the following:

- build an image using the Alpine distribution of Linux, 
- install Go, 
- define a working directory, 
- copy all files from the current host directory to that working directory, 
- compile the Go app,
- run the executable,
- and, finally, expose port 3000.

Build the image: `docker build -t btcclock:0.0.1 .`.

Start the container: `docker run -dp 3000:3000 btcclock:0.0.1`.

Visit the ticker: `nc localhost 3000`.

Stop the ticker: `docker stop <container_id>`

(Find the container ID or name with `docker container ls`)

Now to deploy the image to the local Kubernetes cluster using MiniKube.

Start up Minicube: `minikube start`

Build the docker application within the Minikube registry: `minikube image build -t btcclock:0.0.1 .`

Create the Kubernetes deployment: `kubectl apply -f ./deployment.yaml`

The deployment starts one pod using the `btcclock:0.0.1` image that we built within Minikube's registry. The secret to running this unhosted Docker image is setting `imagePullPolicy` to `Never`.

Create a service: `kubectl expose deployment btcclock-deployment --type=NodePort --port=3000`

Forward the port: `kubectl port-forward service/btcclock-deployment 3000`

Visit the ticker: `nc localhost 3000`.
