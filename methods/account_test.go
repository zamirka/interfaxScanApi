package methods

import (
	"testing"

	"github.com/zamirka/interfaxScanApi/utils"
)

func TestLogin(t *testing.T) {
	var err error
	var ctx utils.AppContext
	if err = utils.InitExecutionContext("../testconf.json", &ctx); err != nil {
		t.Error(err)
	}
	if err = Login(&ctx); err != nil {
		t.Error(err)
	}
}
