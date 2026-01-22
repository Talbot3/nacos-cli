# nacosctl

`nacosctl` 是一个用于管理 Nacos 配置的命令行工具。

## 功能特性

- 支持用户名密码认证，Token 自动缓存和刷新
- 获取、创建、编辑、删除配置
- 列出命名空间中的所有配置
- 兼容无认证模式的 Nacos 服务器
- 基于 [Nacos Open API](https://nacos.io/zh-cn/docs/open-api.html) 实现

## 安装

### Linux

**AMD64**

```bash
curl -L -o /usr/local/bin/nacosctl https://github.com/szpinc/nacos-cli/releases/download/v1.2/nacosctl_linux_amd64
chmod +x /usr/local/bin/nacosctl
```

**ARM64**

```bash
curl -L -o /usr/local/bin/nacosctl https://github.com/szpinc/nacos-cli/releases/download/v1.2/nacosctl_linux_arm64
chmod +x /usr/local/bin/nacosctl
```

### macOS

**AMD64**

```bash
curl -L -o /usr/local/bin/nacosctl https://github.com/szpinc/nacos-cli/releases/download/v1.2/nacosctl_darwin_amd64
chmod +x /usr/local/bin/nacosctl
```

**ARM64**

```bash
curl -L -o /usr/local/bin/nacosctl https://github.com/szpinc/nacos-cli/releases/download/v1.2/nacosctl_darwin_arm64
chmod +x /usr/local/bin/nacosctl
```

### 从源码编译

```bash
go build -o nacosctl .
```

## 配置

### 环境变量

| 变量 | 说明 | 示例 |
|------|------|------|
| `NACOS_ADDR` | Nacos 服务器地址 | `http://localhost:8848/nacos` |
| `NACOS_USERNAME` | 用户名 | `nacos` |
| `NACOS_PASSWORD` | 密码 | `your-password` |
| `NACOS_API_VERSION` | API 版本 | `v1` (默认) |

### 命令行参数

| 参数 | 短参数 | 说明 |
|------|--------|------|
| `--namespace` | `-n` | Nacos 命名空间 ID (必填) |
| `--group` | `-g` | 分组名称 (默认: DEFAULT_GROUP) |
| `--username` | `-u` | 用户名 (覆盖环境变量) |
| `--password` | `-p` | 密码 (覆盖环境变量) |

## 使用示例

### 设置环境变量

```bash
export NACOS_ADDR="http://localhost:8848/nacos"
export NACOS_USERNAME="nacos"
export NACOS_PASSWORD="your-password"
```

### 获取配置

```bash
# 获取指定配置
nacosctl get config app.yaml -n public -g DEFAULT_GROUP

# 保存配置到文件
nacosctl get config app.yaml -n public > app.yaml

# 列出所有配置
nacosctl get config -A -n public
```

### 创建/更新配置

```bash
# 从文件创建或更新配置
nacosctl apply config --file ./app.yaml -n public

# 指定自定义 dataId
nacosctl apply config --file ./config.yaml --id app-config -n public

# 显式指定文件类型
nacosctl apply config --file ./app.conf --type properties -n public
```

### 编辑配置

```bash
# 交互式编辑配置
nacosctl edit config app.yaml -n public

# 使用自定义编辑器
export EDITOR=vim
nacosctl edit config app.yaml -n public
```

### 删除配置

```bash
nacosctl delete config app.yaml -n public
```

## 认证说明

工具支持两种模式：

1. **认证模式**: 配置 `NACOS_USERNAME` 和 `NACOS_PASSWORD` 后，自动登录并缓存 Token
2. **无认证模式**: 不配置用户名密码，直接访问 Nacos 服务器

Token 缓存在 `~/.nacosctl/token_*.json`，有效期 5 小时，过期前自动刷新。

### Nacos 2.4.0+ 管理员密码初始化

从 Nacos 2.4.0 版本开始，**已取消默认密码**。首次启用认证后，需要通过 API 初始化管理员用户 `nacos` 的密码：

```bash
# 初始化管理员密码
curl -X POST 'http://localhost:8849/nacos/v1/auth/users/admin' -d 'password=your-password'

# 响应示例
{"username":"nacos","password":"your-password"}
```

> **注意**：
> - 用户名固定为 `nacos`，不可修改
> - 若不指定 `password` 参数或传空字符串，Nacos 将生成随机密码
> - 该 API 只能在首次初始化时调用一次，成功后将无法再次调用

### Docker 环境初始化步骤

如果使用 Docker Compose 启动带认证的 Nacos，首次启动后请执行：

```bash
# 1. 等待 Nacos 启动完成
docker-compose up -d
sleep 30

# 2. 初始化管理员密码
curl -X POST 'http://localhost:8849/nacos/v1/auth/users/admin' -d 'password=nacos'

# 3. 验证登录
curl -X POST 'http://localhost:8849/nacos/v1/auth/login' -d 'username=nacos&password=nacos'
```

## 常见问题

**Q: 如何查看帮助信息？**

```bash
nacosctl -h              # 查看主帮助
nacosctl get -h          # 查看子命令帮助
nacosctl get config -h   # 查看具体命令帮助
```

**Q: 提示认证失败怎么办？**

检查 `NACOS_ADDR`、`NACOS_USERNAME`、`NACOS_PASSWORD` 是否正确配置。

**Q: Token 缓存在哪里？**

`~/.nacosctl/` 目录下，每个服务器地址对应一个缓存文件。

## 测试环境

项目提供了 Docker Compose 配置，可同时启动两个 Nacos 实例用于测试：

```bash
docker-compose up -d
```

### 服务说明

| 服务 | 地址 | 认证 | 控制台 |
|------|------|------|--------|
| nacos-no-auth | http://localhost:8848/nacos | 无需认证 | 直接访问 |
| nacos-with-auth | http://localhost:8849/nacos | nacos/nacos | 需要登录 |

### 测试无认证模式

```bash
export NACOS_ADDR="http://localhost:8848/nacos"

# 创建测试配置
echo "app.name: test" > test.yaml
nacosctl apply config --file test.yaml -n public

# 获取配置
nacosctl get config test.yaml -n public
```

### 测试带认证模式

> **重要**: Nacos 2.4.0+ 首次启动后需要先初始化管理员密码

```bash
# 1. 启动服务
docker-compose up -d
sleep 30

# 2. 初始化管理员密码（仅首次需要）
curl -X POST 'http://localhost:8849/nacos/v1/auth/users/admin' -d 'password=nacos'

# 3. 测试 nacosctl
export NACOS_ADDR="http://localhost:8849/nacos"
export NACOS_USERNAME="nacos"
export NACOS_PASSWORD="nacos"

# 创建测试配置
echo "app.name: test" > test.yaml
nacosctl apply config --file test.yaml -n public

# 获取配置
nacosctl get config test.yaml -n public

# 列出所有配置
nacosctl get config -A -n public
```

### 停止测试环境

```bash
docker-compose down
```

## License

MIT
