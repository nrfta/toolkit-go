package sqlboiler_test

import (
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/sqlboiler"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Simple String Comparator", func() {
	var products []*models.Product

	BeforeEach(func() {
		var err error
		products, err = SeedProducts(ctx, db, 20)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := shared.TruncateAll(ctx, db)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ModsForSimpleStringComparator", func() {
		It("should filter by Eq (exact match only)", func() {
			targetName := products[0].Name
			mods := sqlboiler.ModsForSimpleStringComparator("products", "name", &comparators.SimpleString{
				Eq: &targetName,
			})

			results, err := models.Products(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].Name).To(Equal(targetName))
		})

		It("should filter by Neq", func() {
			targetName := products[0].Name
			mods := sqlboiler.ModsForSimpleStringComparator("products", "name", &comparators.SimpleString{
				Neq: &targetName,
			})

			results, err := models.Products(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(19))
			for _, result := range results {
				Expect(result.Name).NotTo(Equal(targetName))
			}
		})

		It("should filter by In", func() {
			names := []string{products[0].Name, products[1].Name, products[2].Name}
			mods := sqlboiler.ModsForSimpleStringComparator("products", "name", &comparators.SimpleString{
				In: names,
			})

			results, err := models.Products(mods...).All(ctx, db)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(3))
			for _, result := range results {
				Expect(result.Name).To(BeElementOf(names))
			}
		})

		It("should be case-sensitive for exact match", func() {
			targetName := products[0].Name
			uppercaseName := "PRODUCT 0"

			if targetName != uppercaseName {
				mods := sqlboiler.ModsForSimpleStringComparator("products", "name", &comparators.SimpleString{
					Eq: &uppercaseName,
				})

				results, err := models.Products(mods...).All(ctx, db)

				Expect(err).NotTo(HaveOccurred())
				Expect(results).To(BeEmpty())
			}
		})
	})
})
