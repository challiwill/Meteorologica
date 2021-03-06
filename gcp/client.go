package gcp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/calendar"
	"github.com/challiwill/meteorologica/csv"
	"github.com/challiwill/meteorologica/datamodels"
	"github.com/challiwill/meteorologica/errare"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

var IAAS = "GCP"

//go:generate counterfeiter . StorageService

type StorageService interface {
	DailyUsage(string, string) (*http.Response, error)
	Insert(string, *storage.Object, *os.File) (*storage.Object, error)
}

type DetailedUsageReport [][]byte

type Client struct {
	StorageService StorageService
	BucketName     string
	Log            *logrus.Logger
	Location       *time.Location
}

func NewClient(log *logrus.Logger, location *time.Location, jsonCredentials []byte, bucketName string) (*Client, error) {
	jwtConfig, err := google.JWTConfigFromJSON(jsonCredentials, "https://www.googleapis.com/auth/devstorage.read_write")
	if err != nil {
		return nil, err
	}
	service, err := storage.New(jwtConfig.Client(oauth2.NoContext))
	if err != nil {
		return nil, err
	}
	return &Client{
		StorageService: &storageService{service: service},
		BucketName:     bucketName,
		Log:            log,
		Location:       location,
	}, nil
}

func (c Client) Name() string {
	return IAAS
}

func (c Client) GetNormalizedUsage() (datamodels.Reports, error) {
	c.Log.Info("Getting monthly GCP usage...")
	c.Log.Debug("Entering gcp.GetNormalizedUsage")
	defer c.Log.Debug("Returning gcp.GetNormalizedUsage")

	gcpMonthlyUsage, err := c.GetBillingData()
	if err != nil {
		c.Log.Error("Failed to get GCP monthly usage")
		return datamodels.Reports{}, err
	}
	c.Log.Debug("Got monthly GCP usage")

	monthlyReport := []*Usage{}
	for i, usage := range gcpMonthlyUsage {
		var readerCleaner *csv.ReaderCleaner
		readerCleaner, err = csv.NewReaderCleaner(bytes.NewReader(usage), 18, 14) // ambiguously 18 and 14...
		if err != nil {
			return datamodels.Reports{}, err
		}

		dailyReport := []*Usage{}
		err = csv.GenerateReports(readerCleaner, &dailyReport)
		if err != nil {
			c.Log.Errorf("Failed to parse GCP usage for day: %d %s: %s", i+1, time.Now().In(c.Location).Month().String(), err.Error())
			continue
		}
		dailyReport = c.setDate(dailyReport, i+1)
		monthlyReport = append(monthlyReport, dailyReport...)
	}
	if len(monthlyReport) == 0 {
		return datamodels.Reports{}, csv.NewEmptyReportError("parsing GCP usage")
	}

	normalizer := NewNormalizer(c.Log, c.Location)
	normalizedReports := normalizer.Normalize(monthlyReport)
	normalizedReports = datamodels.ConsolidateReports(normalizedReports)
	return normalizedReports, nil
}

func (c Client) GetBillingData() (DetailedUsageReport, error) {
	c.Log.Debug("Entering gcp.GetBillingData")
	defer c.Log.Debug("Returning gcp.GetBillingData")

	monthlyUsageReport := DetailedUsageReport{}
	year, month, day := calendar.YesterdaysDate(c.Location)
	for i := 1; i < day; i++ {
		dailyUsage, err := c.DailyUsageReport(year, month, i)
		if err != nil {
			c.Log.Warnf("Failed to get GCP Daily Usage for %s, %d: %s", time.Now().In(c.Location).Month().String(), i, err.Error())
			continue
		}
		monthlyUsageReport = append(monthlyUsageReport, dailyUsage)
	}
	return monthlyUsageReport, nil
}

func (c Client) DailyUsageReport(year int, month time.Month, day int) ([]byte, error) {
	c.Log.Debug("Entering gcp.DailyUsageReport")
	defer c.Log.Debug("Returning gcp.DailyUsageReport")

	resp, err := c.StorageService.DailyUsage(c.BucketName, c.dailyBillingFileName(year, month, day))
	if err != nil {
		return nil, errare.NewRequestError(err, IAAS)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errare.NewResponseError(resp.Status, IAAS)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c Client) dailyBillingFileName(year int, month time.Month, day int) string {
	monthStr := calendar.PadMonth(month)
	dayStr := padDay(day)
	return url.QueryEscape(strings.Join([]string{"Billing", strconv.Itoa(year), monthStr, dayStr}, "-") + ".csv")
}

func padDay(day int) string {
	d := strconv.Itoa(day)
	if day < 10 {
		return "0" + d
	}
	return d
}

func (c Client) setDate(usages []*Usage, day int) []*Usage {
	year, month, _ := calendar.YesterdaysDate(c.Location)
	for i, _ := range usages {
		usages[i].TimeFetched = time.Date(year, month, day, 0, 0, 0, 0, c.Location)
	}
	return usages
}
