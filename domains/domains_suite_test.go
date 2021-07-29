package domains

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/cloudfoundry-incubator/cf-performance-tests/helpers"
	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var testConfig helpers.Config = helpers.NewConfig()
var testSetup *workflowhelpers.ReproducibleTestSuiteSetup
var ccdb *sql.DB
var uaadb *sql.DB
var ctx context.Context
const (
	orgs = 1000
	sharedDomains = 1000
	privateDomains = 1000
)

var _ = BeforeSuite(func() {
	testSetup = workflowhelpers.NewTestSuiteSetup(&testConfig)
	testSetup.Setup()
	ccdb, uaadb, ctx = helpers.OpenDbConnections(testConfig.CcdbConnection, testConfig.UaaConnection)

	quotaId := helpers.ExecuteSelectStatementOneRow(ccdb, ctx, "SELECT id FROM quota_definitions WHERE name = 'default'")
	var organizationIds []int

	for i := 0; i < orgs; i++ {
		guid := uuid.New()
		name := testConfig.NamePrefix + "-org-" + guid.String()
		statement := fmt.Sprintf("INSERT INTO organizations (guid, name, quota_definition_id) VALUES ('%s', '%s', %d) RETURNING id", guid.String(), name, quotaId)
		organizationId := helpers.ExecuteInsertStatement(ccdb, ctx, statement)
		organizationIds = append(organizationIds, organizationId)
	}
	for i := 0; i<sharedDomains; i++ {
		sharedDomainGuid := uuid.New()
		sharedDomainName := testConfig.NamePrefix + "-shareddomain-" + sharedDomainGuid.String()
		statement := fmt.Sprintf("INSERT INTO domains (guid, name) VALUES ('%s', '%s') RETURNING id", sharedDomainGuid.String(), sharedDomainName)
		helpers.ExecuteInsertStatement(ccdb, ctx, statement)
	}

	for i := 0; i<privateDomains; i++ {
		privateDomainGuid := uuid.New()
		privateDomainName := testConfig.NamePrefix + "-privatedomain-" + privateDomainGuid.String()
		owningOrganizationId := organizationIds[rand.Intn(len(organizationIds))]
		statement := fmt.Sprintf("INSERT INTO domains (guid, name, owning_organization_id) VALUES ('%s', '%s', %d) RETURNING id", privateDomainGuid.String(), privateDomainName, owningOrganizationId)
		helpers.ExecuteInsertStatement(ccdb, ctx, statement)
	}

})

var _ = AfterSuite(func() {

	helpers.CleanupTestData(ccdb, uaadb, ctx)

	err := ccdb.Close()
	if err != nil {
		log.Print(err)
	}

	err = uaadb.Close()
	if err != nil {
		log.Print(err)
	}
})

func TestDomains(t *testing.T) {
	viper.SetConfigName("config")
	viper.AddConfigPath("..")
	viper.AddConfigPath("$HOME/.cf-performance-tests")
	err := viper.ReadInConfig()
	if err != nil {
		t.Fatalf("error loading config: %s", err.Error())
	}

	err = viper.Unmarshal(&testConfig)
	if err != nil {
		t.Fatalf("error parsing config: %s", err.Error())
	}

	timestamp := time.Now().Unix()
	jsonReporter := helpers.NewJsonReporter(fmt.Sprintf("../test-results/domains-test-results-%d.json", timestamp), testConfig.CfDeploymentVersion, timestamp)

	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "DomainsTest Suite", []Reporter{jsonReporter})
}

