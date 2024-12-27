package unitTests

import "github.com/innovay-software/famapp-main/tests"

const (
	superAdminID       uint64 = tests.SuperAdminID
	superAdminName     string = tests.SuperAdminName
	superAdminMobile   string = tests.SuperAdminMobile
	superAdminPassword string = tests.SuperAdminPassword
)

const (
	runTests                  bool = tests.RunUnitTests
	runFolderFileServiceTests bool = runTests && true
)
