package main

import (
	"fmt"

	"github.com/zamirka/interfaxScanApi/methods"
	"github.com/zamirka/interfaxScanApi/utils"
)

var ctx utils.AppContext

func main() {
	err := utils.InitExecutionContext(&ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = methods.Login(&ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	myq, err := methods.GetAllQueries(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(myq)
}
