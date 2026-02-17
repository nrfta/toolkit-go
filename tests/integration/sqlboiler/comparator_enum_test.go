package sqlboiler_test

import (
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Enum Comparator", func() {
	var repo *ProductRepository

	BeforeEach(func() {
		repo = NewProductRepository(db)
		_, err := SeedProducts(ctx, db, 20)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := shared.TruncateAll(ctx, db)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ModsForEnumComparator", func() {
		It("should filter by Eq", func() {
			status := "active"
			filters := []ProductFilter{
				{Status: &comparators.Enum[string]{Eq: &status}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Status).To(Equal("active"))
			}
		})

		It("should filter by Neq", func() {
			status := "deleted"
			filters := []ProductFilter{
				{Status: &comparators.Enum[string]{Neq: &status}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Status).NotTo(Equal("deleted"))
			}
		})

		It("should filter by In", func() {
			statuses := []string{"active", "inactive"}
			filters := []ProductFilter{
				{Status: &comparators.Enum[string]{In: statuses}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Status).To(BeElementOf(statuses))
			}
		})

		It("should filter by Nin", func() {
			statuses := []string{"deleted", "archived"}
			filters := []ProductFilter{
				{Status: &comparators.Enum[string]{Nin: statuses}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Status).NotTo(BeElementOf(statuses))
			}
		})

		It("should handle empty In list", func() {
			statuses := []string{}
			filters := []ProductFilter{
				{Status: &comparators.Enum[string]{In: statuses}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Empty IN list with no filtering returns all results
			Expect(results).To(HaveLen(20))
		})

		It("should handle single value in In list", func() {
			statuses := []string{"active"}
			filters := []ProductFilter{
				{Status: &comparators.Enum[string]{In: statuses}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.Status).To(Equal("active"))
			}
		})
	})
})
