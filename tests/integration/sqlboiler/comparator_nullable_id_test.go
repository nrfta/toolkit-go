package sqlboiler_test

import (
	"github.com/nrfta/toolkit-go/comparators"
	"github.com/nrfta/toolkit-go/tests/integration/shared"
	"github.com/nrfta/toolkit-go/tests/integration/sqlboiler/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nullable ID Comparator", func() {
	var (
		userRepo *UserRepository
		users    []*models.User
	)

	BeforeEach(func() {
		userRepo = NewUserRepository(db)
		var err error
		users, err = SeedUsers(ctx, db, 20)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := shared.TruncateAll(ctx, db)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ModsForNullableIDComparator", func() {
		It("should filter NULL values with Null=true", func() {
			nullVal := true
			filters := []UserFilter{
				{OrganizationID: &comparators.NullableID{Null: &nullVal}},
			}

			results, err := userRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.OrganizationID.Valid).To(BeFalse())
			}
		})

		It("should filter non-NULL values with Null=false", func() {
			nullVal := false
			filters := []UserFilter{
				{OrganizationID: &comparators.NullableID{Null: &nullVal}},
			}

			results, err := userRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.OrganizationID.Valid).To(BeTrue())
			}
		})

		It("should filter by Eq when value is set", func() {
			// Find a user with an organization ID
			var targetUser *models.User
			for _, u := range users {
				if u.OrganizationID.Valid {
					targetUser = u
					break
				}
			}
			Expect(targetUser).NotTo(BeNil())

			orgID := targetUser.OrganizationID.String
			filters := []UserFilter{
				{OrganizationID: &comparators.NullableID{Eq: &orgID}},
			}

			results, err := userRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.OrganizationID.Valid).To(BeTrue())
				Expect(result.OrganizationID.String).To(Equal(orgID))
			}
		})

		It("should filter by Neq", func() {
			// Find a user with an organization ID
			var targetUser *models.User
			for _, u := range users {
				if u.OrganizationID.Valid {
					targetUser = u
					break
				}
			}
			Expect(targetUser).NotTo(BeNil())

			orgID := targetUser.OrganizationID.String
			filters := []UserFilter{
				{OrganizationID: &comparators.NullableID{Neq: &orgID}},
			}

			results, err := userRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Should return all users except those with the target org ID
			// (including NULL org IDs since Neq doesn't match NULL)
			for _, result := range results {
				if result.OrganizationID.Valid {
					Expect(result.OrganizationID.String).NotTo(Equal(orgID))
				}
			}
		})

		It("should filter by In", func() {
			// Find two users with organization IDs
			var validUsers []*models.User
			for _, u := range users {
				if u.OrganizationID.Valid {
					validUsers = append(validUsers, u)
					if len(validUsers) == 2 {
						break
					}
				}
			}
			Expect(validUsers).To(HaveLen(2))

			orgIDs := []string{
				validUsers[0].OrganizationID.String,
				validUsers[1].OrganizationID.String,
			}
			filters := []UserFilter{
				{OrganizationID: &comparators.NullableID{In: orgIDs}},
			}

			results, err := userRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			for _, result := range results {
				Expect(result.OrganizationID.Valid).To(BeTrue())
				Expect(result.OrganizationID.String).To(BeElementOf(orgIDs))
			}
		})

		It("should filter by Nin", func() {
			// Find two users with organization IDs
			var validUsers []*models.User
			for _, u := range users {
				if u.OrganizationID.Valid {
					validUsers = append(validUsers, u)
					if len(validUsers) == 2 {
						break
					}
				}
			}
			Expect(validUsers).To(HaveLen(2))

			orgIDs := []string{
				validUsers[0].OrganizationID.String,
				validUsers[1].OrganizationID.String,
			}
			filters := []UserFilter{
				{OrganizationID: &comparators.NullableID{Nin: orgIDs}},
			}

			results, err := userRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			// Should return all users except those with the specified org IDs
			for _, result := range results {
				if result.OrganizationID.Valid {
					Expect(result.OrganizationID.String).NotTo(BeElementOf(orgIDs))
				}
			}
		})

		It("should combine Null with Eq for OR logic", func() {
			// Find a user with an organization ID
			var targetUser *models.User
			for _, u := range users {
				if u.OrganizationID.Valid {
					targetUser = u
					break
				}
			}
			Expect(targetUser).NotTo(BeNil())

			orgID := targetUser.OrganizationID.String
			nullVal := true
			filters := []UserFilter{
				{OrganizationID: &comparators.NullableID{
					Eq:   &orgID,
					Null: &nullVal,
				}},
			}

			results, err := userRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			// Results should be either NULL or matching the org ID
			for _, result := range results {
				if result.OrganizationID.Valid {
					Expect(result.OrganizationID.String).To(Equal(orgID))
				}
			}
		})

		It("should combine Null with In for OR logic", func() {
			// Find two users with organization IDs
			var validUsers []*models.User
			for _, u := range users {
				if u.OrganizationID.Valid {
					validUsers = append(validUsers, u)
					if len(validUsers) == 2 {
						break
					}
				}
			}
			Expect(validUsers).To(HaveLen(2))

			orgIDs := []string{
				validUsers[0].OrganizationID.String,
				validUsers[1].OrganizationID.String,
			}
			nullVal := true
			filters := []UserFilter{
				{OrganizationID: &comparators.NullableID{
					In:   orgIDs,
					Null: &nullVal,
				}},
			}

			results, err := userRepo.GetAll(ctx, filters...)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			// Results should be either NULL or in the org ID list
			for _, result := range results {
				if result.OrganizationID.Valid {
					Expect(result.OrganizationID.String).To(BeElementOf(orgIDs))
				}
			}
		})
	})
})
