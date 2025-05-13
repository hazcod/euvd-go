package main

import (
	"context"
	"flag"
	"github.com/hazcod/euvd-go/pkg/euvd"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	defaultLogLevel = logrus.InfoLevel
	dateLayout      = "2006-01-02"
)

func main() {
	ctx := context.Background()

	logger := logrus.New()
	logger.SetLevel(defaultLogLevel)

	logLevel := flag.String("log", defaultLogLevel.String(), "Log level (debug, info, warn, error)")

	assigner := flag.String("assigner", "", "Filter by assigner")
	vendor := flag.String("vendor", "", "Filter by vendor")
	product := flag.String("product", "", "Filter by product")
	text := flag.String("text", "", "Full-text search query")

	fromDateStr := flag.String("from-date", "", "Filter from date (YYYY-MM-DD)")
	toDateStr := flag.String("to-date", "", "Filter to date (YYYY-MM-DD)")

	fromScore := flag.Int("from-score", -1, "Filter minimum CVSS score (0-10)")
	toScore := flag.Int("to-score", -1, "Filter maximum CVSS score (0-10)")

	fromEpss := flag.Int("from-epss", -1, "Filter minimum EPSS (0-100)")
	toEpss := flag.Int("to-epss", -1, "Filter maximum EPSS (0-100)")

	exploited := flag.String("exploited", "", "Filter by exploited status (true/false)")

	page := flag.Int("page", 0, "Pagination page number")
	size := flag.Int("size", 100, "Number of results per page")

	flag.Parse()
	command := strings.ToLower(flag.Arg(0))

	logrusLevel, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logger.WithError(err).Error("invalid log level provided")
		logrusLevel = defaultLogLevel
	}
	logger.SetLevel(logrusLevel)

	//

	var fromDate, toDate time.Time
	if *fromDateStr != "" {
		parsed, err := time.Parse(dateLayout, *fromDateStr)
		if err != nil {
			logger.WithError(err).Fatal("invalid from-date format, expected YYYY-MM-DD")
		}
		fromDate = parsed
	}
	if *toDateStr != "" {
		parsed, err := time.Parse(dateLayout, *toDateStr)
		if err != nil {
			logger.WithError(err).Fatal("invalid to-date format, expected YYYY-MM-DD")
		}
		toDate = parsed
	}

	// Handle optional int pointers
	var fromScorePtr, toScorePtr, fromEpssPtr, toEpssPtr *int
	if *fromScore >= 0 {
		fromScorePtr = fromScore
	}
	if *toScore >= 0 {
		toScorePtr = toScore
	}
	if *fromEpss >= 0 {
		fromEpssPtr = fromEpss
	}
	if *toEpss >= 0 {
		toEpssPtr = toEpss
	}

	var exploitedPtr *bool
	if *exploited != "" {
		val, err := strconv.ParseBool(*exploited)
		if err != nil {
			logger.WithError(err).Fatal("invalid exploited value, must be true or false")
		}
		exploitedPtr = &val
	}

	// Build search options
	opts := euvd.SearchOpts{
		Assigner:  *assigner,
		Vendor:    *vendor,
		Product:   *product,
		Text:      *text,
		FromDate:  fromDate,
		ToDate:    toDate,
		FromScore: fromScorePtr,
		ToScore:   toScorePtr,
		FromEpss:  fromEpssPtr,
		ToEpss:    toEpssPtr,
		Exploited: exploitedPtr,
		Page:      *page,
		Size:      *size,
	}

	// --

	euvdClient := euvd.New(logger)

	switch command {
	case "search":
		resp, err := euvdClient.Search(ctx, opts)
		if err != nil {
			logger.WithError(err).Fatal("failed to search EUVD")
		}

		for _, item := range resp.Items {
			product := ""
			if len(item.EnisaIDProduct) > 0 {
				product = item.EnisaIDProduct[0].Product.Name
			}

			vendor := ""
			if len(item.EnisaIDVendor) > 0 {
				vendor = item.EnisaIDVendor[0].Vendor.Name
			}

			logger.WithFields(logrus.Fields{
				"id":      item.ID,
				"product": product,
				"vendor":  vendor,
				"score":   item.BaseScore,
				"ref":     item.References,
			}).Println()
		}
		return
	case "lookup":
		if len(flag.Args()) < 2 {
			logger.Debug(flag.Args())
			logger.Fatal("missing EUVD ID.. use euvd lookup <id>")
		}

		resp, err := euvdClient.Lookup(ctx, flag.Arg(1))
		if err != nil {
			logger.WithError(err).Fatal("failed to search EUVD")
		}

		product := ""
		if len(resp.EnisaIDProduct) > 0 {
			product = resp.EnisaIDProduct[0].Product.Name
		}

		vendor := ""
		if len(resp.EnisaIDVendor) > 0 {
			vendor = resp.EnisaIDVendor[0].Vendor.Name
		}

		logger.WithFields(logrus.Fields{
			"id":      resp.ID,
			"product": product,
			"vendor":  vendor,
			"score":   resp.BaseScore,
			"ref":     resp.References,
		}).Println()
		return
	default:
		logger.Fatalf("unknown command '%s'. Use search,lookup", flag.Arg(0))
	}

}
