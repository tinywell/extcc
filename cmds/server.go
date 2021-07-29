package cmds

import (
	"extcc/cmds/server"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var seraddr string

func init() {
	InvokeCmd.Flags().StringVar(&seraddr, "listen", ":9999", "服务监听地址")

}

// ServerCmd 链码调用
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "链码仿真服务",
	Long:  "模拟 peer 功能与链码通讯，完成链码调用",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	r := gin.Default()
	r.Use(cors.Default())
	server.Router(r)
	r.Run(seraddr)
}
