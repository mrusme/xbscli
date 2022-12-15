xbscli
------

A command line interface for [xbsapi](https://github.com/mrusme/xbsapi) as well 
as the official [xBrowserSync](https://github.com/xbrowsersync/api) API.


## Build

```sh
go build .
```

## Run

```sh
xbscli \
  -s "https://xbsapi.myserver.com/api/v1" \
  -i $(pass show xbs/id) \
  -p $(pass show xbs/password) \
  -f pretty
```


