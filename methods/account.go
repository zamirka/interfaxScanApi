package methods

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/shopspring/decimal"
	"github.com/zamirka/interfaxScanApi/utils"
)

// BalanceResponse is a response type for balance
type BalanceResponse struct {
	Unlimited             bool            `json:"unlimited"`
	Balance               decimal.Decimal `json:"balance"`
	SearchTermCost        decimal.Decimal `json:"searchTgo ermCost"`
	SearchRateBlockPeriod int             `json:"searchRateBlockPeriod"`
}

//Login is a POST method that authenticates a user in API and writes a token and it's expire date for future calls into provided context
func Login(context *utils.AppContext) (err error) {
	requestURL := context.APIURL + "account/login"
	var payload bytes.Buffer
	payload.WriteString(`{"login":"`)
	payload.WriteString(context.Username)
	payload.WriteString(`","password":"`)
	payload.WriteString(context.Password)
	payload.WriteString(`"}`)
	var req *http.Request
	if req, err = utils.PrepareRequest("POST", requestURL, bytes.NewReader(payload.Bytes()), ""); err != nil {
		return err
	}
	var configData map[string]*json.RawMessage
	if err = utils.MakeRequest(req, &configData); err != nil {
		return err
	}
	var token string
	if err = json.Unmarshal(*configData["AccessToken"], &token); err != nil {
		return err
	}
	context.AccessToken = token
	var expireDate string
	if err = json.Unmarshal(*configData["Expire"], &expireDate); err != nil {
		return err
	}
	context.Expire = expireDate
	return nil
}

// Balance is a GET method that returns Balance structure
func Balance(context *utils.AppContext) (balance *BalanceResponse, err error) {
	requestURL := context.APIURL + "account/balance"
	var req *http.Request
	if req, err = utils.PrepareRequest("GET", requestURL, nil, context.AccessToken); err != nil {
		return nil, err
	}
	var b BalanceResponse
	if err = utils.MakeRequest(req, &b); err != nil {
		return nil, err
	}
	return &b, nil
}
