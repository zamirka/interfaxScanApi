package main

import (
	"fmt"

	"github.com/zamirka/interfaxScanApi/methods"
	"github.com/zamirka/interfaxScanApi/utils"
)

var ctx utils.AppContext

func main() {
	var err error
	if err = utils.InitExecutionContext(&ctx); err != nil {
		fmt.Println(err)
		return
	}

	if err = methods.Login(&ctx); err != nil {
		fmt.Println(err)
		return
	}

	var myq []methods.SearchQuery
	if myq, err = methods.GetAllQueries(ctx); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(myq)
}
