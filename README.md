# Blue

Blue is a little service we built to test BLE devices. It helps you to access BLE services and characteristics through simple and friendly HTTP API.

## Build & Run

Blue is written in Go and has very little amount of dependencies only necessary for service to work. Obviously, you can only run Blue on the computer which has active Bluetooth device.

In order to build and run Blue you need Go >= 1.22.0. At least, we used this one to develop it. If you have it just run the following:

```shell script
go run main.go
```

It will just start a web server but no scanning is performed yet. In order to start scanning you can send request like this one:

```shell script
curl --location 'http://localhost:8080/start' \
--header 'Content-Type: application/json' \
--data '{
    "filter": {
        "kind": "one_of",
        "props": {
            "filters": [
                {
                    "kind": "name",
                    "props": {
                        "pattern": "AXA:([A-Z0-9]+)"
                    }
                },
                {
                    "kind": "service",
                    "props": {
                        "uuid": "98baba41-3e32-4fec-976f-634a8efc6010"
                    }
                }
            ]
        }
    }
}'
```

Using filters you can tell Blue which devices you're interested in and which one it should ignore. Here are kinds available:

1. `name` - only keeps devices matching regular expression given as `pattern` property to the filter
2. `service` - only keeps devices which have certain service speified by it's UUID in props.
3. `one_of` - only keeps devices which are matching one of the filters given to this on in props, basically it is logical OR
4. `all_of` - only keeps devices which are matching all of the filters given to this one in props, basically it is logical AND

You can always stop scanning using this request:

```shell script
curl --location --request POST 'http://localhost:8080/stop'
```

And start it again with another set of filters if you wish so. Next, check out what devices were detected:

```shell script
curl --location 'http://localhost:8080/devices'
```

Select any device you like and try to connect to it. For that take `address` and put it into the following request:

```shell script
curl --location 'http://localhost:8080/devices/b82ff1a8-e9ca-53f6-9392-01580af5710e'
```

If connection is successful as a response you should get a list of services and characteristics advertised by the device. Now  try to read some specific characteristic:

```shell script
curl --location 'http://localhost:8080/devices/91a2d34c-e89b-0b86-7824-529d1efb69a4/services/00001523-e513-11e5-9260-0002a5d5c51b/characteristics/00001524-e513-11e5-9260-0002a5d5c51b'
```

Because exact data type of this specific characteristic is not known Blue is trying to read it's value into 32 bytes buffer and print it out.

That's basically it. Have fun!

## License

Blue is released under the [MIT License](https://opensource.org/license/MIT).
