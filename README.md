# GMG Green mountain grill go library

Green mountain grill library for interacting with a wifi green mountain grill.

## Usage

This is a go library so use 

```
go get github.com:rob121/gmg
```

and then import as the same

## Features

This library allows you to get information from the grill, to poll this data on a user set schedule and to add event triggers when certains events occur

## Events

* Grill is ready - when the initial temp is reached (150)
* Grill is to temp - when the grill hits its set temp
* Grill is low on pellets
* Grill is done cooling off and ready to be turned off

