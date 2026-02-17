package sqlboiler_test

import (
	"time"

	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Date Comparator Integration Tests", func() {
	var (
		orderRepo *OrderRepository
		userRepo  *UserRepository
	)

	BeforeEach(func() {
		orderRepo = NewOrderRepository(db)
		userRepo = NewUserRepository(db)
	})

	AfterEach(func() {
		err := shared.TruncateAll(ctx, db)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ModsForDateComparator with absolute dates", func() {
		BeforeEach(func() {
			users, err := SeedUsers(ctx, db, 5)
			Expect(err).NotTo(HaveOccurred())

			products, err := SeedProducts(ctx, db, 5)
			Expect(err).NotTo(HaveOccurred())

			_, err = SeedOrders(ctx, db, users, products, 10)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should filter by Eq", func() {
			today := time.Now().UTC().Truncate(24 * time.Hour).Format(time.RFC3339)

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Eq: &today}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			for _, result := range results {
				resultDate := result.OrderDate.UTC().Truncate(24 * time.Hour).Format(time.RFC3339)
				Expect(resultDate).To(Equal(today))
			}
		})

		It("should filter by Gte (greater than or equal)", func() {
			fiveDaysAgoTime := time.Now().UTC().Add(-5 * 24 * time.Hour)
			fiveDaysAgo := fiveDaysAgoTime.Format(time.RFC3339)

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Gte: &fiveDaysAgo}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.OrderDate.UTC()).To(BeTemporally(">=", fiveDaysAgoTime.Add(-1*time.Second)))
			}
		})

		It("should filter by Lt (less than)", func() {
			fiveDaysAgo := time.Now().Add(-5 * 24 * time.Hour).Format(time.RFC3339)

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Lt: &fiveDaysAgo}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			for _, result := range results {
				Expect(result.OrderDate).To(BeTemporally("<", time.Now().Add(-5*24*time.Hour)))
			}
		})

		It("should filter by Neq (not equal)", func() {
			specificDate := "2025-01-01T00:00:00Z"

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Neq: &specificDate}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(10))
			targetTime, _ := time.Parse(time.RFC3339, specificDate)
			for _, result := range results {
				Expect(result.OrderDate.UTC()).NotTo(Equal(targetTime))
			}
		})

		It("should filter by Gt (greater than)", func() {
			fiveDaysAgoTime := time.Now().UTC().Add(-5 * 24 * time.Hour)
			fiveDaysAgo := fiveDaysAgoTime.Format(time.RFC3339)

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Gt: &fiveDaysAgo}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			for _, result := range results {
				Expect(result.OrderDate.UTC()).To(BeTemporally(">", fiveDaysAgoTime.Add(-1*time.Second)))
			}
		})

		It("should filter by Lte (less than or equal)", func() {
			threeDaysAgo := time.Now().UTC().Add(-3 * 24 * time.Hour).Format(time.RFC3339)

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Lte: &threeDaysAgo}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			threeDaysAgoTime := time.Now().UTC().Add(-3 * 24 * time.Hour)
			for _, result := range results {
				Expect(result.OrderDate.UTC()).To(BeTemporally("<=", threeDaysAgoTime))
			}
		})

		It("should filter by range (Gte + Lt)", func() {
			sevenDaysAgoTime := time.Now().UTC().Add(-7 * 24 * time.Hour)
			oneDayAgoTime := time.Now().UTC().Add(-1 * 24 * time.Hour)
			sevenDaysAgo := sevenDaysAgoTime.Format(time.RFC3339)
			oneDayAgo := oneDayAgoTime.Format(time.RFC3339)

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{
					Gte: &sevenDaysAgo,
					Lt:  &oneDayAgo,
				}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			for _, result := range results {
				Expect(result.OrderDate.UTC()).To(BeTemporally(">=", sevenDaysAgoTime.Add(-1*time.Second)))
				Expect(result.OrderDate.UTC()).To(BeTemporally("<", oneDayAgoTime.Add(1*time.Second)))
			}
		})
	})

	Describe("ModsForDateComparator with durations", func() {
		BeforeEach(func() {
			users, err := SeedUsers(ctx, db, 5)
			Expect(err).NotTo(HaveOccurred())

			products, err := SeedProducts(ctx, db, 5)
			Expect(err).NotTo(HaveOccurred())

			_, err = SeedOrders(ctx, db, users, products, 20)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should parse negative duration (past) with -P format", func() {
			pastWeek := "-P7D"
			sevenDaysAgoTime := time.Now().Add(-7 * 24 * time.Hour)

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Gte: &pastWeek}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.OrderDate).To(BeTemporally(">=", sevenDaysAgoTime.Add(-1*time.Second)))
			}
		})

		It("should parse positive duration (future) with P format", func() {
			// Orders before 2 weeks from now: < P2W
			twoWeeksFuture := "P2W"

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Lt: &twoWeeksFuture}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
		})

		It("should handle month duration -P1M", func() {
			pastMonth := "-P1M"

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Gte: &pastMonth}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
		})

		It("should handle year duration -P1Y", func() {
			pastYear := "-P1Y"

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Gte: &pastYear}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
		})
	})

	Describe("ModsForNullableDateComparator", func() {
		BeforeEach(func() {
			users, err := SeedUsers(ctx, db, 5)
			Expect(err).NotTo(HaveOccurred())

			products, err := SeedProducts(ctx, db, 5)
			Expect(err).NotTo(HaveOccurred())

			_, err = SeedOrders(ctx, db, users, products, 20)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should filter NULL dates with Null=true", func() {
			nullVal := true
			filters := []OrderFilter{
				{ShippedDate: &comparators.NullableDate{Null: &nullVal}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, order := range results {
				Expect(order.ShippedDate.Valid).To(BeFalse())
			}
		})

		It("should filter non-NULL dates with Null=false", func() {
			nullVal := false
			filters := []OrderFilter{
				{ShippedDate: &comparators.NullableDate{Null: &nullVal}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, order := range results {
				Expect(order.ShippedDate.Valid).To(BeTrue())
			}
		})

		It("should filter by Eq (equals)", func() {
			specificDate := "2025-12-25T00:00:00Z"
			filters := []OrderFilter{
				{ShippedDate: &comparators.NullableDate{Eq: &specificDate}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			targetTime, _ := time.Parse(time.RFC3339, specificDate)
			for _, result := range results {
				if result.ShippedDate.Valid {
					Expect(result.ShippedDate.Time.UTC()).To(Equal(targetTime))
				}
			}
		})

		It("should filter by Neq (not equal)", func() {
			specificDate := "2025-01-01T00:00:00Z"

			filters := []OrderFilter{
				{ShippedDate: &comparators.NullableDate{Neq: &specificDate}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Should return orders with non-NULL shipped dates that don't match the specific date
			// (NULL values are excluded by Neq)
			Expect(results).NotTo(BeEmpty())
			targetTime, _ := time.Parse(time.RFC3339, specificDate)
			for _, result := range results {
				if result.ShippedDate.Valid {
					Expect(result.ShippedDate.Time.UTC()).NotTo(Equal(targetTime))
				}
			}
		})

		It("should filter by Gt (greater than)", func() {
			fiveDaysAgoTime := time.Now().UTC().Add(-5 * 24 * time.Hour)
			fiveDaysAgo := fiveDaysAgoTime.Format(time.RFC3339)

			filters := []OrderFilter{
				{ShippedDate: &comparators.NullableDate{Gt: &fiveDaysAgo}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			for _, result := range results {
				Expect(result.ShippedDate.Valid).To(BeTrue())
				Expect(result.ShippedDate.Time.UTC()).To(BeTemporally(">", fiveDaysAgoTime.Add(-1*time.Second)))
			}
		})

		It("should filter by Lte (less than or equal)", func() {
			threeDaysAgo := time.Now().UTC().Add(-3 * 24 * time.Hour).Format(time.RFC3339)

			filters := []OrderFilter{
				{ShippedDate: &comparators.NullableDate{Lte: &threeDaysAgo}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			threeDaysAgoTime := time.Now().UTC().Add(-3 * 24 * time.Hour)
			for _, result := range results {
				Expect(result.ShippedDate.Valid).To(BeTrue())
				Expect(result.ShippedDate.Time.UTC()).To(BeTemporally("<=", threeDaysAgoTime))
			}
		})

		It("should filter by Lt (less than)", func() {
			fourDaysAgo := time.Now().UTC().Add(-4 * 24 * time.Hour).Format(time.RFC3339)

			filters := []OrderFilter{
				{ShippedDate: &comparators.NullableDate{Lt: &fourDaysAgo}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			fourDaysAgoTime := time.Now().UTC().Add(-4 * 24 * time.Hour)
			for _, result := range results {
				Expect(result.ShippedDate.Valid).To(BeTrue())
				Expect(result.ShippedDate.Time.UTC()).To(BeTemporally("<", fourDaysAgoTime))
			}
		})

		It("should combine Null and Gte for OR logic", func() {
			sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour).Format(time.RFC3339)
			nullVal := true

			filters := []OrderFilter{
				{ShippedDate: &comparators.NullableDate{
					Null: &nullVal,
					Gte:  &sevenDaysAgo,
				}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			// Results should be either NULL or >= 7 days ago
			for _, order := range results {
				if order.ShippedDate.Valid {
					Expect(order.ShippedDate.Time).To(BeTemporally(">=", time.Now().Add(-7*24*time.Hour)))
				}
			}
		})
	})

	Describe("Date comparators on Users created_at", func() {
		BeforeEach(func() {
			_, err := SeedUsers(ctx, db, 15)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should filter users created in the last 10 hours", func() {
			tenHoursAgo := time.Now().Add(-10 * time.Hour).Format(time.RFC3339)

			filters := []UserFilter{
				{CreatedAt: &comparators.Date{Gte: &tenHoursAgo}},
			}

			results, err := userRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			Expect(len(results)).To(BeNumerically("<=", 11)) // Users 0-10
		})

		It("should filter users created more than 5 hours ago", func() {
			fiveHoursAgo := time.Now().Add(-5 * time.Hour).Format(time.RFC3339)

			filters := []UserFilter{
				{CreatedAt: &comparators.Date{Lt: &fiveHoursAgo}},
			}

			results, err := userRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
		})
	})
})
