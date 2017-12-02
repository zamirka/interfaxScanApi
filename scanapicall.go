package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	d "github.com/shopspring/decimal"
)

var context AppContext

// AppContext is a configuration which is read from filw conf.json
type AppContext struct {
	APIURL   string `json:"APIURL"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func initExecutionContext(context *AppContext) error {
	file, err := os.Open("conf.json")
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(context)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// LoginResponse is a structure for login data
type LoginResponse struct {
	AccessToken string `json:"AccessToken"`
	Expire      string `json:"Expire"`
}

// BalanceResponse is a response type for balance
type BalanceResponse struct {
	Unlimited             bool      `json:"unlimited"`
	Balance               d.Decimal `json:"balance"`
	SearchTermCost        d.Decimal `json:"searchTgo ermCost"`
	SearchRateBlockPeriod int       `json:"searchRateBlockPeriod"`
}

// SourceRegion region criterias for search query
type SourceRegion struct {
	RegionID int    `json:"regionId"`
	Type     string `json:"type"`
}

// SearchArea area criterias for search query
type SearchArea struct {
	Sources      []int `json:"sources"`
	SourceGroups []int `json:"sourcesGroups"`
	// SourcesAggregation valid values:
	// includeOnly (elements of contentSources/contentSourcesGroups).Default value
	// except(elements of contentSources/contentSourcesGroups)
	SourcesAggregation                    string         `json:"sourcesAggregation"`
	Topics                                []int          `json:"topics"`
	Regions                               []SourceRegion `json:"regions"`
	Categories                            []int          `json:"categories"`
	Levels                                []int          `json:"levels"`
	ExcludeDocumentTypes                  []string       `json:"excludeDocumentTypes"`
	Languages                             []string       `json:"languages"`
	ExcludeDocumentsWithoutFullTextAccess bool           `json:"excludeDocumentsWithoutFullTextAccess"`
}

// SearchTerm term criterias for search query
type SearchTerm struct {
	Operator   string       `json:"operator"`
	Type       string       `json:"type"`
	ChildTerms []SearchTerm `json:"childTerms"`
	// WordsDistance valid values:
	// undefined – все термы рядом.
	// withinWord – в пределах одного слова.
	// within2Words – в пределах двух слов.
	// within3Words – в пределах трёх слов.
	// within4Words – в пределах четырёх слов.
	// within5Words – в пределах пяти слов.
	// withinSentence – в пределах предложения.
	// withinParagraph – в пределах абзаца.
	WordsDistance     string `json:"wordsDistance"`
	TextToSearch      string `json:"textToSearch"`
	ExactlyPhrase     bool   `json:"exactlyPhrase"`
	ConsiderWordOrder bool   `json:"considerWordOrder"`
	// SearchDocumentComponent valid values
	// any - везде Default value
	// title – только по заголовку
	// content – только по тексту
	SearchDocumentComponent string `json:"searchDocumentComponent"`
	EntityID                int    `json:"entityId"`
	WithStaff               bool   `json:"withStaff"`
	FindQuotes              bool   `json:"findQuotes"`
	DirectSpeechText        string `json:"directSpeechText"`
	// MentionContext valid values
	// any – контекст не задан. Default Value
	// positiveTone – позитивная тональность.
	// negativeTone – негативная тональность.
	// main – главная роль.
	// intentions – намерения.
	MentionContext string `json:"mentionContext"`
	// MentionContext valid values
	// maxFullness – максимальная полнота.
	// maxPrecision – максимальная точность
	TonePrecision string `json:"tonePrecision"`
	Subjects      []int  `json:"subjects"`
}

// QueryParams are all possible parameters
type QueryParams struct {
	SearchArea SearchArea   `json:"searchArea"`
	Terms      []SearchTerm `json:"terms"`
}

// SearchQuery is a query with all possible parameters which can be saved
type SearchQuery struct {
	UserQueryID int         `json:"userQueryId"`
	Name        string      `json:"name"`
	Query       QueryParams `json:"query"`
}

func main() {
	err := initExecutionContext(&context)
	if err != nil {
		fmt.Println(err)
		return
	}

	accessData, err := login(context.Username, context.Password)
	if err != nil {
		fmt.Println(err)
		return
	}

	myq, err := getAllQueries(accessData.AccessToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(myq)
}

func prepareRequest(httpMethod string, url string, payload io.Reader, accessToken string) (request *http.Request, err error) {
	req, err := http.NewRequest(httpMethod, url, payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	if accessToken != "" {
		var buffer bytes.Buffer
		buffer.WriteString("Bearer ")
		buffer.WriteString(accessToken)
		req.Header.Add("Authorization", buffer.String())
	}
	return req, nil
}

func makeRequest(request *http.Request, dataObject interface{}) error {
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// bodyStr := string(body)
	// fmt.Println(bodyStr)
	err = json.Unmarshal(body, &dataObject)

	if err != nil {
		return err
	}
	return nil
}

func login(userLogin string, userPassword string) (token *LoginResponse, err error) {
	requestURL := context.APIURL + "account/login"
	var payload bytes.Buffer
	payload.WriteString(`{"login":"`)
	payload.WriteString(userLogin)
	payload.WriteString(`","password":"`)
	payload.WriteString(userPassword)
	payload.WriteString(`"}`)

	req, err := prepareRequest("POST", requestURL, bytes.NewReader(payload.Bytes()), "")

	if err != nil {
		return nil, err
	}

	var t LoginResponse
	err = makeRequest(req, &t)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func balance(accessToken string) (balance *BalanceResponse, err error) {
	requestURL := context.APIURL + "account/balance"

	req, err := prepareRequest("GET", requestURL, nil, accessToken)

	if err != nil {
		return nil, err
	}

	var b BalanceResponse
	err = makeRequest(req, &b)

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func getSearchQueryByID(accessToken string, queryID int) (query *SearchQuery, err error) {
	requestURL := context.APIURL + "userQuery/" + strconv.Itoa(queryID)

	req, err := prepareRequest("GET", requestURL, nil, accessToken)

	if err != nil {
		return nil, err
	}

	var qs SearchQuery
	err = makeRequest(req, &qs)

	if err != nil {
		return nil, err
	}

	return &qs, nil
}

func createOrUpdateSearchQuery(accessToken string, query SearchQuery) (queryID int, err error) {
	requestURL := context.APIURL + "userQuery"
	dataToSend, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}
	req, err := prepareRequest("POST", requestURL, bytes.NewReader(dataToSend), accessToken)

	if err != nil {
		return 0, err
	}

	var objmap map[string]*json.RawMessage
	err = makeRequest(req, &objmap)

	if err != nil {
		return 0, err
	}

	var qid int
	err = json.Unmarshal(*objmap["id"], &qid)
	if err != nil {
		return 0, err
	}
	return qid, nil
}

func getAllQueries(accessToken string) (queries []SearchQuery, err error) {
	requestURL := context.APIURL + "userQuery"
	req, err := prepareRequest("GET", requestURL, nil, accessToken)

	var myQueries []SearchQuery
	err = makeRequest(req, &myQueries)

	if err != nil {
		return nil, err
	}

	return myQueries, nil
}
