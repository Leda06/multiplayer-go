# multiplayer-go

### Requirements

- go v1.23.2

#### Proxy

[proxy.cn](https://github.com/goproxy/goproxy.cn/blob/master/README.zh-CN.md)

```bash
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

### Running

```
go run cmd/entry/main.go
```

### Build

```
docker build -t multiplayer-go .
```
