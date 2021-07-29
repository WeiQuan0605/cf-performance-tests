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
var ctx context.Context

var _ = BeforeSuite(func() {
	testSetup = workflowhelpers.NewTestSuiteSetup(&testConfig)
	testSetup.Setup()
	ccdb, err := sql.Open("pgx", testConfig.CcdbConnection)
	if err != nil {
		log.Fatal(err)
	}
	defer ccdb.Close()

	ctx = context.Background()

	quotaId := helpers.ExecuteSelectStatementOneRow(ccdb, ctx, "SELECT id FROM quota_definitions WHERE name = 'default'")
	var organizationIds []int

	for i := 0; i < 100; i++ {
		guid := uuid.New()
		name := testConfig.NamePrefix + "-org-" + guid.String()
		statement := fmt.Sprintf("INSERT INTO organizations (guid, name, quota_definition_id) VALUES ('%s', '%s', %d) RETURNING id", guid.String(), name, quotaId)
		organizationId := helpers.ExecuteInsertStatement(ccdb, ctx, statement)
		organizationIds = append(organizationIds, organizationId)
	}
	for i := 0; i<100;i++ {
		sharedDomainGuid := uuid.New()
		sharedDomainName := testConfig.NamePrefix + "-shareddomain-" + sharedDomainGuid.String()
		statement := fmt.Sprintf("INSERT INTO domains (guid, name) VALUES ('%s', '%s') RETURNING id", sharedDomainGuid.String(), sharedDomainName)
		helpers.ExecuteInsertStatement(ccdb, ctx, statement)
	}

	for i := 0; i<100;i++ {
		privateDomainGuid := uuid.New()
		privateDomainName := testConfig.NamePrefix + "-privatedomain-" + privateDomainGuid.String()
		owningOrganizationId := organizationIds[rand.Intn(len(organizationIds))]
		statement := fmt.Sprintf("INSERT INTO domains (guid, name, owning_organization_id) VALUES ('%s', '%s', %d) RETURNING id", privateDomainGuid.String(), privateDomainName, owningOrganizationId)
		helpers.ExecuteInsertStatement(ccdb, ctx, statement)
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

