Tradologics Golang SDK
======================

This is the initial version of Tradologics' Golang SDK.

At the moment, it only supports a wrapper for the `net/http` library that will automatically:

- prepend the full endpoint url to your calls
- attach your token to the request headers
- add `datetime` to your order when in backtesting mode


Install
------------------

```sh
go get github/tradologics/go-sdk
```
