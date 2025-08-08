# 部署说明 - File Browser My Branch

这个分支包含了多项性能优化和功能改进，使用专业化的CI/CD流程进行自动构建和部署。

## 🚀 专业化CI/CD流程

基于官方工作流设计，当向 `my` 分支推送代码时，会自动触发完整的质量检查和构建流程：

### 📋 质量保证阶段
1. **代码检查**: 前端 ESLint + 后端 go fmt/vet
2. **类型检查**: TypeScript 编译验证
3. **单元测试**: 前端和后端测试套件
4. **并行执行**: 所有检查同时进行以节省时间

### 🔨 构建和发布阶段
1. **前端构建**: pnpm 优化构建，资源压缩
2. **多平台编译**: Linux/macOS/Windows (amd64/arm64)
3. **Docker镜像**: 多架构Alpine镜像 (amd64/arm64)
4. **自动发布**: GitHub Release + DockerHub

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
- ✅ **质量保证**: 代码检查、类型检查、单元测试
- ✅ **并行执行**: 前端和后端任务同时进行
- ✅ **多平台编译**: Linux/macOS/Windows 全平台支持
- ✅ **Docker镜像**: 多架构 Alpine 镜像自动构建
- ✅ **自动发布**: GitHub Release 和 DockerHub 同步发布
- ✅ **健康检查**: Docker 镜像包含健康检查端点
- ✅ **语义化版本**: 自动生成版本标签和发布说明

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