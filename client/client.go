package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

type rclient struct {
	c *http.Client
	config *viper.Viper
}

const urlKey = "remoteit_url"
const devKeyHeader = "developerkey"
const tokenKeyHeader = "token"

type LoginRequest struct {
	APIKey string `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status string `json:"status"`
	Token string `json:"token"`
	Email string `json:"email"`
	GUID string `json:"guid"`
	ServiceToken string `json:"service_token"`
	ServiceLevel string `json:"service_level"`
	StoragePlan string `json:"storage_plan"`
	SecondaryAuth string `json:"secondary_auth"`
	APIKey string `json:"apikey"`
	AuthToken string `json:"auth_token"`
	AuthExpiration int `json:"auth_expiration"`
	ServiceAuthhash string `json:"service_authhash"`
	DeveloperPlan string `json:"developer_plan"`
	PortalPlan string `json:"portal_plan"`
	PortalPlanExpiration string `json:"portal_plan_expires"`
	ServiceFeatures string `json:"service_features"`
}

type BaseReqeust struct {
	APIKey string `json:"-"`
	Token string `json:"-"`
}

type ListDevicesRequest struct {
	BaseReqeust
}

type DeviceEntry struct {
	Address string `json:"deviceaddress"`
	Alias string `json:"devicealias"`
	Username  string `json:"deviceusername"`
	Type  string `json:"devicetype"`
	LastIP string `json:"devicelastip"`
	ServiceTitle string `json:"servicetitle"`
	WebEnabled string `json:"webenabled"`
	//WebURI string `json:"weburi"`
	LocalURI string `json:"localuri"`
}

type ListDevicesResponse struct {
	Status string `json:"status"`
	Devices []*DeviceEntry `json:"devices"`
}

type ConnectRequest struct {
	BaseReqeust
	DeviceAddress string `json:"deviceaddress"`
	HostIP string `json:"hostip"`
	Wait string `json:"wait"`
}

type ConnectResponse struct {
	Status string `json:"status"`
	Connection ConnectionInfo `json:"connection"`
}

type ConnectionInfo struct {
	DeviceAddress string `json:"deviceaddress"`
	ExpirationInSeconds string `json:"expirationsec"`
	ImageIntervalMS string `json"imageintervalms"`
	Proxy string `json:"proxy"`
	RequestedAt string `json:"requested"`
}

type Client interface {
	Login(request LoginRequest) (*LoginResponse, error)
	ListDevices(request ListDevicesRequest) (*ListDevicesResponse, error)
	Connect(ConnectRequest) (*ConnectResponse, error)
}

func (rc *rclient) Connect(request ConnectRequest) (*ConnectResponse, error) {
	url := fmt.Sprintf("%s/device/connect", rc.config.GetString(urlKey))

	b, err := json.Marshal(&request)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Add(devKeyHeader, request.APIKey)
	req.Header.Add(tokenKeyHeader, request.Token)
	req.Header.Add("content-type", "application/json")

	if err != nil {
		return nil, err
	}

	resp := new(ConnectResponse)

	r, err := rc.do(req, resp)

	if err != nil {
		return resp, err
	}

	return r.(*ConnectResponse), nil
}

func (rc *rclient) ListDevices(request ListDevicesRequest) (*ListDevicesResponse, error) {

	url := fmt.Sprintf("%s/device/list/all", rc.config.GetString(urlKey))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add(devKeyHeader, request.APIKey)
	req.Header.Add(tokenKeyHeader, request.Token)
	req.Header.Add("content-type", "application/json")

	if err != nil {
		return nil, err
	}

	resp := new(ListDevicesResponse)

	r, err := rc.do(req, resp)

	if err != nil {
		return resp, err
	}

	return r.(*ListDevicesResponse), nil
}

func (rc *rclient) Login(request LoginRequest) (*LoginResponse, error) {

	url := fmt.Sprintf("%s/user/login", rc.config.GetString(urlKey))

	b, err := json.Marshal(&request)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	req.Header.Add(devKeyHeader, request.APIKey)
	req.Header.Add("content-type", "application/json")

	resp := new(LoginResponse)

	r, err := rc.do(req, resp)

	if err != nil {
		return resp, err
	}

	return r.(*LoginResponse), nil

	return resp, err
}

func (rc *rclient) do(req *http.Request, resp interface{}) (interface{}, error) {

	hresp, err := rc.c.Do(req)

	if err != nil {
		return nil, err
	}

	defer hresp.Body.Close()

	if hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status from server: %v %s", hresp.StatusCode, hresp.Status)
	}

	body, err := ioutil.ReadAll(hresp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &resp)

	return resp, err
}

func NewClient(config *viper.Viper, httpclient *http.Client) Client {
	rc := new(rclient)

	rc.config = config

	if httpclient == nil {
		rc.c = http.DefaultClient
	} else {
		rc.c = httpclient
	}

	return rc
}