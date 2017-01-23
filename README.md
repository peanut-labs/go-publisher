# Peanut Labs GoLang package
[![Build Status](https://travis-ci.org/peanut-labs/go-publisher.svg?branch=master)](https://travis-ci.org/peanut-labs/go-publisher)
[![GoDoc](https://godoc.org/github.com/peanut-labs/go-publisher?status.svg)](https://godoc.org/github.com/peanut-labs/go-publisher)


GoLang utility provided as an integration utility for Peanut Labs publishers.

Only supports the Reward Center integration for now.

### Integration Guide
[Publisher Integration Guide](http://peanut-labs.github.io/publisher-doc/)

### Installation

```
go get github.com/peanut-labs/go-publisher
```

### Usage

To generate PL User ID or the Reward Center URL:

```go
import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	peanutlabs "github.com/peanut-labs/go-publisher"
)

//See publisher guide on how to get AppId and Keys
appID := 1
secKey := "123"
transKey := "123"
pl, err := peanutlabs.New(appID, secKey, transKey)
if err != nil {
	fmt.Println(err)
}

endUserID := "john_doe"

//To generate Peanut Labs User ID
uid, _ := pl.GenerateUserID(endUserID)
fmt.Println(uid)

//The URL returned should be used for an iFrame
url, _ := pl.GenerateRewardCenterURL(endUserID)
fmt.Println(url)
```


To process Reward Notification

```go
func HandleNotification(w http.ResponseWriter, r *http.Request) {
	pl, err := peanutlabs.New(appID, secKey, secKey)
	if err != nil {
		fmt.Fprintf(w, "%d", peanutlabs.NotificationResponseFailure)
		return
	}
	//pass the http.Request 
	rn, err := pl.ProcessRewardNotification(r)
	if err != nil {
		fmt.Fprintf(w, "%d", peanutlabs.NotificationResponseFailure)
		return
	}

	//reward notification
	fmt.Println("Amount: %s", rn.Amount)

	fmt.Fprintf(w, "%d", peanutlabs.NotificationResponseSuccess)
}
```


### Status
This package only supports Reward Center at the moment but the newer versions will support other PL integrations as well.

Contact: publishers@peanutlabs.com