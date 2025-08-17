# FileBrowser 现代化主题配置指南

本指南将帮助您为 FileBrowser 应用现代化的橙色/粉色主题，提供更具现代感的用户界面。

## 🎨 主题特色

- **现代化配色**: 橙色主色调配合粉色辅助色
- **流畅动画**: 使用现代 CSS 动画和过渡效果
- **卡片设计**: 现代化的卡片布局和阴影效果
- **毛玻璃效果**: 头部和侧边栏采用毛玻璃背景
- **优化交互**: 改进的按钮、悬浮效果和焦点状态
- **响应式设计**: 适配移动设备和不同屏幕尺寸

## 📁 配置步骤

### 1. 创建品牌文件夹

```bash
mkdir -p /path/to/filebrowser-branding
```

### 2. 复制主题文件

将 `custom.css` 文件复制到品牌文件夹中：

```bash
cp custom.css /path/to/filebrowser-branding/
```

### 3. 配置 FileBrowser

使用以下命令配置 FileBrowser 使用自定义品牌：

```bash
filebrowser config set --branding.files "/path/to/filebrowser-branding"
```

### 4. 重启 FileBrowser

重启 FileBrowser 服务以应用新主题：

```bash
# 如果使用 systemd
sudo systemctl restart filebrowser

# 或者如果直接运行
# 停止当前进程，然后重新启动
filebrowser
```

## 🐳 Docker 配置

如果您使用 Docker 运行 FileBrowser，需要挂载品牌文件夹：

```bash
docker run -d \
  --name filebrowser \
  -v /path/to/files:/srv \
  -v /path/to/filebrowser-branding:/branding \
  -p 8080:80 \
  filebrowser/filebrowser:latest \
  --branding.files /branding
```

或者在 `docker-compose.yml` 中：

```yaml
version: '3.8'
services:
  filebrowser:
    image: filebrowser/filebrowser:latest
    container_name: filebrowser
    ports:
      - "8080:80"
    volumes:
      - /path/to/files:/srv
      - /path/to/filebrowser-branding:/branding
    command: --branding.files /branding
```

## 🎨 主题自定义

### 颜色变量

您可以修改 `custom.css` 中的颜色变量来自定义主题：

```css
:root {
  --modern-primary: #f47750;     /* 主色调 - 橙色 */
  --modern-secondary: #f8bbd9;   /* 辅助色 - 粉色 */
  --modern-accent: #ec4899;      /* 强调色 - 深粉色 */
}
```

### 添加自定义图标

1. 在品牌文件夹中创建 `img` 目录：
```bash
mkdir -p /path/to/filebrowser-branding/img/icons
```

2. 添加自定义图标文件：
```
img/
├── logo.svg          # 主 logo
└── icons/
    ├── favicon.ico    # 网站图标
    └── favicon.svg    # SVG 图标
```

### 自定义实例名称

```bash
filebrowser config set \
  --branding.name "我的文件管理器" \
  --branding.files "/path/to/filebrowser-branding"
```

## 🛠️ 高级配置

### 禁用外部链接

```bash
filebrowser config set --branding.disableExternal
```

### 禁用磁盘使用量显示

```bash
filebrowser config set --branding.disableUsage
```

### 完整配置示例

```bash
filebrowser config set \
  --branding.name "现代文件管理器" \
  --branding.files "/opt/filebrowser/branding" \
  --branding.disableExternal \
  --address 0.0.0.0 \
  --port 8080 \
  --database /opt/filebrowser/database.db \
  --root /srv
```

## 📱 移动端优化

主题包含了移动端优化：

- 响应式卡片布局
- 触摸友好的按钮尺寸
- 优化的间距和排版
- 减少动画（如果用户偏好）

## 🔧 故障排除

### 主题未生效

1. 检查品牌文件夹路径是否正确
2. 确认 `custom.css` 文件存在且可读
3. 重启 FileBrowser 服务
4. 清除浏览器缓存

### 性能优化

如果遇到性能问题，可以在 `custom.css` 中注释掉某些动画效果：

```css
/* 注释掉复杂动画 */
/*
#listing .item:hover i {
  transform: scale(1.1) rotate(5deg) !important;
}
*/
```

### 兼容性

- 支持现代浏览器（Chrome 80+, Firefox 75+, Safari 13+）
- 包含 `-webkit-` 前缀以支持 Safari
- 提供降级方案用于不支持的浏览器

## 📋 文件结构

```
/path/to/filebrowser-branding/
├── custom.css              # 主题样式文件
└── img/                     # 可选：自定义图标
    ├── logo.svg
    └── icons/
        ├── favicon.ico
        └── favicon.svg
```

## 🎯 更新主题

要更新主题，只需替换 `custom.css` 文件并重启 FileBrowser：

```bash
cp new-custom.css /path/to/filebrowser-branding/custom.css
sudo systemctl restart filebrowser
```

## 💡 提示

- 主题使用 CSS 变量，易于自定义
- 所有样式都使用 `!important` 确保优先级
- 包含深色模式适配
- 支持减少动画偏好设置
- 使用现代 CSS 特性如 `backdrop-filter` 和 `cubic-bezier`

享受您的现代化 FileBrowser 体验！🚀