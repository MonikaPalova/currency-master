package coinapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/MonikaPalova/currency-master/httputils"
)

const (
	baseUrl      = "https://rest.coinapi.io"
	apiKeyHeader = "X-CoinAPI-Key"
	apiKey       = "D8096E91-86D8-4998-B5B8-C785CE5D58AD"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	var c Client
	c.httpClient = &http.Client{}

	return &c
}

func (c Client) GetAssets() ([]Asset, *httputils.CoinApiError) {
	request, err := setUpRequest(http.MethodGet, "/v1/assets", nil)
	if err != nil {
		return nil, &httputils.CoinApiError{Err: err, Message: "could not set up get assets request", StatusCode: http.StatusInternalServerError}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, &httputils.CoinApiError{Err: err, Message: "could not execute get assets request", StatusCode: http.StatusInternalServerError}
	}
	defer response.Body.Close()

	if err := validateResponseCode(response); err != nil {
		return nil, err
	}

	// TODO should I use pointers here or the garbage collector handles the copies?
	var assets []Asset
	if err = json.NewDecoder(response.Body).Decode(&assets); err != nil {
		return nil, &httputils.CoinApiError{Err: err, Message: "could not parse get assets JSON response", StatusCode: http.StatusInternalServerError}
	}

	assets = removeInvalidAssets(assets)
	return assets, nil
}

func removeInvalidAssets(assets []Asset) []Asset {
	var filtered []Asset
	for _, asset := range assets {
		// if asset.PriceUSD > 0 {
		filtered = append(filtered, asset)
		// }
	}

	return filtered
}

func (c Client) GetAssetById(id string) (*Asset, *httputils.CoinApiError) {
	request, err := setUpRequest(http.MethodGet, "/v1/assets/"+id, nil)
	if err != nil {
		return nil, &httputils.CoinApiError{Err: err, Message: "could not set up request get asset by id with assetId [" + id + "]", StatusCode: http.StatusInternalServerError}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, &httputils.CoinApiError{Err: err, Message: "could not execute request get asset by id [" + id + "]", StatusCode: http.StatusInternalServerError}
	}
	defer response.Body.Close()

	if err := validateResponseCode(response); err != nil {
		return nil, err
	}

	var assets []Asset
	if err = json.NewDecoder(response.Body).Decode(&assets); err != nil {
		return nil, &httputils.CoinApiError{Err: err, Message: "could not parse get asset by id JSON response", StatusCode: http.StatusInternalServerError}
	}

	if len(assets) == 0 {
		return nil, &httputils.CoinApiError{Err: nil, Message: "Asset with id [" + id + "] doesn't exist", StatusCode: http.StatusNotFound}
	}
	return &assets[0], nil
}

func validateResponseCode(response *http.Response) *httputils.CoinApiError {
	if response.StatusCode != http.StatusOK {
		var responseBody struct {
			Error string `json:"error"`
		}

		var err error
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return &httputils.CoinApiError{Err: err, Message: "could not parse JSON response with error", StatusCode: http.StatusInternalServerError}
		}

		return &httputils.CoinApiError{Err: err, Message: "coin API returned code [" + strconv.Itoa(response.StatusCode) + "] with messsage [" + responseBody.Error + "]", StatusCode: response.StatusCode}
	}
	return nil
}

func setUpRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	url := baseUrl + endpoint
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add(apiKeyHeader, apiKey)
	request.Header.Add("Accept", "application/json")
	return request, nil
}
