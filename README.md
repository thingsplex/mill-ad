# Futurehome Mill Adapter - WORK IN PROGRESS

Adapter connects with api and retrieves list of homes, rooms and devices connected to your mill user.

For testing: 
Configure mqtt url, username and password in `testdata/data/config.json`to your Futurehome gateway. 
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

The program will then send a http request to the mill api which returns a unique authorization code. The program will then send a new http request to the api to receive access_token, refresh_token, expireTime and refresh_expireTime. When access_token is received you can get information about all homes, rooms and devices, as well as setting temperature and/or switching devices on/off. 

***

## Services and interfaces
#### Service name
`thermostat`
#### Interfaces
Type | Interface               | Value type | Description
-----|-------------------------|------------|------------------
in   | cmd.mode.get_report     | null       |
in   | cmd.mode.set            | string     |  set thermostat mode
out  | evt.mode.report         | string     |
-|||
in   | cmd.setpoint.get_report | string     | value is a set-point type
in   | cmd.setpoint.set        | str_map    | val = {"type":"heat", "temp":"21.5", "unit":"C"}
out  | evt.setpoint.report     | str_map    | val = {"type":"heat", "temp":"21.5", "unit":"C"}

#### Service name
`sensor_temp`
#### Interfaces
Type | Interface               | Value type | Description
-----|-------------------------|------------|------------------
in   | cmd.sensor.get_report   | null       | 
in   | evt.sensor.report       | float      | measured temperature



