# Streamfiles
An example client/server grpc application that can stream a files and store it.

## Requirements
Make sure you have `protoc` binary in your $PATH, as its used to generate from protobuf files.

## Usage
```$ make build
$ ./streamfiles server # start a grpc server
$ ./streamfiles client localhost:7777 <filename> # uploads filename to grpc server
```

## Data Directory
This is the location you want files uploaded. Defaults to `data`