package cmd

import (
	"github/szpinc/nacosctl/pkg/nacos"

	"github.com/spf13/cobra"
)

var (
	file   string
	dataId string
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "创建或更新配置",
	Long: `在 Nacos 服务器上创建或更新配置。

apply 命令会创建新配置或更新现有配置。
默认情况下，dataId 从文件名派生，但可以通过 --id 参数覆盖。
文件类型会从文件扩展名自动检测。`,
	Example: `  # 使用文件创建或更新配置
  nacosctl apply config --file ./app.yaml -n public -g DEFAULT_GROUP

  # 指定自定义 dataId
  nacosctl apply config --file ./config.yaml --id app-config -n public

  # 显式指定文件类型
  nacosctl apply config --file ./app.conf --type properties -n public

  # 使用环境变量
  export NACOS_ADDR="http://localhost:8848/nacos"
  export NACOS_USERNAME="nacos"
  export NACOS_PASSWORD="nacos"
  nacosctl apply config --file ./app.yaml -n public

  # 使用自定义分组
  nacosctl apply config --file ./app.yaml -n public -g PROD_GROUP`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nacosClient.ApplyConfig(nacos.ConfigApplyOperation{
			NacosOperation: &nacos.NacosOperation{
				Namespace: namespace,
				Group:     group,
			},
			DataId: dataId,
			File:   file,
			Type:   fileType,
		})
	},
}

func init() {
	applyCmd.Flags().StringVarP(&file, "file", "f", "", "配置文件路径 (必填)")
	applyCmd.Flags().StringVarP(&dataId, "id", "d", "", "自定义 dataId (默认为文件名)")
	applyCmd.Flags().StringVarP(&fileType, "type", "t", "", "配置文件类型 (如: yaml, properties, json)。默认从文件扩展名自动检测")

	applyCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(applyCmd)
}
