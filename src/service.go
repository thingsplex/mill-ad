package main

import (
	"flag"
	"fmt"
	"reflect"
	"time"

	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/discovery"
	log "github.com/sirupsen/logrus"
	mill "github.com/thingsplex/mill/millapi"
	"github.com/thingsplex/mill/model"
	"github.com/thingsplex/mill/router"
	"github.com/thingsplex/mill/utils"
)

func main() {
	var workDir string
	flag.StringVar(&workDir, "c", "", "Work dir")
	flag.Parse()
	if workDir == "" {
		workDir = "./"
	} else {
		fmt.Println("Work dir ", workDir)
	}
	appLifecycle := model.NewAppLifecycle()
	configs := model.NewConfigs(workDir)
	states := model.NewStates(workDir)
	err := configs.LoadFromFile()
	if err != nil {
		fmt.Print(err)
		panic("Can't load config file.")
	}
	err = states.LoadFromFile()
	if err != nil {
		fmt.Print(err)
		panic("Can't load state file.")
	}
	mill := mill.Client{}

	utils.SetupLog(configs.LogFile, configs.LogLevel, configs.LogFormat)
	log.Info("--------------Starting mill----------------")
	log.Info("Work directory : ", configs.WorkDir)
	appLifecycle.PublishEvent(model.EventConfiguring, "main", nil)

	mqtt := fimpgo.NewMqttTransport(configs.MqttServerURI, configs.MqttClientIdPrefix, configs.MqttUsername, configs.MqttPassword, true, 1, 1)
	err = mqtt.Start()
	responder := discovery.NewServiceDiscoveryResponder(mqtt)
	responder.RegisterResource(model.GetDiscoveryResource())
	responder.Start()

	fimpRouter := router.NewFromFimpRouter(mqtt, appLifecycle, configs, states)
	fimpRouter.Start()

	appLifecycle.SetConnectionState(model.ConnStateDisconnected)
	if configs.IsConfigured() && err == nil {
		appLifecycle.SetConfigState(model.ConfigStateConfigured)
		appLifecycle.SetAppState(model.AppStateRunning, nil)
		appLifecycle.SetConnectionState(model.ConnStateConnected)
	} else {
		appLifecycle.SetConfigState(model.ConfigStateNotConfigured)
		appLifecycle.SetAppState(model.AppStateNotConfigured, nil)
		appLifecycle.SetConnectionState(model.ConnStateDisconnected)
	}

	if configs.IsAuthenticated() && err == nil {
		appLifecycle.SetAuthState(model.AuthStateAuthenticated)
	} else {
		appLifecycle.SetAuthState(model.AuthStateNotAuthenticated)
	}
	//------------------ Sample code --------------------------------------

	msg := fimpgo.NewFloatMessage("evt.sensor.report", "temp_sensor", float64(23), nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeDevice, ResourceName: "mill", ResourceAddress: "1", ServiceName: "temp_sensor", ServiceAddress: "300"}
	mqtt.Publish(&adr, msg)
	if err != nil {
		log.Error("Can't connect to broker. Error:", err.Error())
	} else {
		log.Info("Connected")
	}
	appLifecycle.SetAppState(model.AppStateRunning, nil)
	//------------------ Sample code --------------------------------------

	for {
		appLifecycle.WaitForState("main", model.AppStateRunning)
		log.Info("Starting ticker")

		ticker := time.NewTicker(time.Duration(configs.PollTimeMin) * time.Minute)
		for ; true; <-ticker.C {
			states.DeviceCollection, states.RoomCollection, states.HomeCollection, states.IndependentDeviceCollection = nil, nil, nil, nil
			states.HomeCollection, states.RoomCollection, states.DeviceCollection, states.IndependentDeviceCollection = mill.UpdateLists(configs.Auth.AccessToken, states.HomeCollection, states.RoomCollection, states.DeviceCollection, states.IndependentDeviceCollection)
			states.SaveToFile()
			states.LoadFromFile()
			for device := range states.DeviceCollection {
				device := reflect.ValueOf(states.DeviceCollection[device])
				deviceId := fmt.Sprintf("%v", device.Interface().(map[string]interface{})["deviceId"])
				currentTemp := device.Interface().(map[string]interface{})["currentTemp"]

				tempVal := currentTemp
				props := fimpgo.Props{}
				props["unit"] = "C"

				adr := &fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeDevice, ResourceName: model.ServiceName, ResourceAddress: "1", ServiceName: "sensor_temp", ServiceAddress: deviceId}
				msg := fimpgo.NewMessage("evt.sensor.report", "sensor_temp", fimpgo.VTypeFloat, tempVal, props, nil, nil)
				mqtt.Publish(adr, msg)

				setpointTemp := fmt.Sprintf("%f", device.Interface().(map[string]interface{})["holidayTemp"])

				if setpointTemp != "0" {
					setpointVal := map[string]interface{}{
						"type": "heat",
						"temp": setpointTemp,
						"unit": "C",
					}
					adr = &fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeDevice, ResourceName: model.ServiceName, ResourceAddress: "1", ServiceName: "thermostat", ServiceAddress: deviceId}
					msg = fimpgo.NewMessage("evt.setpoint.report", "thermostat", fimpgo.VTypeStrMap, setpointVal, nil, nil, nil)
					mqtt.Publish(adr, msg)
				}
				// -----------------------------------------------------------------------------------------------
				states.SaveToFile()
			}
		}
		appLifecycle.WaitForState(model.AppStateNotConfigured, "main")
	}

	mqtt.Stop()
	time.Sleep(5 * time.Second)
}
