package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type credentials struct {
	APIKey string `yaml:"api-key"`
	APIId  string `yaml:"api-id"`
}

type ArquiveiClient struct {
	endpoint string
	page     Page
	cred     *credentials
	client   httpInterface
}

type httpInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewArquiveiClient(credentialsFilepath string) (*ArquiveiClient, error) {
	cred, err := getCredentials(credentialsFilepath)
	if err != nil {
		return nil, err
	}
	return &ArquiveiClient{
		endpoint: "https://sandbox-api.arquivei.com.br",
		client:   &http.Client{},
		page: Page{
			Next:     "",
			Previous: "",
		},
		cred: cred,
	}, nil
}

func getCredentials(filepath string) (*credentials, error) {
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	cred := &credentials{}
	if err := yaml.Unmarshal(f, &cred); err != nil {
		return nil, err
	}

	return cred, nil
}

func (cli *ArquiveiClient) RequestNFCs() ([]NFC, error) {
	var (
		err  error
		body []byte
		br   = &BodyResponse{}
	)

	if cli.page.Next == "" {
		body, err = cli.makeRequest("GET", fmt.Sprintf("%s/v1/nfe/received", cli.endpoint), []byte(""))
	} else {
		body, err = cli.makeRequest("GET", cli.page.Next, []byte(""))
	}
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, br)
	if err != nil {
		return nil, err
	}
	cli.page = br.Page
	return br.NFCs, nil
}

func (cli *ArquiveiClient) newAuthorizedRequest(method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-id", cli.cred.APIId)
	req.Header.Add("x-api-key", cli.cred.APIKey)
	req.Header.Set("Content-Type", "application/json")
	return req, nil

}

func (cli *ArquiveiClient) makeRequest(method, url string, body []byte) ([]byte, error) {
	var resp *http.Response
	var err error
	var attempts int

	maxAttempts := 5
	for attempts = 1; attempts < maxAttempts; attempts++ {
		req, err := cli.newAuthorizedRequest(method, url, body)
		if err != nil {
			return nil, fmt.Errorf("failed to create authorized request: %v", err)
		}

		resp, err = cli.client.Do(req)
		if err == nil && shouldRetry(resp) == false {
			if attempts > 1 {
				log.Info().
					Str("URL", req.URL.String()).
					Msgf("finished retry [%d/%d]", attempts, maxAttempts)
			}
			break
		}

		if err != nil {
			log.Error().Err(err).
				Str("URL", req.URL.String()).
				Msg("request error")
		}

		backoff := time.Duration(math.Pow(2, float64(attempts))) * time.Second
		if resp != nil {
			if resp.StatusCode == http.StatusForbidden {
				log.Warn().Msgf("received forbidden status code. Aborting...")
				return nil, fmt.Errorf("user forbidden")
			}

			if resp.StatusCode == http.StatusUnauthorized {
				log.Warn().Msgf("received unauthorized status code. Aborting...")
				return nil, fmt.Errorf("user unauthorized")
			}

			if resp.StatusCode == http.StatusNotFound {
				return nil, fmt.Errorf("endpoint not found")
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}
			resp.Body.Close()

			log.Error().
				Err(err).
				Str("status_code", fmt.Sprint(resp.StatusCode)).
				Str("backoff", fmt.Sprint(backoff)).
				Str("URL", req.URL.String()).
				Msg(string(data))
		}
		time.Sleep(backoff)
	}

	if attempts == maxAttempts {
		return nil, fmt.Errorf("maximum retry attempts reached")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	// fmt.Println(string(data))
	resp.Body.Close()

	return data, nil
}

func shouldRetry(resp *http.Response) bool {
	return resp == nil || (resp.StatusCode != 204 && resp.StatusCode != 200)
}
