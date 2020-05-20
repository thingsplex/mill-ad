# Futurehome Mill Adapter - WORK IN PROGRESS

Adapter connects with api and retrieves list of homes, rooms and devices connected to your mill user.

Testing: 
Configure mqtt uri, username and password in `testdata/data/config.json`to your Futurehome gateway. 
Use `make run` to run the program. 


Configure this message and send to your mqtt broker (e.g. through thingsplex)

`Topic`
```
pt:j1/mt:evt/rt:dev/rn:mill/ad:1/sv:mill/ad:1
```

`Payload`
```json
    {
    "type": "cmd.auth.login",
    "serv": "mill",
    "props": {
        "username": "appUsername",
        "password": "appPassword",
        "access_key": "futurehome_access_key",
        "secret_token": "futurehome_secret_token"
        },
    }
```

The program will then send a http request to the mill api which returns a unique authorization code. (work in progress->) The program will then send a new http request to the api to receive access_token and refresh_token. When access_token is received you can get information about all homes, rooms and devices, as well as setting temperature and/or switching devices on/off. 