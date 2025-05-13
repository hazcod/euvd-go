package euvd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type LookupResponse struct {
	Items []Item `json:"items"`
	Total int    `json:"total"`
}

func (e *EUVD) Lookup(ctx context.Context, euvdID string) (*Item, error) {
	params := url.Values{}
	params.Add("id", euvdID)

	fullURL := fmt.Sprintf("%s/enisaid?%s", baseURL, params.Encode())

	timeoutCtx, _ := context.WithTimeout(ctx, getTimeout)

	req, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create API request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if e.logger.IsLevelEnabled(logrus.DebugLevel) {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			e.logger.WithError(err).Warn("Failed to dump HTTP request")
		} else {
			e.logger.Debug(string(dump))
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute API request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		response, _ := io.ReadAll(resp.Body)
		e.logger.WithField("status_code", resp.StatusCode).Debug(string(response))
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	responseByes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %w", err)
	}

	if e.logger.IsLevelEnabled(logrus.DebugLevel) {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			e.logger.WithError(err).Warn("Failed to dump HTTP response")
		} else {
			e.logger.Debug(string(dump))
		}
	}

	var result Item
	if err := json.NewDecoder(bytes.NewReader(responseByes)).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return &result, nil
}
