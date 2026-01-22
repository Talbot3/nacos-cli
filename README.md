# nacosctl

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Nacos Version](https://img.shields.io/badge/Nacos-2.4.0+-2E86C1?logo=apache-nifi&logoColor=white)](https://nacos.io/)
[![Tests](https://github.com/Talbot3/nacos-cli/actions/workflows/test.yml/badge.svg)](https://github.com/Talbot3/nacos-cli/actions/workflows/test.yml)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

`nacosctl` 是一个用于管理 Nacos 配置的命令行工具，提供简洁高效的配置管理体验。

## 为什么选择 nacosctl？

- **简单易用** - 直观的命令语法，无需学习复杂的 API
- **自动化认证** - Token 自动缓存和刷新，无需手动管理会话
- **多环境支持** - 轻松管理开发、测试、生产等多套环境配置
- **灵活集成** - 可轻松集成到 CI/CD 流程和自动化脚本中
- **安全可靠** - 基于 [Nacos Open API](https://nacos.io/zh-cn/docs/open-api.html) 实现，支持认证和权限控制

## 功能特性

- **配置管理** - 创建、更新、删除、查询配置
- **列表查看** - 列出命名空间中的所有配置
- **交互式编辑** - 直接在终端编辑远程配置
- **多种格式** - 支持 YAML、JSON、Properties、TXT 等格式
- **认证支持** - 支持用户名密码认证，Token 自动缓存和刷新
- **多命名空间** - 支持不同命名空间和分组管理
- **兼容性强** - 兼容无认证模式的 Nacos 服务器

## 快速开始

### 安装

#### Linux

**AMD64**

```bash
curl -L -o /usr/local/bin/nacosctl https://github.com/Talbot3/nacos-cli/releases/download/v1.2/nacosctl_linux_amd64
chmod +x /usr/local/bin/nacosctl
```

**ARM64**

```bash
curl -L -o /usr/local/bin/nacosctl https://github.com/Talbot3/nacos-cli/releases/download/v1.2/nacosctl_linux_arm64
chmod +x /usr/local/bin/nacosctl
```

#### macOS

**AMD64**

```bash
curl -L -o /usr/local/bin/nacosctl https://github.com/Talbot3/nacos-cli/releases/download/v1.2/nacosctl_darwin_amd64
chmod +x /usr/local/bin/nacosctl
```

**ARM64 (Apple Silicon)**

```bash
curl -L -o /usr/local/bin/nacosctl https://github.com/Talbot3/nacos-cli/releases/download/v1.2/nacosctl_darwin_arm64
chmod +x /usr/local/bin/nacosctl
```

#### 从源码编译

```bash
go build -o nacosctl .
```

### 配置

#### 环境变量

| 变量 | 说明 | 示例 | 必填 |
|------|------|------|------|
| `NACOS_ADDR` | Nacos 服务器地址 | `http://localhost:8848/nacos` | 是 |
| `NACOS_USERNAME` | 用户名 | `nacos` | 否* |
| `NACOS_PASSWORD` | 密码 | `your-password` | 否* |
| `NACOS_API_VERSION` | API 版本 | `v1` (默认) | 否 |

*当 Nacos 启用认证时必填

#### 命令行参数

| 参数 | 短参数 | 说明 |
|------|--------|------|
| `--namespace` | `-n` | Nacos 命名空间 ID (如: public) |
| `--group` | `-g` | 分组名称 (默认: DEFAULT_GROUP) |
| `--username` | `-u` | 用户名 (覆盖环境变量) |
| `--password` | `-p` | 密码 (覆盖环境变量) |

## 使用场景

### 场景一：日常开发配置管理

快速查看和修改本地开发环境的配置：

```bash
# 设置开发环境地址
export NACOS_ADDR="http://localhost:8848/nacos"
export NACOS_USERNAME="nacos"
export NACOS_PASSWORD="nacos"

# 查看应用配置
nacosctl get config application.yaml -n public

# 列出所有配置
nacosctl get config -A -n public

# 交互式编辑配置
nacosctl edit config application.yaml -n public
```

### 场景二：批量导入配置

将本地配置文件批量导入到 Nacos：

```bash
# 导入应用主配置
nacosctl apply config --file ./application.yaml -n public -g DEFAULT_GROUP

# 导入数据库配置
nacosctl apply config --file ./database.yaml -n public -g DEFAULT_GROUP

# 导入 Redis 配置
nacosctl apply config --file ./redis.yaml -n public -g DEFAULT_GROUP

# 导入到生产环境分组
nacosctl apply config --file ./application.yaml -n public -g PROD_GROUP
```

### 场景三：多环境部署

管理不同环境的配置，确保配置一致性：

```bash
# 开发环境
export NACOS_ADDR="http://dev-nacos:8848/nacos"
nacosctl apply config --file ./app-config.yaml -n public -g DEV_GROUP

# 测试环境
export NACOS_ADDR="http://test-nacos:8848/nacos"
nacosctl apply config --file ./app-config.yaml -n public -g TEST_GROUP

# 生产环境
export NACOS_ADDR="http://prod-nacos:8848/nacos"
export NACOS_PASSWORD="${PROD_NACOS_PASSWORD}"
nacosctl apply config --file ./app-config.yaml -n public -g PROD_GROUP
```

### 场景四：CI/CD 集成

在 CI/CD 流程中自动更新配置：

```bash
#!/bin/bash
# deploy.sh

# 配置 Nacos 连接信息
export NACOS_ADDR="${NACOS_URL}"
export NACOS_USERNAME="${NACOS_USER}"
export NACOS_PASSWORD="${NACOS_PASS}"

# 获取当前 Git 分支作为配置标识
BRANCH=$(git rev-parse --abbrev-ref HEAD)
CONFIG_ID="app-${BRANCH}"

# 发布配置
nacosctl apply config \
  --file "./config/application.yaml" \
  --id "${CONFIG_ID}" \
  -n public \
  -g DEFAULT_GROUP

echo "配置已发布: ${CONFIG_ID}"
```

### 场景五：配置迁移

将配置从一个 Nacos 集群迁移到另一个：

```bash
#!/bin/bash
# migrate.sh

# 从源集群导出所有配置
export NACOS_ADDR="http://old-nacos:8848/nacos"
export NACOS_USERNAME="nacos"
export NACOS_PASSWORD="nacos"

# 获取所有配置列表
nacosctl get config -A -n public | awk '{print $1}' | while read dataId; do
  # 下载配置内容
  nacosctl get config "$dataId" -n public > "backup/$dataId"
done

# 导入到新集群
export NACOS_ADDR="http://new-nacos:8848/nacos"
for file in backup/*; do
  nacosctl apply config --file "$file" -n public
done
```

### 场景六：配置备份与恢复

定期备份重要配置：

```bash
#!/bin/bash
# backup.sh
BACKUP_DIR="./backups/$(date +%Y%m%d)"
mkdir -p "$BACKUP_DIR"

export NACOS_ADDR="http://nacos:8848/nacos"
export NACOS_USERNAME="nacos"
export NACOS_PASSWORD="nacos"

# 备份所有命名空间的配置
for ns in public dev test; do
  mkdir -p "$BACKUP_DIR/$ns"
  nacosctl get config -A -n "$ns" | awk '{print $1}' | while read dataId; do
    nacosctl get config "$dataId" -n "$ns" > "$BACKUP_DIR/$ns/$dataId"
  done
done

echo "备份完成: $BACKUP_DIR"
```

### 场景七：使用不同文件类型

nacosctl 自动识别文件类型，也支持手动指定：

```bash
# YAML 类型 (自动识别)
nacosctl apply config --file ./app.yaml -n public

# JSON 类型
echo '{"database": {"host": "localhost"}}' > db.json
nacosctl apply config --file ./db.json -n public

# Properties 类型
echo "database.host=localhost\ndatabase.port=3306" > db.properties
nacosctl apply config --file ./db.properties -n public

# 手动指定类型
nacosctl apply config --file ./app.conf --type properties -n public
```

### 场景八：自定义 dataId

使用与文件名不同的 dataId：

```bash
# 文件名 local.yaml，但 dataId 为 application-prod.yaml
nacosctl apply config \
  --file ./local.yaml \
  --id application-prod.yaml \
  -n public \
  -g PROD_GROUP
```

## 认证说明

### 认证模式

工具支持两种模式：

1. **认证模式**: 配置 `NACOS_USERNAME` 和 `NACOS_PASSWORD` 后，自动登录并缓存 Token
2. **无认证模式**: 不配置用户名密码，直接访问 Nacos 服务器

Token 缓存在 `~/.nacosctl/token_*.json`，基于服务器地址和用户名分别缓存，有效期 5 小时，过期前自动刷新。

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

## 常见问题

**Q: 如何查看帮助信息？**

```bash
nacosctl -h              # 查看主帮助
nacosctl get -h          # 查看子命令帮助
nacosctl get config -h   # 查看具体命令帮助
```

**Q: 提示认证失败怎么办？**

检查 `NACOS_ADDR`、`NACOS_USERNAME`、`NACOS_PASSWORD` 是否正确配置。如果是 Nacos 2.4.0+，确保已初始化管理员密码。

**Q: Token 缓存在哪里？**

`~/.nacosctl/` 目录下，每个服务器地址和用户名组合对应一个缓存文件。

**Q: 如何切换用户？**

修改 `NACOS_USERNAME` 和 `NACOS_PASSWORD` 环境变量，工具会自动使用新用户登录。

**Q: 支持哪些配置文件格式？**

支持 YAML、JSON、Properties、TXT 等常见格式。工具会根据文件扩展名自动识别。

## License

MIT
