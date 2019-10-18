package handlers

import (
	"github.com/button-tech/utils-node-tool/shared/responses"
	"github.com/qiangxue/fasthttp-routing"
	"os"
	"os/exec"
	"strings"
)

var workdir = os.Getenv("WORKDIR")

func GetBalance(c *routing.Context) error {

	stdout, err := exec.Command(workdir+"wrappers/get_balance.py", workdir, c.Param("address")).Output()
	if err != nil {
		return err
	}

	if string(stdout) == "error\n" {
		return routing.NewHTTPError(400, "Bad request")
	}

	balance := strings.TrimSuffix(string(stdout), "\n")

	response := new(responses.BalanceResponse)

	response.Balance = balance

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}
