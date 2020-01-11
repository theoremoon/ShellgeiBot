# ShellGeiBot

[![Go Report Card](https://goreportcard.com/badge/github.com/theoldmoon0602/ShellgeiBot)](https://goreportcard.com/report/github.com/theoldmoon0602/ShellgeiBot)

## Twitter

- [@minyoruminyon](https://twitter.com/minyoruminyon)

## Specification

- https://furutsuki.hatenablog.com/entry/2018/07/13/221806

## Official Docker Image

- DockerHub: https://hub.docker.com/r/theoldmoon0602/shellgeibot
- GitHub: https://github.com/theoldmoon0602/ShellgeiBot-Image

## Development

* go version go1.12.7

Build shellgei bot.

```bash
make build
```

Testing.

```bash
make test
```

## Run

An example of TwitterConfig.json

```json
{
	"ConsumerKey": "<your twitter app's consumer key>",
	"ConsumerSecret": "<your twitter app's consumer secret>",
	"AccessToken": "<your account's access token>",
	"AccessSecret": "<your account's access secret>"
}
```

An example of ShellgeiConfig.json

```js
{
	"dockerimage": "theoldmoon0602/shellgeibot:master",
	"timeout": "20s",  // timeout
	"workdir": ".",    // where to make temporary directories
	"memory": "100M",  // max memory size of docker container
	"mediasize": 250,  // max media size to be able to creaate
	"tags": ["シェル芸", "危険シェル芸", "ゆるシェル"]  // trigger tags
}
```


```
<Usage>./ShellgeiBot: TwitterConfig.json ShellgeiConfig.json | -test ShellgeiConfig.json script
```

## 3rd parties

- SGWeb (https://shellgei-web.net) by [@kekeho](https://github.com/kekeho): https://github.com/kekeho/SGWeb
- websh (https://websh.jiro4989.com/) by [@jiro4989](https://github.com/jiro4989): https://github.com/jiro4989/websh

## Author

theoldmooon0602

## LICENSE

Apache License

