package mill

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thingsplex/mill/model"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultBaseURL is mill api url
	baseURL = "https://api.millheat.com/"
	// applyAccessTokenURL is mill api to get access_token and refresh_token
	applyAccessTokenURL = baseURL + "share/applyAccessToken"
	// authURL is mill api to get authorization_code
	authURL = baseURL + "share/applyAuthCode"
	// refreshURL is mill api to update access_token and refresh_token
	refreshURL = baseURL + "share/refreshtoken"

	// deviceControlForOpenApiURL is mill api to controll individual devices
	deviceControlURL = baseURL + "uds/deviceControlForOpenApi"
	// getIndependentDevicesURL is mill api to get list of devices in unassigned room
	getIndependentDevicesURL = baseURL + "uds/getIndependentDevices"
	// selectDevicebyRoomURL is mill api to search device list by room
	selectDevicebyRoomURL = baseURL + "uds/selectDevicebyRoom"
	// selectHomeListURL is mill api to search housing list
	selectHomeListURL = baseURL + "uds/selectHomeList"
	// selectRoombyHomeURL is mill api to search room list by home
	selectRoombyHomeURL = baseURL + "uds/selectRoombyHome"
)

// Config is used to specify credential to Mill API
// AccessKey : Access Key from api registration at http://api.millheat.com. Key is sent to mail.
// SecretToken: Secret Token from api registration at http://api.millheat.com. Token is sent to mail.
// Username: Your mill app account username
// Password: Your mill app account password
type Config struct {
	ErrorCode  int    `json:"errorCode"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Success    bool   `json:"success"`

	Props struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		AccessKey   string `json:"access_key"`
		SecretToken string `json:"secret_token"`
	} `json:"props"`

	Data struct {
		AuthorizationCode string `json:"authorization_code"`
		AccessToken       string `json:"access_token"`
		RefreshToken      string `json:"refresh_token"`
		ExpireTime        int64  `json:"expireTime"`
		RefreshExpireTime int64  `json:"refresh_expireTime"`
	} `json:"data"`
}

// Client to make request to Mill API
type Client struct {
	configs      *model.Configs
	httpResponse *http.Response

	Data struct {
		Homes   []*Home   `json:"homeList"`
		Rooms   []*Room   `json:"roomList"`
		Devices []*Device `json:"deviceList"`
	} `json:"data"`
}

// Device is a mill heater
type Device struct {
	MaxTemperature       int    `json:"maxTemperature"`
	MaxTemperatureMsg    string `json:"maxTemperatureMsg"`
	ChangeTemperature    int    `json:"changeTemperature"`
	CanChangeTemp        int    `json:"canChangeTemp"`
	DeviceID             int64  `json:"deviceId"`
	DeviceName           string `json:"deviceName"`
	ChangeTemperatureMsg string `json:"changeTemperatureMsg"`
	Mac                  string `json:"mac"`
	DeviceStatus         int    `json:"deviceStatus"`
	HeaterFlag           int    `json:"heaterFlag"`
	SubDomainID          int    `json:"subDomainId"`
	ControlType          int    `json:"controlType"`
	CurrentTemp          int    `json:"currentTemp"`
}

type Home struct {
	HomeName         string      `json:"homeName"`
	IsHoliday        int         `json:"isHoliday"`
	HolidayStartTime int         `json:"holidayStartTime"`
	TimeZone         string      `json:"timeZone"`
	ModeMinute       int         `json:"modeMinute"`
	ModeStartTime    int64       `json:"modeStartTime"`
	HolidayTemp      int         `json:"holidayTemp"`
	ModeHour         int         `json:"modeHour"`
	CurrentMode      int         `json:"currentMode"`
	HolidayEndTime   int         `json:"holidayEndTime"`
	HomeType         interface{} `json:"homeType"`
	HomeID           int64       `json:"homeId"`
	ProgramID        int64       `json:"programId"`
}

type Room struct {
	MaxTemperature       int           `json:"maxTemperature"`
	IndependentDeviceIds []interface{} `json:"independentDeviceIds"`
	MaxTemperatureMsg    string        `json:"maxTemperatureMsg"`
	ChangeTemperature    int           `json:"changeTemperature"`
	ControlSource        string        `json:"controlSource"`
	ComfortTemp          int           `json:"comfortTemp"`
	RoomProgram          string        `json:"roomProgram"`
	AwayTemp             int           `json:"awayTemp"`
	AvgTemp              int           `json:"avgTemp"`
	ChangeTemperatureMsg string        `json:"changeTemperatureMsg"`
	RoomID               int64         `json:"roomId"`
	RoomName             string        `json:"roomName"`
	CurrentMode          int           `json:"currentMode"`
	HeatStatus           int           `json:"heatStatus"`
	OffLineDeviceNum     int           `json:"offLineDeviceNum"`
	Total                int           `json:"total"`
	IndependentCount     int           `json:"independentCount"`
	SleepTemp            int           `json:"sleepTemp"`
	OnlineDeviceNum      int           `json:"onlineDeviceNum"`
	IsOffline            int           `json:"isOffline"`
}

// NewClient create a handle authentication to Mill API
func (config *Config) NewClient(accessKey string, secretToken string, password string, username string) (string, string, string, int64, int64) {
	req, err := http.NewRequest("POST", authURL, nil)
	if err != nil {
		// handle err
		log.Debug("request error")
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Access_key", accessKey)
	req.Header.Set("Secret_token", secretToken)

	resp, err := http.DefaultClient.Do(req)
	processHTTPResponse(resp, err, config)

	authorizationCode := config.Data.AuthorizationCode

	// have authorization code, send new curl request to get tokens
	url := applyAccessTokenURL + "?password=" + password + "&username=" + username
	req, err = http.NewRequest("POST", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization_code", authorizationCode)

	resp, err = http.DefaultClient.Do(req)
	processHTTPResponse(resp, err, config)

	accessToken := config.Data.AccessToken
	refreshToken := config.Data.RefreshToken
	expireTime := config.Data.ExpireTime
	refreshExpireTime := config.Data.RefreshExpireTime

	defer resp.Body.Close()
	return authorizationCode, accessToken, refreshToken, expireTime, refreshExpireTime
}

// GetHomeList sends curl request to get list of homes connected to user
func (c *Client) GetHomeList(accessToken string) (*Client, error) {
	req, err := http.NewRequest("POST", selectHomeListURL, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Access_token", accessToken)

	resp, err := http.DefaultClient.Do(req)
	processHTTPResponse(resp, err, c)

	return c, nil
}

// GetRoomList sends curl request to get list of rooms by home
func (c *Client) GetRoomList(accessToken string, homeID int64) (*Client, error) {
	url := fmt.Sprintf("%s%s%d", selectRoombyHomeURL, "?homeId=", homeID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Access_token", accessToken)

	resp, err := http.DefaultClient.Do(req)
	processHTTPResponse(resp, err, c)
	return c, nil
}

// GetDeviceList sends curl request to get list of devices by room
func (c *Client) GetDeviceList(accessToken string, roomID int64) (*Client, error) {
	url := fmt.Sprintf("%s%s%d", selectDevicebyRoomURL, "?roomId=", roomID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Access_token", accessToken)

	resp, err := http.DefaultClient.Do(req)
	processHTTPResponse(resp, err, c)
	return c, nil
}

// Unmarshall received data into holder struct
func processHTTPResponse(resp *http.Response, err error, holder interface{}) error {
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	// check http return code
	if resp.StatusCode != 200 {
		//bytes, _ := ioutil.ReadAll(resp.Body)
		log.Debug("Bad HTTP return code ", resp.StatusCode)
		return fmt.Errorf("Bad HTTP return code %d", resp.StatusCode)
	}

	// Unmarshall response into given struct
	if err = json.NewDecoder(resp.Body).Decode(holder); err != nil {
		return err
	}
	return nil
}
