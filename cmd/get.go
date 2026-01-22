package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "获取 Nacos 配置",
	Long: `从 Nacos 服务器获取配置。

可以指定 dataId 获取单个配置，或使用 --all 参数列出命名空间中的所有配置。`,
	Example: `  # 获取指定配置
  nacosctl get config app.yaml -n public -g DEFAULT_GROUP

  # 获取配置并保存到文件
  nacosctl get config app.yaml -n public > app.yaml

  # 列出命名空间中的所有配置
  nacosctl get config -A -n public

  # 列出指定分组中的所有配置
  nacosctl get config -A -n public -g PROD_GROUP

  # 使用环境变量进行认证
  export NACOS_ADDR="http://localhost:8848/nacos"
  export NACOS_USERNAME="nacos"
  export NACOS_PASSWORD="nacos"
  nacosctl get config app.yaml -n public`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
