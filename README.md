# MNZ chain

a toy blockchain implementation, for learning purposes.

## Background

inspired from awesome

[Learn Blockchains by Building One](https://github.com/dvf/blockchain)

## Building

```
$ go build 
```

## Usage

start instances(nodes) like this

You can start as many nodes as you want with the following command

```
./mnzchain -port 8001
./mnzchain -port 8002
./mnzchain -port 8003
```

## Json Endpoints


get full blockchain

* `GET 127.0.0.1:8001/chain`

mine a new block

* `GET 127.0.0.1:8001/mine`

Adding a new transaction

* `POST 127.0.0.1:8001/transactions/new`

* __Body__:

  ```json
  {
    "sender": "15PP7EWmqnctqpxjxXxi8Bh6kWVdhkgWzV",
    "recipient": "15X8skrTNBR8TKxBMa6axPL9iGTv6bG9mA",
    "amount": 1000
  }
  ```

Register a node in the blockchain network

* `POST 127.0.0.1:8001/nodes/register`

* __Body__:

  ```json
  {
    "nodes": [
        "http://127.0.0.1:8002",
        "http://127.0.0.1:8003"
    ]
}
  ```

Resolving Blockchain

* `GET 127.0.0.1:8001/nodes/resolve`


## TODO

- a lot!

