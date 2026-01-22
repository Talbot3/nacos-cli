package cmd

import (
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "交互式编辑配置",
	Long: `交互式编辑 Nacos 服务器上的配置。

edit 命令会下载配置，在默认编辑器中打开，
如果进行了修改，会将修改后的版本上传回服务器。

使用的编辑器由 EDITOR 环境变量决定，
Unix 系统默认为 vi，Windows 系统默认为 notepad。`,
	Example: `  # 编辑配置
  nacosctl edit config app.yaml -n public -g DEFAULT_GROUP

  # 设置自定义编辑器
  export EDITOR=vim
  nacosctl edit config app.yaml -n public

  # 通过命令行参数进行认证
  nacosctl edit config app.yaml -n public -u nacos -p nacos

  # 使用环境变量进行认证
  export NACOS_ADDR="http://localhost:8848/nacos"
  export NACOS_USERNAME="nacos"
  export NACOS_PASSWORD="nacos"
  nacosctl edit config app.yaml -n public`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	ValidArgs: []string{"config"},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
