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
	"strconv"
	"time"
)

var (
	getTimeout = 10 * time.Second
)

type SearchOpts struct {
	Assigner  string
	Vendor    string
	Product   string
	Text      string
	FromDate  time.Time
	ToDate    time.Time
	FromScore *int
	ToScore   *int
	FromEpss  *int
	ToEpss    *int
	Exploited *bool
	Page      int
	Size      int
}

type SearchResponse struct {
	Items []Item `json:"items"`
	Total int    `json:"total"`
}
type Product struct {
	Name string `json:"name"`
}
type EnisaIDProduct struct {
	ID             string  `json:"id"`
	Product        Product `json:"product"`
	ProductVersion string  `json:"product_version,omitempty"`
}
type Vendor struct {
	Name string `json:"name"`
}
type EnisaIDVendor struct {
	ID     string `json:"id"`
	Vendor Vendor `json:"vendor"`
}
type Item struct {
	ID               string           `json:"id"`
	Description      string           `json:"description"`
	DatePublished    string           `json:"datePublished"`
	DateUpdated      string           `json:"dateUpdated"`
	BaseScore        float64          `json:"baseScore"`
	BaseScoreVersion string           `json:"baseScoreVersion"`
	BaseScoreVector  string           `json:"baseScoreVector"`
	References       string           `json:"references"`
	Aliases          string           `json:"aliases"`
	Assigner         string           `json:"assigner"`
	Epss             float64          `json:"epss"`
	EnisaIDProduct   []EnisaIDProduct `json:"enisaIdProduct"`
	EnisaIDVendor    []EnisaIDVendor  `json:"enisaIdVendor"`
}

func (e *EUVD) Search(ctx context.Context, opts SearchOpts) (*SearchResponse, error) {
	params := url.Values{}
	if opts.Assigner != "" {
		params.Add("assigner", opts.Assigner)
	}
	if opts.Vendor != "" {
		params.Add("vendor", opts.Vendor)
	}
	if opts.Product != "" {
		params.Add("product", opts.Product)
	}
	if opts.Text != "" {
		params.Add("text", opts.Text)
	}
	if !opts.FromDate.IsZero() {
		params.Add("fromDate", opts.FromDate.Format("2006-01-02"))
	}
	if !opts.ToDate.IsZero() {
		params.Add("toDate", opts.ToDate.Format("2006-01-02"))
	}
	if opts.FromScore != nil {
		params.Add("fromScore", strconv.Itoa(*opts.FromScore))
	}
	if opts.ToScore != nil {
		params.Add("toScore", strconv.Itoa(*opts.ToScore))
	}
	if opts.FromEpss != nil {
		params.Add("fromEpss", strconv.Itoa(*opts.FromEpss))
	}
	if opts.ToEpss != nil {
		params.Add("toEpss", strconv.Itoa(*opts.ToEpss))
	}
	if opts.Exploited != nil {
		params.Add("exploited", strconv.FormatBool(*opts.Exploited))
	}
	if opts.Page > 0 {
		params.Add("page", strconv.Itoa(opts.Page))
	}
	if opts.Size > 0 {
		params.Add("size", strconv.Itoa(opts.Size))
	}

	fullURL := fmt.Sprintf("%s/vulnerabilities?%s", baseURL, params.Encode())

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

	var result SearchResponse
	if err := json.NewDecoder(bytes.NewReader(responseByes)).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return &result, nil
}
