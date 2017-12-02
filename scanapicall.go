package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zamirka/interfaxScanApi/methods"
	"github.com/zamirka/interfaxScanApi/utils"
)

var ctx utils.AppContext

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
	err := utils.InitExecutionContext(&ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	accessData, err := methods.Login(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	myq, err := getAllQueries(accessData.AccessToken, ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(myq)
}

func getSearchQueryByID(accessToken string, queryID int, context utils.AppContext) (query *SearchQuery, err error) {
	requestURL := context.APIURL + "userQuery/" + strconv.Itoa(queryID)

	req, err := utils.PrepareRequest("GET", requestURL, nil, accessToken)

	if err != nil {
		return nil, err
	}

	var qs SearchQuery
	err = utils.MakeRequest(req, &qs)

	if err != nil {
		return nil, err
	}

	return &qs, nil
}

func createOrUpdateSearchQuery(accessToken string, query SearchQuery, context utils.AppContext) (queryID int, err error) {
	requestURL := context.APIURL + "userQuery"
	dataToSend, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}
	req, err := utils.PrepareRequest("POST", requestURL, bytes.NewReader(dataToSend), accessToken)

	if err != nil {
		return 0, err
	}

	var objmap map[string]*json.RawMessage
	err = utils.MakeRequest(req, &objmap)

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

func getAllQueries(accessToken string, context utils.AppContext) (queries []SearchQuery, err error) {
	requestURL := context.APIURL + "userQuery"
	req, err := utils.PrepareRequest("GET", requestURL, nil, accessToken)

	var myQueries []SearchQuery
	err = utils.MakeRequest(req, &myQueries)

	if err != nil {
		return nil, err
	}

	return myQueries, nil
}
