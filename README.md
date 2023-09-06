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

## Supported formats (-f)

* json (default)
* pretty - formatted text
* html - produces a basic HTML file similar to browser export of bookmarks as html. Handy as an input to [static-marks](https://darekkay.com/static-marks/)

## Docker

An alternate way to run xbscli

### Build

```sh
 docker build -t xbscli .
```

### Run

```sh
 docker run --rm xbscli \
  -s "https://xbsapi.myserver.com/api/v1" \
  -i $(pass show xbs/id) \
  -p $(pass show xbs/password) \
  -f pretty
```
