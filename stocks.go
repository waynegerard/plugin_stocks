package stocks

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "strings"
)

type Quote struct {
    Symbol string `json:"symbol"`
    EPS float64 `json:"EarningsShare,string"`
    EPSEstCY float64 `json:"EPSEstimateCurrentYear,string"`
    EPSEstNY float64 `json:"EPSEstimateNextYear,string"`
    EPSEstNQ float64 `json:"EPSEstimateNextQuarter,string"`
    PreviousClose float64 `json:"PreviousClose,string"`
    PE float64
}

type QuotesWrapper struct {
    Quotes []Quote `json:"quote"`
}

type QuoteResponse struct {
    Count int `json:"count"`
    Created string `json:"created"`
    Lang string `json:"lang"`
    Results QuotesWrapper `json:"results"`
}

type QueryWrapper struct {
    Response QuoteResponse `json:"query"`
}

func urlWithQuery(symbols string) (string) {
    return fmt.Sprintf("%s%s%s", "https://query.yahooapis.com/v1/public/yql?q=select%20Symbol%2CEarningsShare%2CEPSEstimateCurrentYear%2CEPSEstimateNextYear%2CEPSEstimateNextQuarter%2CPreviousClose%20from%20yahoo.finance.quotes%20where%20symbol%20in%20(", symbols, ")&format=json&env=http%3A%2F%2Fdatatables.org%2Falltables.env&callback=")
}

func createSymbolQuery(symbols []string) (string) {
    return fmt.Sprintf("%s%s%s","%22", strings.Join(symbols, "%22%2C%22"), "%22")
}

func getRequestData(url string) ([]Quote) {
    req, err := http.Get(url)

    if err != nil {
        fmt.Printf("%s\n", err)
        panic(err)
    }
    defer req.Body.Close()

    contents, err := ioutil.ReadAll(req.Body)
    if err != nil {
        fmt.Printf("%s", err)
        panic(err)
    }

    // Decode
    var queryResponse QueryWrapper
    err = json.NewDecoder(strings.NewReader(string(contents))).Decode(&queryResponse)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    return queryResponse.Response.Results.Quotes
}

func getDataForSymbols(symbols []string) []Quote {
    // Insert 25 symbols into URL
    symbolsString := createSymbolQuery(symbols)
    url := urlWithQuery(symbolsString)

    // Fetch data for said symbols
    data := getRequestData(url)

    return data
}
