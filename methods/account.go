package methods

import (
	"bytes"

	"github.com/shopspring/decimal"
	"github.com/zamirka/interfaxScanApi/utils"
)

// LoginResponse is a structure for login data
type LoginResponse struct {
	AccessToken string `json:"AccessToken"`
	Expire      string `json:"Expire"`
}

// BalanceResponse is a response type for balance
type BalanceResponse struct {
	Unlimited             bool            `json:"unlimited"`
	Balance               decimal.Decimal `json:"balance"`
	SearchTermCost        decimal.Decimal `json:"searchTgo ermCost"`
	SearchRateBlockPeriod int             `json:"searchRateBlockPeriod"`
}

//Login is a POST method that authenticates a user in API and returns a token for future calls
func Login(context utils.AppContext) (token *LoginResponse, err error) {
	requestURL := context.APIURL + "account/login"
	var payload bytes.Buffer
	payload.WriteString(`{"login":"`)
	payload.WriteString(context.Username)
	payload.WriteString(`","password":"`)
	payload.WriteString(context.Password)
	payload.WriteString(`"}`)

	req, err := utils.PrepareRequest("POST", requestURL, bytes.NewReader(payload.Bytes()), "")

	if err != nil {
		return nil, err
	}

	var t LoginResponse
	err = utils.MakeRequest(req, &t)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Balance is a GET method that returns Balance structure
func Balance(accessToken string, context *utils.AppContext) (balance *BalanceResponse, err error) {
	requestURL := context.APIURL + "account/balance"

	req, err := utils.PrepareRequest("GET", requestURL, nil, accessToken)

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
