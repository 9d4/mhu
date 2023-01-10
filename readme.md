## Mikrotik Hotspot User

### Build Frontend

```shell
$ npm install
$ npm run build
```

Output located in views/dist.

### Build Backend

``` shell
$ go get
$ go build # or go run . 
```


### Runtime

Optional Environments:
- ADDRESS (default: 0.0.0.0:4500): Where server listen on
- MIKROTIK_URI (default:https://admin:@192.168.104.12/rest/): mikrotik rest address.


Notes:

- Set mikrotik certificate. READ [Mikrotik REST API](https://help.mikrotik.com/docs/display/ROS/REST+API)

- Mikrotik certificate's Subject Alt. Name must be available.
