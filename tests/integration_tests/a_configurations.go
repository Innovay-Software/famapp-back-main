package integrationTests

import "github.com/innovay-software/famapp-main/tests"

const (
	superAdminID       int64  = tests.SuperAdminID
	superAdminName     string = tests.SuperAdminName
	superAdminMobile   string = tests.SuperAdminMobile
	superAdminPassword string = tests.SuperAdminPassword
)

const (
	runTests     bool = tests.RunIntegrationTests
	runWorkflow1 bool = runTests && true
	runWorkflow2 bool = runTests && true
	runWorkflow3 bool = runTests && true
	runWorkflow4 bool = runTests && true
	runWorkflow5 bool = runTests && true
	runWorkflow6 bool = runTests && true
)
