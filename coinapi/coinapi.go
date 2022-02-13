package coinapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MonikaPalova/currency-master/config"
)

type client struct {
	httpClient *http.Client
	config     *config.CoinAPI
}

func newClient() *client {
	return &client{&http.Client{}, config.NewCoinAPI()}
}

func (c client) getAssets() ([]Asset, error) {
	request, err := c.setUpRequest(http.MethodGet, c.config.AssetsUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("could not set up get assets request, %v", err.Error())
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not execute get assets request, %v", err.Error())
	}
	defer response.Body.Close()

	if err := validateResponseCode(response); err != nil {
		return nil, err
	}

	var assets []Asset
	if err = json.NewDecoder(response.Body).Decode(&assets); err != nil {
		return nil, fmt.Errorf("could not parse get assets JSON response, %v", err.Error())
	}

	assets = removeInvalidAssets(assets)
	return assets, nil
}

func removeInvalidAssets(assets []Asset) []Asset {
	var filtered []Asset
	for _, asset := range assets {
		if asset.PriceUSD > 0 {
			filtered = append(filtered, asset)
		}
	}

	return filtered
}

func validateResponseCode(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		var responseBody struct {
			Error string `json:"error"`
		}

		if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			return fmt.Errorf("could not parse JSON response with error, %v", err.Error())
		}

		return fmt.Errorf("coin API returned code %d with messsage %s", response.StatusCode, responseBody.Error)
	}
	return nil
}

func (c client) setUpRequest(method, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add(c.config.ApiKeyHeader, c.config.ApiKey)
	request.Header.Add("Accept", "application/json")
	return request, nil
}
