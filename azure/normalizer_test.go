package azure_test

import (
	"time"

	"github.com/Sirupsen/logrus"
	. "github.com/challiwill/meteorologica/azure"
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
		normalizer = NewNormalizer(log, loc)
	})

	Describe("Normalize", func() {
		var (
			usageReports []*Usage
			reports      datamodels.Reports
		)

		BeforeEach(func() {
			usageReports = []*Usage{
				&Usage{
					AccountOwnerId:         "some-owner",
					AccountName:            "some-account",
					ServiceAdministratorId: "some-administrator-id",
					SubscriptionId:         "some-subscription-id",
					SubscriptionGuid:       "some-guid",
					SubscriptionName:       "some-name",
					Date:                   "10/01/2016",
					Month:                  "10",
					Day:                    "1",
					Year:                   "2016",
					Product:                "some-product",
					MeterID:                "some-meter",
					MeterCategory:          "some-category",
					MeterSubCategory:       "some-sub-category",
					MeterRegion:            "some-region",
					MeterName:              "some-meter-name",
					ConsumedQuantity:       "24.00",
					ResourceRate:           "0.01",
					ExtendedCost:           "0.02",
					ResourceLocation:       "westus",
					ConsumedService:        "some-service-type",
					InstanceID:             "some-instance-id",
					ServiceInfo1:           "some-info",
					ServiceInfo2:           "some-other-info",
					AdditionalInfo:         "some-really-other-info",
					Tags:                   "some-tags",
					StoreServiceIdentifier: "some-identifier",
					DepartmentName:         "some-department-name",
					CostCenter:             "some-cost-center",
					UnitOfMeasure:          "Hours",
					ResourceGroup:          "some-group",
				},
				&Usage{
					AccountOwnerId:         "some-owner",
					AccountName:            "some-other-account",
					ServiceAdministratorId: "some-other-administrator-id",
					SubscriptionId:         "some-other-subscription-id",
					SubscriptionGuid:       "some-other-guid",
					SubscriptionName:       "some-other-name",
					Date:                   "10/01/2016",
					Month:                  "10",
					Day:                    "1",
					Year:                   "2016",
					Product:                "some-other-product",
					MeterID:                "some-other-meter",
					MeterCategory:          "some-other-category",
					MeterSubCategory:       "some-other-sub-category",
					MeterRegion:            "some-other-region",
					MeterName:              "some-other-meter-name",
					ConsumedQuantity:       "22.00",
					ResourceRate:           "2.01",
					ExtendedCost:           "4.02",
					ResourceLocation:       "eastus",
					ConsumedService:        "some-other-service-type",
					InstanceID:             "some-other-instance-id",
					ServiceInfo1:           "some-other-info",
					ServiceInfo2:           "other-info",
					AdditionalInfo:         "really-other-info",
					Tags:                   "some-other-tags",
					StoreServiceIdentifier: "some-other-identifier",
					DepartmentName:         "some-other-department-name",
					CostCenter:             "some-other-cost-center",
					UnitOfMeasure:          "Hours",
					ResourceGroup:          "some-other-group",
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
						AccountNumber: "some-guid",
						AccountName:   "some-name",
						Day:           1,
						Month:         "October",
						Year:          2016,
						ServiceType:   "some-service-type",
						UsageQuantity: 24.00,
						Cost:          0.02,
						Region:        "some-region",
						UnitOfMeasure: "Hours",
						IAAS:          "Azure",
					}))
					Expect(reports[1]).To(Equal(datamodels.Report{
						AccountNumber: "some-other-guid",
						AccountName:   "some-other-name",
						Day:           1,
						Month:         "October",
						Year:          2016,
						ServiceType:   "some-other-service-type",
						UsageQuantity: 22.00,
						Cost:          4.02,
						Region:        "some-other-region",
						UnitOfMeasure: "Hours",
						IAAS:          "Azure",
					}))
				})
			})

			Context("with invalid extended cost", func() {
				BeforeEach(func() {
					usageReports[0].ExtendedCost = "not-a-float"
				})

				It("returns the same number of reports", func() {
					Expect(reports).To(HaveLen((len(usageReports))))
				})

				It("warns that extended cost is invalid", func() {
					Expect(log.Out).To(Say("extended cost 'not-a-float' is invalid"))
				})

				It("returns normalized Reports cost set to neutral value 0", func() {
					Expect(reports[0].Cost).To(Equal(float64(0)))
				})
			})

			Context("with invalid consumed quantity", func() {
				BeforeEach(func() {
					usageReports[0].ConsumedQuantity = "not-a-float"
				})

				It("returns the same number of reports", func() {
					Expect(reports).To(HaveLen((len(usageReports))))
				})

				It("warns that consumed quantity is invalid", func() {
					Expect(log.Out).To(Say("consumed quantity 'not-a-float' is invalid"))
				})

				It("returns normalized Reports usage quantity set to neutral value 0", func() {
					Expect(reports[0].UsageQuantity).To(Equal(float64(0)))
				})
			})

			Context("with invalid day", func() {
				Context("when not an int", func() {
					BeforeEach(func() {
						usageReports[0].Day = "not-an-int"
					})

					It("returns the same number of reports", func() {
						Expect(reports).To(HaveLen((len(usageReports))))
					})

					It("warns that day is invalid", func() {
						Expect(log.Out).To(Say("day is invalid"))
					})

					It("returns normalized Reports with day set to today", func() {
						Expect(reports[0].Day).To(Equal(time.Now().Day()))
					})
				})

				Context("when invalid int", func() {
					BeforeEach(func() {
						usageReports[0].Day = "33"
					})

					It("returns the same number of reports", func() {
						Expect(reports).To(HaveLen((len(usageReports))))
					})

					It("warns that day is invalid", func() {
						Expect(log.Out).To(Say("day is invalid"))
					})

					It("returns normalized Reports with day set to today", func() {
						Expect(reports[0].Day).To(Equal(time.Now().Day()))
					})
				})
			})

			Context("with invalid month", func() {
				Context("when not an int", func() {
					BeforeEach(func() {
						usageReports[0].Month = "not-an-int"
					})

					It("returns the same number of reports", func() {
						Expect(reports).To(HaveLen((len(usageReports))))
					})

					It("warns that month is invalid", func() {
						Expect(log.Out).To(Say("month is invalid"))
					})

					It("returns normalized Reports with month set to today", func() {
						Expect(reports[0].Month).To(Equal(time.Now().Month().String()))
					})
				})

				Context("when invalid int", func() {
					BeforeEach(func() {
						usageReports[0].Month = "13"
					})

					It("returns the same number of reports", func() {
						Expect(reports).To(HaveLen((len(usageReports))))
					})

					It("warns that month is invalid", func() {
						Expect(log.Out).To(Say("month is invalid"))
					})

					It("returns normalized Reports with month set to today", func() {
						Expect(reports[0].Month).To(Equal(time.Now().Month().String()))
					})
				})

				Context("with invalid year", func() {
					Context("when not an int", func() {
						BeforeEach(func() {
							usageReports[0].Year = "not-an-int"
						})

						It("returns the same number of reports", func() {
							Expect(reports).To(HaveLen((len(usageReports))))
						})

						It("warns that year is invalid", func() {
							Expect(log.Out).To(Say("year is invalid"))
						})

						It("returns normalized Reports with year set to today", func() {
							Expect(reports[0].Year).To(Equal(time.Now().Year()))
						})
					})
				})
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
