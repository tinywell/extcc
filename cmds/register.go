package cmds

import (
	"encoding/json"
	"extcc/cmds/server"
	"fmt"

	"github.com/spf13/cobra"
)

// RegisterURL 链码注册地址
const RegisterURL = server.URLRegister

var ccaddr string
var ccname string

// RegCmd 链码调用
var RegCmd = &cobra.Command{
	Use:   "register",
	Short: "链码注册",
	Run: func(cmd *cobra.Command, args []string) {
		register()
	},
}

func init() {
	RegCmd.Flags().StringVar(&address, "addr", "localhost:9999", "仿真服务地址")
	RegCmd.Flags().StringVar(&ccaddr, "ccaddr", "localhost:9998", "外部链码服务地址")
	RegCmd.Flags().StringVar(&ccname, "ccname", "mycc", "外部链码名称")
}

func register() {
	req := &server.ReqChaincode{
		Name: ccname,
		Addr: ccaddr,
	}
	body, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("http://%s%s", address, RegisterURL)
	rsp, err := Post(url, body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rsp))
}
