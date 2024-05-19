package poynt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MikeB1124/trimana-dashboard-api/configuration"
)

func GetPoyntTransactions(startAt string, endAt string) ([]Transaction, error) {
	// Get the Poynt configuration
	config := configuration.GetConfig()
	poyntAccessToken, err := configuration.GetPoyntJWTAccessToken()
	if err != nil {
		return nil, fmt.Errorf("error getting Poynt access token: %v", err)
	}
	// Construct the URL for the API request
	endpoint_url := fmt.Sprintf("/businesses/%s/transactions?limit=100&startAt=%s&endAt=%s&orderBy=ASC", config.Poynt.BusinessID, startAt, endAt)
	poyntTransactionsResponse, err := fetchTransactions(config.Poynt.URL, endpoint_url, poyntAccessToken)
	if err != nil {
		return nil, fmt.Errorf("error fetching transactions: %v", err)
	}

	return poyntTransactionsResponse.Transactions, nil
}

func fetchTransactions(base_url string, endpoint_url string, accessToken string) (*PoyntTransactionsResponse, error) {
	req, err := http.NewRequest("GET", base_url+endpoint_url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Add("api-version", "1.2")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	var poyntTransactionsResponse PoyntTransactionsResponse
	err = json.Unmarshal(body, &poyntTransactionsResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return nil, err
	}
	return &poyntTransactionsResponse, nil
}
