package mill

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	httpResponse *http.Response
	Dc           *DeviceCollection
	Rc           *RoomCollection
	Hc           *HomeCollection
}

// DeviceCollection hold all devices from mill account
type DeviceCollection struct {
	Body struct {
		Devices []*Device `json:"deviceList"`
	}
}

// RoomCollection hold all rooms from mill account
type RoomCollection struct {
	Body struct {
		Rooms []*Room `json:"roomList"`
	}
}

type HomeCollection struct {
	Body struct {
		Homes []*Home `json:"homeList"`
	}
}

// Device is a mill heater
type Device struct {
	maxTemperature       *int32 `json:"maxTemperature"`
	maxTemperatureMsg    string `json:"maxTemperatureMsg"`
	changeTemperature    *int32 `json:"changeTemperature"`
	canChangeTemp        bool   `json:"canChangeTemp"`
	deviceId             *int32 `json:"deviceId"`
	deviceName           string `json:"deviceName"`
	changeTemperatureMsg string `json:"changeTemperatureMsg"`
	mac                  string `json:"mac"`
	deviceStatus         bool   `json:"deviceStatus"`
	heaterFlag           string `json:"heaterFlag"`
	subDomainId          *int32 `json:"subDomainId"`
	controlType          bool   `json:"controlType"`
	currentTemp          *int32 `json:"currentTemp"`
}

// Room is a room containing one or more mill heaters
type Room struct {
	maxTemperature       *int32 `json:"maxTemperature"`
	independentDeviceIds *int32 `json:"independentDeviceIds"`
	maxTemperatureMsg    string `json:"maxTemperatureMsg"`
	changeTemperature    *int32 `json:"changeTemperature"`
	controlSource        string `json:"controlSource"`
	comfortTemp          *int32 `json:"comfortTemp"`
	roomProgram          string `json:"roomProgram"`
	awayTemp             *int32 `json:"awayTemp"`
	avgTemp              *int32 `json:"avgTemp"`
	changeTemperatureMsg string `json:"changeTemperatureMsg"`
	roomId               *int32 `json:"roomId"`
	roomName             string `json:"roomName"`
	currentMode          *int32 `json:"currentMode"`
	heatStatus           bool   `json:"heatStatus"`
	offLineDeviceNum     *int32 `json:"offLineDeviceNum"`
	total                *int32 `json:"total"`
	independentCount     *int32 `json:independentCount"`
	sleepTemp            *int32 `json:sleepTemp"`
	onlineDeviceNum      *int32 `json:onlineDeviceNum"`
	isOffline            *int32 `json:isOffline"`
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

// GetAuth sends curl request to get authorization_code
func (config *Config) GetAuth(accessKey string, secretToken string) string {
	req, err := http.NewRequest("POST", authURL, nil)
	if err != nil {
		// handle err
		log.Debug("request error")
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Access_key", accessKey)
	req.Header.Set("Secret_token", secretToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		log.Debug("do error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle err
		log.Debug("read error")
	}
	conf := Config{}
	jsonErr := json.Unmarshal(body, &conf)
	if jsonErr != nil {
		// handle err
		log.Debug("json error")
	}

	defer resp.Body.Close()
	return conf.Data.AuthorizationCode
}

// GetAccessAndRefresh sends curl request to get access_token and refresh_token, as well as expireTime and refresh_expireTime
func (config *Config) GetAccessAndRefresh(authorizationCode string, password string, username string) (string, string, int64, int64) {
	url := applyAccessTokenURL + "?password=" + password + "&username=" + username
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization_code", authorizationCode)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		log.Debug("do error")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle err
		log.Debug("read error")
	}

	conf := Config{}
	jsonErr := json.Unmarshal(body, &conf)
	if jsonErr != nil {
		// handle err
		log.Debug("json error")
	}

	accessToken := conf.Data.AccessToken
	refreshToken := conf.Data.RefreshToken
	expireTime := conf.Data.ExpireTime
	refreshExpireTime := conf.Data.RefreshExpireTime

	defer resp.Body.Close()
	return accessToken, refreshToken, expireTime, refreshExpireTime
}

// Reading json does not work
func (c *Client) GetHomeList(accessToken string) (*HomeCollection, error) {
	req, err := http.NewRequest("POST", selectHomeListURL, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Access_token", accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}

	log.Debug(resp.Body)

	if err = processHTTPResponse(resp, err, c.Hc); err != nil {
		return nil, err
	}

	return c.Hc, nil
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
