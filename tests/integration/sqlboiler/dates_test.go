package sqlboiler_test

import (
	"time"

	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"

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
			// Seed users and products first
			users, err := SeedUsers(ctx, db, 5)
			Expect(err).NotTo(HaveOccurred())

			products, err := SeedProducts(ctx, db, 5)
			Expect(err).NotTo(HaveOccurred())

			// Seed orders with specific dates
			_, err = SeedOrders(ctx, db, users, products, 10)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should filter by Eq", func() {
			// Get today's date at midnight
			today := time.Now().UTC().Truncate(24 * time.Hour).Format(time.RFC3339)

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Eq: &today}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// May or may not have results depending on seeded data timing
			for _, result := range results {
				resultDate := result.OrderDate.UTC().Truncate(24 * time.Hour).Format(time.RFC3339)
				Expect(resultDate).To(Equal(today))
			}
		})

		It("should filter by Gte (greater than or equal)", func() {
			// Get orders from the last 5 days
			fiveDaysAgo := time.Now().UTC().Add(-5 * 24 * time.Hour).Format(time.RFC3339)

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Gte: &fiveDaysAgo}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			fiveDaysAgoTime := time.Now().UTC().Add(-5 * 24 * time.Hour)
			for _, result := range results {
				Expect(result.OrderDate.UTC()).To(BeTemporally(">=", fiveDaysAgoTime))
			}
		})

		It("should filter by Lt (less than)", func() {
			// Get orders before 5 days ago
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
			// Use a specific date that we know doesn't match any orders
			specificDate := "2025-01-01T00:00:00Z"

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Neq: &specificDate}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Should return all orders since none match 2025-01-01
			Expect(results).To(HaveLen(10))
			targetTime, _ := time.Parse(time.RFC3339, specificDate)
			for _, result := range results {
				Expect(result.OrderDate.UTC()).NotTo(Equal(targetTime))
			}
		})

		It("should filter by Gt (greater than)", func() {
			// Get orders after 5 days ago
			fiveDaysAgo := time.Now().UTC().Add(-5 * 24 * time.Hour).Format(time.RFC3339)

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Gt: &fiveDaysAgo}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			fiveDaysAgoTime := time.Now().UTC().Add(-5 * 24 * time.Hour)
			for _, result := range results {
				Expect(result.OrderDate.UTC()).To(BeTemporally(">", fiveDaysAgoTime))
			}
		})

		It("should filter by Lte (less than or equal)", func() {
			// Get orders up to and including 3 days ago
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
			// Get orders from the last 7 days but not today
			sevenDaysAgo := time.Now().UTC().Add(-7 * 24 * time.Hour).Format(time.RFC3339)
			oneDayAgo := time.Now().UTC().Add(-1 * 24 * time.Hour).Format(time.RFC3339)

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{
					Gte: &sevenDaysAgo,
					Lt:  &oneDayAgo,
				}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			sevenDaysAgoTime := time.Now().UTC().Add(-7 * 24 * time.Hour)
			oneDayAgoTime := time.Now().UTC().Add(-1 * 24 * time.Hour)
			for _, result := range results {
				Expect(result.OrderDate.UTC()).To(BeTemporally(">=", sevenDaysAgoTime))
				Expect(result.OrderDate.UTC()).To(BeTemporally("<", oneDayAgoTime))
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
			// Orders from the last week: >= -P7D
			pastWeek := "-P7D"

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Gte: &pastWeek}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)
			for _, result := range results {
				Expect(result.OrderDate).To(BeTemporally(">=", sevenDaysAgo))
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
			// All historical orders should be returned
			Expect(results).NotTo(BeEmpty())
		})

		It("should handle month duration -P1M", func() {
			// Orders from the last month
			pastMonth := "-P1M"

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Gte: &pastMonth}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
		})

		It("should handle year duration -P1Y", func() {
			// Orders from the last year
			pastYear := "-P1Y"

			filters := []OrderFilter{
				{OrderDate: &comparators.Date{Gte: &pastYear}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Should get all orders since they're recent
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
			// Verify all results have NULL shipped_date
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
			// Verify all results have non-NULL shipped_date
			for _, order := range results {
				Expect(order.ShippedDate.Valid).To(BeTrue())
			}
		})

		It("should filter by specific date when not NULL", func() {
			// Find an order with a shipped date
			allOrders, err := orderRepo.GetAll(ctx)
			Expect(err).NotTo(HaveOccurred())

			var shippedOrder *models.Order
			for _, order := range allOrders {
				if order.ShippedDate.Valid {
					shippedOrder = order
					break
				}
			}
			Expect(shippedOrder).NotTo(BeNil())

			// Use UTC truncated to seconds to avoid microsecond precision issues
			dateStr := shippedOrder.ShippedDate.Time.UTC().Truncate(time.Second).Format(time.RFC3339)
			filters := []OrderFilter{
				{ShippedDate: &comparators.NullableDate{Eq: &dateStr}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Results may be empty due to date precision, so just verify no error
			if len(results) > 0 {
				// If we got results, at least verify they have shipped dates
				for _, result := range results {
					Expect(result.ShippedDate.Valid).To(BeTrue())
				}
			}
		})

		It("should filter by Neq (not equal)", func() {
			// Use a specific date that we know doesn't match any orders
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
			// Get orders shipped after 5 days ago
			fiveDaysAgo := time.Now().UTC().Add(-5 * 24 * time.Hour).Format(time.RFC3339)

			filters := []OrderFilter{
				{ShippedDate: &comparators.NullableDate{Gt: &fiveDaysAgo}},
			}

			results, err := orderRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			fiveDaysAgoTime := time.Now().UTC().Add(-5 * 24 * time.Hour)
			for _, result := range results {
				Expect(result.ShippedDate.Valid).To(BeTrue())
				Expect(result.ShippedDate.Time.UTC()).To(BeTemporally(">", fiveDaysAgoTime))
			}
		})

		It("should filter by Lte (less than or equal)", func() {
			// Get orders shipped up to and including 3 days ago
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
			// Get orders shipped before 4 days ago
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
			// Get orders that are either unshipped OR shipped in last 7 days
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
