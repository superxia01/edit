# Edit Business Chrome 插件安装说明

## 📦 安装步骤

### 方法一：直接安装（推荐）

1. **下载插件包**
   - 插件文件：`edit-business-chrome-plugin.zip`

2. **解压插件**
   ```bash
   unzip edit-business-chrome-plugin.zip
   ```

3. **在 Chrome 中加载插件**
   - 打开 Chrome 浏览器
   - 访问：`chrome://extensions/`
   - 打开右上角"开发者模式"
   - 点击"加载已解压的扩展程序"
   - 选择 `chrome-plugin` 文件夹

4. **验证安装**
   - 浏览器工具栏会出现插件图标
   - 访问小红书页面，点击图标测试

### 方法二：从源码安装

```bash
cd chrome-plugin
# 在 Chrome 中加载此文件夹
```

---

## ⚙️ 配置说明

### API 配置

插件已配置生产环境地址：`https://edit.crazyaigc.com`

如需切换到开发环境，编辑 `api-config.js`：
```javascript
// 取消注释开发环境
BASE_URL: 'http://localhost:8084',
```

---

## 🚀 使用方法

### 1. 单篇笔记采集
- 访问小红书笔记详情页
- 点击插件图标
- 点击"采集笔记"按钮

### 2. 批量采集博主笔记
- 访问小红书博主主页
- 点击插件图标
- 点击"批量采集笔记"按钮

### 3. 采集博主信息
- 访问小红书博主主页
- 点击插件图标
- 点击"采集博主信息"按钮

---

## 🔗 相关链接

- **Web 管理后台**：https://edit.crazyaigc.com
- **API 文档**：https://edit.crazyaigc.com/api/v1/docs（如有）

---

## ⚠️ 注意事项

1. **需要先登录**：使用微信登录系统后才能采集数据
2. **网络要求**：确保能访问 edit.crazyaigc.com
3. **采集限制**：受用户设置中的采集开关和每日限制控制

---

## 📋 版本信息

- **版本**：2.0.0
- **更新日期**：2026-02-05
- **适用系统**：Edit Business (小红书数据采集系统)
