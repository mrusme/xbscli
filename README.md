xbscli
------

[![Static Badge](https://img.shields.io/badge/Donate-Support_this_Project-orange?style=for-the-badge&logo=buymeacoffee&logoColor=%23ffffff&labelColor=%23333&link=https%3A%2F%2Fxn--gckvb8fzb.com%2Fsupport%2F)](https://xn--gckvb8fzb.com/support/) [![Static Badge](https://img.shields.io/badge/Join_on_Matrix-green?style=for-the-badge&logo=element&logoColor=%23ffffff&label=Chat&labelColor=%23333&color=%230DBD8B&link=https%3A%2F%2Fmatrix.to%2F%23%2F%2521PHlbgZTdrhjkCJrfVY%253Amatrix.org)](https://matrix.to/#/%21PHlbgZTdrhjkCJrfVY%3Amatrix.org)

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
