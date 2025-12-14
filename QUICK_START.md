# UniFocus 快速启动指南

## 🎯 项目状态

**后端完成度: 70%** | **前端完成度: 30%**

---

## 🚀 快速启动

### 1. 启动后端服务

```bash
cd unifocus/backend

# 安装依赖
go mod tidy

# 启动服务（默认端口8080）
go run cmd/api/main.go
```

后端服务将在 `http://localhost:8080` 启动

### 2. 启动NLP服务

```bash
cd unifocus/nlp-service

# 安装依赖
pip install -r requirements.txt

# 启动服务（默认端口8000）
python -m app.main
```

### 3. 启动前端监控面板

```bash
cd unifocus/web

# 安装依赖
npm install

# 启动开发服务器（默认端口3000）
npm run dev
```

### 4. 访问监控面板

打开浏览器访问: **http://localhost:3000/dashboard**

---

## 📊 监控面板功能

监控面板提供以下实时信息：

1. **系统状态** - 后端服务健康状态
2. **运行时间** - 服务运行时长
3. **机会总数** - 当前数据库中的机会数量
4. **API状态** - API服务可用性
5. **最近机会列表** - 显示最新的机会数据
6. **开发进度** - 已完成和待开发功能列表

**自动刷新**: 每30秒自动刷新数据

---

## 🔌 API端点

### 健康检查
```bash
curl http://localhost:8080/health
```

### 获取系统指标
```bash
curl http://localhost:8080/api/v1/metrics
```

### 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "school": "清华大学",
    "major": "计算机科学",
    "grade": 3
  }'
```

### 获取机会列表
```bash
curl http://localhost:8080/api/v1/opportunities?limit=10
```

---

## 📁 项目结构

```
unifocus/
├── backend/              # Go后端 (端口8080)
│   ├── cmd/api/         # API服务入口
│   ├── internal/
│   │   ├── api/         # HTTP handlers
│   │   ├── service/     # 业务逻辑
│   │   ├── repository/  # 数据访问
│   │   └── crawler/     # 爬虫模块
│   └── configs/         # 配置文件
│
├── nlp-service/         # Python NLP服务 (端口8000)
│   └── app/
│       ├── services/    # NLP服务
│       └── api/routes/  # API路由
│
└── web/                 # Next.js前端 (端口3000)
    ├── app/             # 页面
    ├── lib/api/         # API客户端
    └── package.json
```

---

## ✅ 已完成功能

- ✅ 数据库连接池和Redis客户端
- ✅ 用户认证系统（注册/登录/JWT）
- ✅ 机会管理服务（CRUD + 查询）
- ✅ 基础爬虫框架
- ✅ NLP文本提取服务
- ✅ 用户画像服务
- ✅ API监控系统

---

## ⏳ 待开发功能

- ⏳ 双维度评分引擎
- ⏳ 竞赛级别识别服务
- ⏳ 推送服务和定时匹配任务
- ⏳ 前端完整页面（登录、机会列表等）

---

## 🐛 故障排查

### 后端无法启动
1. 检查PostgreSQL是否运行
2. 检查Redis是否运行
3. 检查配置文件 `configs/config.dev.yaml`

### 前端无法连接后端
1. 确认后端服务在8080端口运行
2. 检查 `web/next.config.js` 中的API_BASE_URL配置
3. 检查CORS设置

### NLP服务无法启动
1. 检查Python版本（需要3.11+）
2. 确认所有依赖已安装
3. 检查端口8000是否被占用

---

## 📝 下一步

1. 完善前端页面（登录、注册、机会列表）
2. 实现评分引擎
3. 实现推送服务
4. 添加更多爬虫站点

---

**最后更新**: 2024年


