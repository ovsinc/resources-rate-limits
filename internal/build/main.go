package main

import (
	"fmt"

	ratelimits "github.com/ovsinc/resources-rate-limits"
	rescommon "github.com/ovsinc/resources-rate-limits/pkg/resources/common"
)

func main() {
	res := ratelimits.Check()

	var out string

	switch res.Type() {
	case rescommon.ResourceType_CG1:
		out = "cg1"

	case rescommon.ResourceType_CG2:
		out = "cg2"

	case rescommon.ResourceType_OS:
		out = "os"
	}

	fmt.Printf("%s", out)
}
