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
	"github.com/challiwill/meteorologica/csv"
	"github.com/challiwill/meteorologica/datamodels"
	"github.com/challiwill/meteorologica/errare"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

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
	return "GCP"
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
			c.Log.Errorf("Failed to parse GCP usage for day: %d %s: %s", i+1, datamodels.MONTH.String(), err.Error())
			continue
		}
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
	for i := 1; i < 31; i++ {
		dailyUsage, err := c.DailyUsageReport(i)
		if err != nil {
			c.Log.Warnf("Failed to get GCP Daily Usage for %s, %d: %s", datamodels.MONTH.String(), i, err.Error())
			continue
		}
		monthlyUsageReport = append(monthlyUsageReport, dailyUsage)
	}
	return monthlyUsageReport, nil
}

func (c Client) DailyUsageReport(day int) ([]byte, error) {
	c.Log.Debug("Entering gcp.DailyUsageReport")
	defer c.Log.Debug("Returning gcp.DailyUsageReport")
	var resp *http.Response
	var err error
	if datamodels.MONTH == time.Month(7) || datamodels.MONTH == time.Month(8) {
		resp, err = c.StorageService.DailyUsage(c.BucketName, datamodels.MONTH.String()+"/"+c.dailyBillingFileName(day))
	} else {
		resp, err = c.StorageService.DailyUsage(c.BucketName, c.dailyBillingFileName(day))
	}
	if err != nil {
		return nil, errare.NewRequestError(err, "GCP")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errare.NewResponseError(resp.Status, "GCP")
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c Client) PublishFileToBucket(name string) error {
	c.Log.Debug("Entering gcp.PublishFileToBucket")
	defer c.Log.Debug("Returning gcp.PublishFileToBucket")

	object := &storage.Object{
		Name:        name,
		ContentType: "text/csv",
	}
	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		c.Log.Errorf("Failed to open normalized file: %s", name)
		return err
	}

	res, err := c.StorageService.Insert(c.BucketName, object, file)
	if err != nil {
		c.Log.Errorf("Objects.Insert to bucket '%s' failed", c.BucketName)
		return err
	}
	c.Log.Infof("Created object %v at location %v", res.Name, res.SelfLink)

	return nil
}

func (c Client) dailyBillingFileName(day int) string {
	year, _, _ := time.Now().In(c.Location).Date()
	monthStr := padMonth(datamodels.MONTH)
	dayStr := padDay(day)
	return url.QueryEscape(strings.Join([]string{"Billing", strconv.Itoa(year), monthStr, dayStr}, "-") + ".csv")
}

func padMonth(month time.Month) string {
	m := strconv.Itoa(int(month))
	if month < 10 {
		return "0" + m
	}
	return m
}

func padDay(day int) string {
	d := strconv.Itoa(day)
	if day < 10 {
		return "0" + d
	}
	return d
}
