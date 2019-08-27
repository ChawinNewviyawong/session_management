#sample go api server
use 
- <gin> framwork to backend server base golang
dependant
- <redis> database store manage session 
- <nodejs> backend server auth user & SDK for hyperledger fabric base javascript
- <expressDB> database store logs
connect 
- <angular> framwork to frontend

## go install lib
```
go get github.com/gin-contrib/cors github.com/gin-gonic/gin github.com/go-redis/redis
```

## How run gin api server
golang version 1.10.X ^
run command to terminal
```
go run server.go  redis.go  sessionmgmt.go  logger.go 
```

Example for run time 
```
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.(*operation).FirstPage-fm (5 handlers)
[GIN-debug] POST   /login                    --> main.(*operation).Login-fm (5 handlers)
[GIN-debug] POST   /getAllCars               --> main.(*operation).GetAllCars-fm (5 handlers)
[GIN-debug] Listening and serving HTTP on 0.0.0.0:3000
INFO 2019-08-27T16:41:56+07:00 |Start API Server 0.0.0.0:3000| "actor":"" "component":"sample_server" "function":"setupRouter"   "uuid":""
```