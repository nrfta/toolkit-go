package sqlboiler_test

import (
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Boolean Comparator", func() {
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

	Describe("ModsForBooleanComparator", func() {
		It("should filter by Eq=true", func() {
			isAvailable := true
			filters := []ProductFilter{
				{IsAvailable: &comparators.Boolean{Eq: &isAvailable}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.IsAvailable).To(BeTrue())
			}
		})

		It("should filter by Eq=false", func() {
			isAvailable := false
			filters := []ProductFilter{
				{IsAvailable: &comparators.Boolean{Eq: &isAvailable}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(10))
			for _, result := range results {
				Expect(result.IsAvailable).To(BeFalse())
			}
		})

		It("should filter by Neq=true", func() {
			isAvailable := true
			filters := []ProductFilter{
				{IsAvailable: &comparators.Boolean{Neq: &isAvailable}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			for _, result := range results {
				Expect(result.IsAvailable).To(BeFalse())
			}
		})

		It("should filter by Neq=false", func() {
			isAvailable := false
			filters := []ProductFilter{
				{IsAvailable: &comparators.Boolean{Neq: &isAvailable}},
			}

			results, err := repo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.IsAvailable).To(BeTrue())
			}
		})

		It("should handle both Eq and Neq together", func() {
			isAvailable := true
			filters := []ProductFilter{
				{IsAvailable: &comparators.Boolean{
					Eq:  &isAvailable,
					Neq: &isAvailable,
				}},
			}

			results, err := repo.GetAll(ctx, filters...)

			if err == nil {
				Expect(results).To(BeEmpty())
			}
		})
	})
})
