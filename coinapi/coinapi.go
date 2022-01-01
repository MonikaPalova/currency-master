package coinapi

import (
	"currency-master/model/wallet"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseUrl      = "https://rest.coinapi.io"
	apiKeyHeader = "X-CoinAPI-Key"
	apiKey       = "D8096E91-86D8-4998-B5B8-C785CE5D58AD"
)

type Client struct {
	HttpClient *http.Client
}

func (c Client) GetAssets() ([]wallet.Asset, error) {
	request, err := setUpRequest(http.MethodGet, "/v1/assets/DUB,DOGE", nil)
	if err != nil {
		return nil, fmt.Errorf("could not set up get assets request: %s", err.Error())
	}

	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not execute get assets request: %s", err.Error())
	}
	defer response.Body.Close()

	if err := validateResponseCode(response); err != nil {
		return nil, fmt.Errorf("response of get assets is not valid: %s", err.Error())
	}

	var assets []wallet.Asset
	if err = json.NewDecoder(response.Body).Decode(&assets); err != nil {
		return nil, fmt.Errorf("could not parse get assets JSON response: %s", err.Error())
	}

	return assets, nil
}

func validateResponseCode(response *http.Response) (err error) {
	if response.StatusCode != http.StatusOK {
		var responseBody struct {
			Error string `json:"error"`
		}
		if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return fmt.Errorf("could not parse JSON response with error: %s", err.Error())
		}

		err = fmt.Errorf("coin API returned code [%d] with messsage [%s]", response.StatusCode, responseBody.Error)
	}
	return err
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
