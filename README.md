Mercury
====

What is Mercury
----
Mercury is a room based chat prject which is aiming to help deveplors building a chat service in a fast way. Mercury is an independent service from your app which make you can take more concentration on your app. We are also planning to make Mercury to a distributed chat service for a more scalability.

How Mercury Work
----
The Mercury is designed into two part, **Rest API** server and **Websocket** server.

### What Rest API Server does
| API | Method | Description |
| ---- | :----: | :----: |
| /api/token | GET  | Get token of a user for a websocket connection |
| /api/room/add  | POST | add members into a chat room |

### What Websocket Server does
Connect to the server using the token.
```
ws://<ip>:<port>/ws/connect?token=xxxx
```
Merucy server receive a json format message. The client should send a json data to the server in the websocket connection.
#### send message to a chat room
```
{
  "type": 1,
  "rid": "<room id>",
  "text" "<message>",
}
```

#### get history message
```
{
  "type": 2,
  "msgid": "<start message id>",
  "offset" "<number of message before msgid>",
}
```


Configuration
---
### command-line flags

#### Server
Basic server configuration

| config | default value |
| ----- | :----: |
| --server.api.address | 127.0.0.1/32  |
| --server.ws.address | 0.0.0.0/0  |
| --server.port  | 6010 |

* Mercury only recive a IPv4 remote ip now.

#### Log

Logging in Mercury is using the interface provided by [go-kit](https://github.com/go-kit/kit/tree/master/log)

| config | available value |
| ---- | :----: |
| --log.format | json(default) \| logfmt   |
| --log.level  | info(default) \| warn \| error \| debug  |

#### Storage
You can choose from multiple back end storage to store the data in the Mercury.

MySQL

| config | default value |
| ---- | :----: |
| --mysql.host | ""  |
| --mysql.port | "3306"  |
| --mysql.user | "root"   |
| --mysql.password |  ""  |

### Configure File
You can also configure the Mercury through a config file using the https://github.com/jinzhu/configor

The YAML, JSON, TOML format of configure file are supported.

The default file path is `./mc.cnf.toml` and can be set by `--config.file` command line option.

Here is an example of configure file https://github.com/leeif/mercury/blob/master/config/config_file.go


Simple Demo
----
```
~$ go get github.com/leeif/mercury
~$ cd $GOPATH/src/github.com/leeif/mercury
~$ dep ensure
```
Start Mrcury server
```
~$ go run ./
```

Test client
```
// Terminal 1
~$ go run ./test --member=1

// Terminal 2
~$ go run ./test --member=2
```

Send Message
```
// Terminal 1
send> Hello World 1

// Termianl 2
send> Hello World 2
```

Contact Us (联系我们)
----
Slack : [#mushare-dev](https://mushare-dev.slack.com)
