package cmd

import (
	"github/szpinc/nacosctl/pkg/nacos"
	"os"

	"github.com/spf13/cobra"
)

var namespace string
var group string
var username string
var password string

var nacosClient *nacos.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nacosctl",
	Short: "Nacos 配置管理命令行工具",
	Long: `nacosctl 是一个用于管理 Nacos 配置的命令行工具。

支持读取、创建、更新和删除 Nacos 服务器上的配置。
可以通过用户名密码或环境变量进行身份验证。`,
	Example: `  # 通过环境变量设置 Nacos 服务器地址
  export NACOS_ADDR="http://localhost:8848/nacos"

  # 通过环境变量设置认证信息
  export NACOS_USERNAME="nacos"
  export NACOS_PASSWORD="nacos"

  # 或者通过命令行参数传递认证信息
  nacosctl get config app.yaml -n public -u nacos -p nacos

  # 获取单个配置
  nacosctl get config app.yaml -n public -g DEFAULT_GROUP

  # 列出命名空间中的所有配置
  nacosctl get config -A -n public

  # 应用配置文件
  nacosctl apply config --file ./app.yaml -n public -g DEFAULT_GROUP

  # 交互式编辑配置
  nacosctl edit config app.yaml -n public

  # 删除配置
  nacosctl delete config app.yaml -n public`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Nacos 命名空间 ID (必填)")
	rootCmd.PersistentFlags().StringVarP(&group, "group", "g", "DEFAULT_GROUP", "Nacos 分组名称")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Nacos 用户名 (覆盖 NACOS_USERNAME 环境变量)")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Nacos 密码 (覆盖 NACOS_PASSWORD 环境变量)")

	_ = rootCmd.MarkFlagRequired("namespace")

	// 初始化客户端
	nacosClient = initNacosClient()
}

// initNacosClient 初始化Nacos客户端
func initNacosClient() *nacos.Client {
	addr := os.Getenv("NACOS_ADDR")
	apiVersion := os.Getenv("NACOS_API_VERSION")
	envUsername := os.Getenv("NACOS_USERNAME")
	envPassword := os.Getenv("NACOS_PASSWORD")

	// 命令行参数优先级高于环境变量
	if username != "" {
		envUsername = username
	}
	if password != "" {
		envPassword = password
	}

	// 如果命令行或环境变量指定了用户名密码，使用NewClient
	if envUsername != "" && envPassword != "" {
		return nacos.NewClient(addr, apiVersion, envUsername, envPassword)
	}

	return nacos.NewDefaultClient()
}
