# SQS Mock

A library used to mimic sqs for local development

## Get Started

Clone this project into a go workspace. See https://golang.org/doc/code.html#Workspaces for details.

**Install Dependencies**
From inside the `sqsmock` directory run
```
$ go get
```

## Use
The library supports and sqs queue and it's corresponding worker.

Start the app by running:
```
$ go run sqsmock.go
```

You must include the url that the mock server is running on with the `--url` flag.

If you want messages forwarded to a worker, use the `--workerUrl` flag to specify the url your worker is running on.
