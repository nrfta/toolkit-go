package sqlboiler_test

import (
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/sqlboiler"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ID Comparator With OR", func() {
	var (
		users    []*models.User
		products []*models.Product
	)

	BeforeEach(func() {
		var err error
		users, err = SeedUsers(ctx, db, 5)
		Expect(err).NotTo(HaveOccurred())

		products, err = SeedProducts(ctx, db, 5)
		Expect(err).NotTo(HaveOccurred())

		_, err = SeedOrders(ctx, db, users, products, 10)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := shared.TruncateAll(ctx, db)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ModsForIDComparatorWithOr", func() {
		It("should match records where either column equals the value", func() {
			// Find orders where user_id OR product_id matches the first user's ID
			targetID := users[0].ID

			mods := sqlboiler.ModsForIDComparatorWithOr(
				"orders",
				"user_id",
				"product_id",
				&comparators.ID{Eq: &targetID},
			)

			results, err := models.Orders(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			// Each result should have user_id == targetID OR product_id == targetID
			for _, order := range results {
				hasMatch := order.UserID == targetID || order.ProductID == targetID
				Expect(hasMatch).To(BeTrue())
			}
		})

		It("should match records where either column is in the list", func() {
			// Find orders where user_id OR product_id is in the list
			ids := []string{users[0].ID, users[1].ID}

			mods := sqlboiler.ModsForIDComparatorWithOr(
				"orders",
				"user_id",
				"product_id",
				&comparators.ID{In: ids},
			)

			results, err := models.Orders(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			// Each result should have user_id in ids OR product_id in ids
			for _, order := range results {
				hasMatch := shared.Contains(ids, order.UserID) || shared.Contains(ids, order.ProductID)
				Expect(hasMatch).To(BeTrue())
			}
		})

		It("should match records where BOTH columns do not equal the value (AND logic)", func() {
			// With AND logic: user_id != targetID AND product_id != targetID
			// This excludes records where EITHER column matches
			targetID := users[0].ID

			mods := sqlboiler.ModsForIDComparatorWithOr(
				"orders",
				"user_id",
				"product_id",
				&comparators.ID{Neq: &targetID},
			)

			results, err := models.Orders(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			// Each result should have user_id != targetID AND product_id != targetID
			for _, order := range results {
				Expect(order.UserID).NotTo(Equal(targetID))
				Expect(order.ProductID).NotTo(Equal(targetID))
			}
		})

		It("should match records where BOTH columns are not in the list (AND logic)", func() {
			// With AND logic: user_id NOT IN ids AND product_id NOT IN ids
			// This excludes records where EITHER column matches
			ids := []string{users[0].ID, users[1].ID}

			mods := sqlboiler.ModsForIDComparatorWithOr(
				"orders",
				"user_id",
				"product_id",
				&comparators.ID{Nin: ids},
			)

			results, err := models.Orders(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			// Each result should have user_id NOT IN ids AND product_id NOT IN ids
			for _, order := range results {
				Expect(shared.Contains(ids, order.UserID)).To(BeFalse())
				Expect(shared.Contains(ids, order.ProductID)).To(BeFalse())
			}
		})

		It("should handle nil comparator", func() {
			mods := sqlboiler.ModsForIDComparatorWithOr(
				"orders",
				"user_id",
				"product_id",
				nil,
			)

			Expect(mods).To(BeNil())
		})

		It("should combine with other filters using AND", func() {
			// Find orders where (user_id=X OR product_id=X) AND status='pending'
			targetID := users[0].ID

			orMods := sqlboiler.ModsForIDComparatorWithOr(
				"orders",
				"user_id",
				"product_id",
				&comparators.ID{Eq: &targetID},
			)

			status := "pending"
			statusMods := sqlboiler.ModsForEnumComparator(
				"orders",
				"status",
				&comparators.Enum[string]{Eq: &status},
			)

			allMods := append(orMods, statusMods...)
			results, err := models.Orders(allMods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			// Each result must satisfy: (user_id==X OR product_id==X) AND status=='pending'
			for _, order := range results {
				hasMatch := order.UserID == targetID || order.ProductID == targetID
				Expect(hasMatch).To(BeTrue())
				Expect(order.Status).To(Equal("pending"))
			}
		})

		It("should work for polymorphic relationships", func() {
			// Common use case: finding records by polymorphic owner
			// Example: WHERE created_by_id = X OR updated_by_id = X
			// We'll simulate this with user_id OR product_id (same concept)

			// Find orders associated with first user's ID in any capacity
			targetID := users[0].ID

			mods := sqlboiler.ModsForIDComparatorWithOr(
				"orders",
				"user_id",
				"product_id",  // In real polymorphic case, might be same column checked twice
				&comparators.ID{Eq: &targetID},
			)

			results, err := models.Orders(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
		})

		It("should handle Eq and In in same comparator", func() {
			// When both Eq and In are set, both conditions should be applied
			eqID := users[0].ID
			inIDs := []string{users[1].ID, users[2].ID}

			mods := sqlboiler.ModsForIDComparatorWithOr(
				"orders",
				"user_id",
				"product_id",
				&comparators.ID{
					Eq: &eqID,
					In: inIDs,
				},
			)

			results, err := models.Orders(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			// Should match orders where:
			// (user_id==eqID OR product_id==eqID) AND
			// (user_id IN inIDs OR product_id IN inIDs)
			for _, order := range results {
				hasEqMatch := order.UserID == eqID || order.ProductID == eqID
				hasInMatch := shared.Contains(inIDs, order.UserID) || shared.Contains(inIDs, order.ProductID)
				// Both conditions must be true (AND)
				Expect(hasEqMatch && hasInMatch).To(BeTrue())
			}
		})
	})
})
