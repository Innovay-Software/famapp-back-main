package integrationTests

import "github.com/innovay-software/famapp-main/tests"

const (
	superAdminID       uint64 = tests.SuperAdminID
	superAdminName     string = tests.SuperAdminName
	superAdminMobile   string = tests.SuperAdminMobile
	superAdminPassword string = tests.SuperAdminPassword
)

const (
	runTests     bool = tests.RunIntegrationTests
	runWorkflow1 bool = runTests && false
	runWorkflow2 bool = runTests && false
	runWorkflow3 bool = runTests && false
	runWorkflow4 bool = runTests && true
	runWorkflow5 bool = runTests && false
	runWorkflow6 bool = runTests && false
)
