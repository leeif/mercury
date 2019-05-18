Mercury
====

What is Mercury
----
Mercury is a room based chat server which is aiming to help deveplors building a chat service in a fast way. Mercury is an independent service from your app which make you can take more concentration to your app. We are also planning to make Mercury to a distributed chat server for a more scalability.

How Mercury Work
----
(In Progress...)

Feature of Mercury
----
(In Progress...)

Simple Demo
----
```
~$ go get github.com/leeif/mercury
~$ cd $GOPATH/src/github.com/leeif/mercury
~$ dep ensure
```
Start mercury server in 127.0.0.1
```
~$ env=development go run ./
```

Start client 1 and 2
```
// Terminal 1
~$ env=development go run ./test --member=1

// Terminal 2
~$ env=development go run ./test --member=2
```

Send Message
```
// Terminal 1
> Nice to meet you!

// Termianl 2
> Have a good day!
```