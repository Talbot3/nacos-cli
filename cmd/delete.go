package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除 Nacos 配置",
	Long: `删除 Nacos 服务器上的配置。

delete 命令会从 Nacos 服务器删除指定的配置。
此操作无法撤销。`,
	Example: `  # 删除配置
  nacosctl delete config app.yaml -n public -g DEFAULT_GROUP

  # 通过命令行参数进行认证
  nacosctl delete config app.yaml -n public -u nacos -p nacos

  # 使用环境变量进行认证
  export NACOS_ADDR="http://localhost:8848/nacos"
  export NACOS_USERNAME="nacos"
  export NACOS_PASSWORD="nacos"
  nacosctl delete config app.yaml -n public

  # 删除指定分组中的配置
  nacosctl delete config app.yaml -n public -g PROD_GROUP`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
