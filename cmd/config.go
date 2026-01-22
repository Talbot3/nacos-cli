package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github/szpinc/nacosctl/pkg/editor"
	"github/szpinc/nacosctl/pkg/nacos"
	"github/szpinc/nacosctl/pkg/util"
	"os"
	"path/filepath"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

var (
	getAllConfig bool   // 获取所有配置
	fileType     string // 配置类型
)

var getConfig = &cobra.Command{
	Use:   "config",
	Short: "获取 Nacos 配置",
	Long: `从 Nacos 服务器获取配置。

可以指定 dataId 获取单个配置，或使用 --all 参数列出命名空间中的所有配置。`,
	Example: `  # 获取指定配置
  nacosctl get config app.yaml -n public -g DEFAULT_GROUP

  # 列出所有配置
  nacosctl get config -A -n public

  # 保存配置到文件
  nacosctl get config app.yaml -n public > app.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if getAllConfig {
			dataIds, err := nacosClient.AllConfig(nacos.ConfigGetOperation{
				NacosOperation: &nacos.NacosOperation{
					Namespace: namespace,
				},
			})

			if err != nil {
				return err
			}

			printTable(dataIds)
			return nil
		}

		if len(args) == 0 {
			return errors.New("请指定 dataId")
		}

		dataId := args[0]

		configData, err := nacosClient.Get(nacos.ConfigGetOperation{
			NacosOperation: &nacos.NacosOperation{
				Namespace: namespace,
				Group:     group,
			},
			DataId: dataId,
		})

		if err != nil {
			return err
		}

		fmt.Println(configData.Content)
		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		dataIds, err := nacosClient.AllConfig(nacos.ConfigGetOperation{
			NacosOperation: &nacos.NacosOperation{
				Namespace: namespace,
			},
		})
		for _, dataId := range dataIds {
			println(dataId.DataId)
		}
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		names := []string{}

		for _, id := range dataIds {
			names = append(names, id.DataId)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	},
}

var editConfig = &cobra.Command{
	Use:   "config",
	Short: "交互式编辑配置",
	Long: `交互式编辑 Nacos 服务器上的配置。

配置会被下载并在默认编辑器中打开，保存并关闭编辑器后，更改会自动上传到服务器。`,
	Example: `  # 编辑配置
  nacosctl edit config app.yaml -n public -g DEFAULT_GROUP

  # 指定编辑器
  export EDITOR=vim
  nacosctl edit config app.yaml -n public`,
	Run: func(cmd *cobra.Command, args []string) {

		var dataId = args[0]

		configData, err := nacosClient.Get(nacos.ConfigGetOperation{
			NacosOperation: &nacos.NacosOperation{
				Namespace: namespace,
				Group:     group,
			},
			DataId: dataId,
		})

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		e := editor.NewDefaultEditor([]string{})

		buf := &bytes.Buffer{}
		buf.Write([]byte(configData.Content))

		edited, file, err := e.LaunchTempFile(fmt.Sprintf("%s-edit-", filepath.Base(os.Args[0])), configData.Type, buf)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		editedMd5 := util.Md5BytesToString(edited)

		if configData.Md5 == editedMd5 {
			fmt.Println("配置未修改")
			return
		}

		defer func(f string) {
			if e := os.Remove(f); e != nil {
				fmt.Println("删除临时文件错误:", e)
			}
		}(file)

		if fileType == "" {
			fileType = configData.Type
		}

		err = nacosClient.Edit(nacos.ConfigEditOperation{
			NacosOperation: &nacos.NacosOperation{
				Namespace: namespace,
				Group:     group,
			},
			DataId:  dataId,
			Content: string(edited),
			Type:    fileType,
		})

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("配置已更新")
	},
}

var deleteConfig = &cobra.Command{
	Use:   "config",
	Short: "删除 Nacos 配置",
	Long: `删除 Nacos 服务器上的配置。

此操作会永久删除配置，无法撤销。`,
	Example: `  # 删除配置
  nacosctl delete config app.yaml -n public -g DEFAULT_GROUP`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return errors.New("请指定 dataId")
		}

		err := nacosClient.DeleteConfig(nacos.ConfigDeleteOperation{
			NacosOperation: &nacos.NacosOperation{
				Namespace: namespace,
				Group:     group,
			},
			DataId: args[0],
		})
		if err != nil {
			return err
		}

		fmt.Println("配置已删除")
		return nil
	},
}

func init() {

	editConfig.Flags().StringVarP(&fileType, "type", "t", "", "配置文件类型 (如: yaml, properties, json)")

	getConfig.Flags().BoolVarP(&getAllConfig, "all", "A", false, "列出命名空间中的所有配置")

	editCmd.AddCommand(editConfig)
	getCmd.AddCommand(getConfig)
	deleteCmd.AddCommand(deleteConfig)
}

func printTable(items []nacos.NacosPageItem) {
	table := uitable.New()
	table.MaxColWidth = 50

	table.AddRow("DataID", "GROUP", "NAMESPACE")

	for _, item := range items {
		if item.Tenant == "" {
			item.Tenant = "public"
		}
		table.AddRow(item.DataId, item.Group, item.Tenant)
	}

	fmt.Println(table)
}
