package sqlboiler_test

import (
	"context"
	"database/sql"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nrfta/toolkit-go/tests/integration/shared"
)

var (
	db        *sql.DB
	container *shared.PostgresContainer
	ctx       context.Context
)

func TestSQLBoilerIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SQLBoiler Integration Suite")
}

var _ = BeforeSuite(func() {
	ctx = context.Background()

	// Start PostgreSQL container
	var err error
	container, err = shared.StartPostgresContainer(ctx)
	Expect(err).NotTo(HaveOccurred())

	// Get database connection
	db = container.DB

	// Load schema
	err = shared.LoadSchema(ctx, db, "schema.sql")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	if container != nil {
		err := container.Terminate(ctx)
		Expect(err).NotTo(HaveOccurred())
	}
})
