package methods

import (
	"bytes"
	"encoding/json"

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

	req, err := utils.PrepareRequest("POST", requestURL, bytes.NewReader(payload.Bytes()), "")

	if err != nil {
		return err
	}

	var configData map[string]*json.RawMessage
	err = utils.MakeRequest(req, &configData)

	if err != nil {
		return err
	}

	var token string
	err = json.Unmarshal(*configData["AccessToken"], &token)
	if err != nil {
		return err
	}

	context.AccessToken = token

	var expireDate string
	err = json.Unmarshal(*configData["Expire"], &expireDate)
	if err != nil {
		return err
	}

	context.Expire = expireDate

	return nil
}

// Balance is a GET method that returns Balance structure
func Balance(context *utils.AppContext) (balance *BalanceResponse, err error) {
	requestURL := context.APIURL + "account/balance"

	req, err := utils.PrepareRequest("GET", requestURL, nil, context.AccessToken)

	if err != nil {
		return nil, err
	}

	var b BalanceResponse
	err = utils.MakeRequest(req, &b)

	if err != nil {
		return nil, err
	}

	return &b, nil
}
