FROM golang:1.11

WORKDIR $GOPATH/src/sqsmock
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 4242

# Looks like Go takes localhost literally and only binds to 127.0.0.1
# But inside a container need to bind to all addresses otherwise can't connect from outside
# Need to use sh -c to get shell expansion of the WORKER_URL environment variable
ENTRYPOINT [ "sh", "-c", "sqsmock --url 0.0.0.0:4242 --workerUrls $WORKER_URL" ]