package main

import (
	"errors"
	"fmt"
	"github.com/dghubble/sling"
	"log"
	"net/http"
)

type AnalyzerApi struct {
	sling *sling.Sling
}

func NewAnalyzerApi(baseUrl string, client *http.Client) *AnalyzerApi {
	return &AnalyzerApi{
		sling: sling.New().Client(client).Base(baseUrl),
	}
}

func (a *AnalyzerApi) AnalyzeImage(url string) (*AnalyzeResponse, error) {
	req := &AnalyzeRequest{Url: url}
	scoreResponse := &AnalyzeResponse{}

	res, err := a.sling.New().Get("/labels").QueryStruct(req).ReceiveSuccess(scoreResponse)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		log.Print("Error: Status Code: ", res.StatusCode)
		return nil, errors.New(fmt.Sprintf("analyzer-service returned status code: %d ", res.StatusCode))
	}

	return scoreResponse, nil
}

type AnalyzeRequest struct {
	Url string `url:"url"`
}

type AnalyzeResponse struct {
	Labels []*LabelResult `json:"labels"`
}

type LabelResult struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}
