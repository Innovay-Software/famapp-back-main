package unitTests

import (
	"github.com/innovay-software/famapp-main/app/utils"
)

func logWorkflowSuccess(content string) {
	utils.LogSuccess("\n***" + content + "***\n")
}
