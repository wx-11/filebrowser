# 部署说明 - File Browser My Branch

这个分支包含了多项性能优化和功能改进，提供自动化的构建和部署流程。

## 🚀 自动构建流程

当向 `my` 分支推送代码时，会自动触发构建流程：

1. **前端构建**: 使用 pnpm 构建 Vue.js 前端应用
2. **后端构建**: 为多个平台构建 Go 二进制文件
3. **发布到 GitHub**: 创建带有所有平台二进制文件的 Release
4. **发布到 DockerHub**: 构建并推送优化的 Alpine 镜像

## 📦 支持的平台

### 二进制文件
- Linux (amd64, arm64)
- macOS (amd64, arm64) 
- Windows (amd64, arm64)
- FreeBSD (amd64)

### Docker 镜像
- `linux/amd64`
- `linux/arm64`

## 🐳 Docker 部署

### 快速启动
```bash
docker run -d \
  --name filebrowser \
  -p 8080:8080 \
  -v /path/to/data:/srv \
  ${DOCKER_USERNAME}/filebrowser:latest
```

### 使用 Docker Compose
1. 复制示例文件:
   ```bash
   cp docker-compose.example.yml docker-compose.yml
   ```

2. 编辑 `docker-compose.yml` 修改:
   - 将 `${DOCKER_USERNAME}` 替换为实际的 Docker Hub 用户名
   - 修改数据路径 `/path/to/your/data`
   - 根据需要调整端口和环境变量

3. 启动服务:
   ```bash
   docker-compose up -d
   ```

### 环境变量配置
| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| FB_DATABASE | 数据库文件路径 | `./filebrowser.db` |
| FB_CONFIG | 配置文件路径 | - |
| FB_ROOT | 文件根目录 | `.` |
| FB_PORT | 监听端口 | `8080` |
| FB_ADDRESS | 监听地址 | `127.0.0.1` |

## 🔧 GitHub Actions 配置

### 必需的变量和密钥

在你的 GitHub 仓库设置中添加：

#### Variables (Settings → Secrets and variables → Actions → Variables)
- `DOCKER_USERNAME`: Docker Hub 用户名

#### Secrets (Settings → Secrets and variables → Actions → Secrets)  
- `DOCKER_PASSWORD`: Docker Hub 密码或访问令牌

### 工作流特性
- ✅ 自动前端构建和优化
- ✅ 多平台二进制文件编译
- ✅ 自动创建 GitHub Release
- ✅ 多架构 Docker 镜像构建
- ✅ 自动清理临时文件
- ✅ 详细的发布说明生成

## 📋 功能特性

### 🚀 性能优化
- WebSocket 缓冲区从 1KB 增加到 16KB
- 使用 goroutine 并发读取 stdout/stderr  
- 8KB I/O 缓冲区提升文件操作效率
- 30分钟命令执行超时控制

### 🖱️ 用户界面改进
- 修复了 Shell 界面文本选择功能
- 优化了滚动性能使用 requestAnimationFrame
- 改进了光标行为和视觉反馈
- 移除了不必要的换行符处理

### 🔗 功能增强
- 恢复了分享文件的直接下载链接
- 增加了 CLI 用户管理标志
- 移除了限制性的命令白名单检查

## 🔨 本地构建

### 前端构建
```bash
cd frontend
pnpm install
pnpm run build
```

### 后端构建
```bash
go mod download
go build -ldflags "-s -w" -o filebrowser .
```

### Docker 构建
```bash
docker build -f Dockerfile.alpine -t filebrowser:my .
```

## 📊 版本标签

Docker 镜像标签:
- `latest`: 最新的 my 分支构建
- `my`: my 分支专用标签
- `{commit-hash}`: 特定提交的构建版本

## 🆘 故障排除

### 常见问题

1. **命令执行缓慢**
   - 新版本已优化，如仍有问题检查网络和系统资源

2. **无法选择文本**  
   - 已修复，确保使用最新版本

3. **Docker 镜像拉取失败**
   - 检查 Docker Hub 用户名是否正确
   - 确认镜像标签存在

4. **构建失败**
   - 检查 GitHub Secrets 配置
   - 查看 Actions 日志获取详细错误信息

## 🤝 贡献

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)  
5. 创建 Pull Request 到 `my` 分支

## 📝 许可证

本项目遵循原项目的许可证条款。