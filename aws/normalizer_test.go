package aws_test

import (
	"time"

	"github.com/Sirupsen/logrus"
	. "github.com/challiwill/meteorologica/aws"
	"github.com/challiwill/meteorologica/datamodels"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Normalizer", func() {
	var (
		normalizer *Normalizer
		log        *logrus.Logger
		loc        *time.Location
	)

	BeforeEach(func() {
		log = logrus.New()
		log.Out = NewBuffer()
		loc = time.Now().Location()
		normalizer = NewNormalizer(log, loc, "my-region")
	})

	Describe("Normalize", func() {
		var (
			usageReports []*Usage
			reports      datamodels.Reports
		)

		BeforeEach(func() {
			usageReports = []*Usage{
				&Usage{
					InvoiceID:              "some-invoice-id",
					PayerAccountId:         "some-payer-account-id",
					LinkedAccountId:        "some-linked-account-id",
					RecordType:             "LinkedLineItem",
					RecordID:               "some-record-id",
					BillingPeriodStartDate: "some-start-date",
					BillingPeriodEndDate:   "some-end-date",
					InvoiceDate:            "some-invoice-date",
					PayerAccountName:       "some-payer-account-name",
					LinkedAccountName:      "some-linked-account-name",
					TaxationAddress:        "some-address",
					PayerPONumber:          "some-payer-number",
					ProductCode:            "some-product-code",
					ProductName:            "some-product-name",
					SellerOfRecord:         "some-seller-record",
					UsageType:              "some-usage-type",
					Operation:              "some-operation",
					RateId:                 "some-rate-id",
					ItemDescription:        "some-item-description",
					UsageStartDate:         "some-usage-start-date",
					UsageEndDate:           "some-usage-end-date",
					UsageQuantity:          0.51,
					BlendedRate:            "some-blended-rate",
					CurrencyCode:           "some-currency-code",
					CostBeforeTax:          "some-cost",
					Credits:                "some-credits",
					TaxAmount:              "some-tax-amount",
					TaxType:                "some-tax-type",
					TotalCost:              1.20,
				},
			}
		})

		JustBeforeEach(func() {
			reports = normalizer.Normalize(usageReports)
		})

		Context("with at least one report", func() {
			Context("with valid data", func() {
				It("returns the same number of reports", func() {
					Expect(reports).To(HaveLen((len(usageReports))))
				})

				It("returns properly converted reports", func() {
					Expect(reports[0]).To(Equal(datamodels.Report{
						ID:            usageReports[0].Hash("my-region"),
						AccountNumber: "some-linked-account-id",
						AccountName:   "some-linked-account-name",
						Day:           time.Now().Day() - 1,
						Month:         time.Now().Month(),
						Year:          time.Now().Year(),
						ServiceType:   "some-product-name",
						UsageQuantity: 0.51,
						Cost:          1.20,
						Region:        "my-region",
						UnitOfMeasure: "",
						Resource:      "AWS",
					}))
				})
			})

			Context("with rows that are not line items", func() {
				BeforeEach(func() {
					usageReports = append(usageReports,
						&Usage{
							InvoiceID:              "some-invoice-id",
							PayerAccountId:         "failure",
							LinkedAccountId:        "failure",
							RecordType:             "InvoiceTotal",
							RecordID:               "some-record-id",
							BillingPeriodStartDate: "some-start-date",
							BillingPeriodEndDate:   "some-end-date",
							InvoiceDate:            "some-invoice-date",
							PayerAccountName:       "some-payer-account-name",
							LinkedAccountName:      "some-linked-account-name",
							TaxationAddress:        "some-address",
							PayerPONumber:          "some-payer-number",
							ProductCode:            "some-product-code",
							ProductName:            "some-product-name",
							SellerOfRecord:         "some-seller-record",
							UsageType:              "some-usage-type",
							Operation:              "some-operation",
							RateId:                 "some-rate-id",
							ItemDescription:        "some-item-description",
							UsageStartDate:         "some-usage-start-date",
							UsageEndDate:           "some-usage-end-date",
							UsageQuantity:          0.51,
							BlendedRate:            "some-blended-rate",
							CurrencyCode:           "some-currency-code",
							CostBeforeTax:          "some-cost",
							Credits:                "some-credits",
							TaxAmount:              "some-tax-amount",
							TaxType:                "some-tax-type",
							TotalCost:              1.20,
						},
					)
				})

				It("skips them", func() {
					Expect(reports).To(HaveLen(1))
					Expect(reports[0].AccountName).To(Equal("some-linked-account-name"))
				})
			})
		})

		Context("with no reports", func() {
			It("returns empty", func() {
				reports := normalizer.Normalize(nil)

				Expect(reports).To(HaveLen(0))
			})
		})

		Context("with rows that are not linked line items", func() {
			BeforeEach(func() {
				usageReports = append(usageReports,
					&Usage{
						InvoiceID:              "some-invoice-id",
						PayerAccountId:         "some-payer-account-id",
						LinkedAccountId:        "",
						RecordType:             "PayerLineItem",
						RecordID:               "some-other-record-id",
						BillingPeriodStartDate: "some-start-date",
						BillingPeriodEndDate:   "some-end-date",
						InvoiceDate:            "some-invoice-date",
						PayerAccountName:       "",
						LinkedAccountName:      "some-other-linked-account-name",
						TaxationAddress:        "some-address",
						PayerPONumber:          "some-payer-number",
						ProductCode:            "some-other-product-code",
						ProductName:            "some-other-product-name",
						SellerOfRecord:         "some-other-seller-record",
						UsageType:              "some-other-usage-type",
						Operation:              "some-other-operation",
						RateId:                 "some-other-rate-id",
						ItemDescription:        "some-other-item-description",
						UsageStartDate:         "some-other-usage-start-date",
						UsageEndDate:           "some-other-usage-end-date",
						UsageQuantity:          0.12345,
						BlendedRate:            "some-other-blended-rate",
						CurrencyCode:           "some-other-currency-code",
						CostBeforeTax:          "some-other-cost",
						Credits:                "some-other-credits",
						TaxAmount:              "some-other-tax-amount",
						TaxType:                "some-other-tax-type",
						TotalCost:              13.37,
					},
				)
			})

			It("skips them", func() {
				Expect(reports).To(HaveLen(1))
				Expect(reports[0].AccountName).To(Equal("some-linked-account-name"))
			})
		})

		Context("with no reports", func() {
			It("returns empty", func() {
				reports := normalizer.Normalize(nil)

				Expect(reports).To(HaveLen(0))
			})
		})
	})
})
