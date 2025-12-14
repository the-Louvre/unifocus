# UniFocus 高校机会助手

<div align="center">

**一站式高校机会发现与管理平台**

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Python Version](https://img.shields.io/badge/Python-3.11+-3776AB?logo=python)](https://www.python.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

[技术蓝图](readme/technical_blueprint.md) • [开发文档](docs/architecture.md) • [API文档](docs/api/README.md)

</div>

---

## 📖 项目简介

UniFocus 是一个专为大学生设计的智能机会发现与管理平台，帮助学生：

- 🔍 **自动发现机会**：从官网、公众号等多渠道自动采集竞赛、实习、奖学金等信息
- 🎯 **智能匹配推荐**：基于用户画像的个性化评分与推荐
- 📊 **可视化管理**：日历视图、DDL提醒、标签宇宙等多维度展示
- 🏆 **竞赛级别认定**：自动识别国家级A类/B类、省级、校级竞赛
- 💡 **能力缺口诊断**：对比简历与机会要求，提供成长建议

---

## 🏗️ 技术架构

### 技术栈

**后端服务 (Go)**
- 框架: Gin
- 数据库: PostgreSQL 15+
- 缓存: Redis 7+
- 爬虫: Colly + Playwright
- 任务队列: Asynq

**NLP/AI 服务 (Python)**
- 框架: FastAPI
- NLP: SpaCy, Sentence-Transformers
- OCR: PaddleOCR, Tesseract
- 向量化: Sentence-BERT

**前端 (待开发)**
- 浏览器插件: React + TypeScript + Vite
- Web后台: Next.js 14

**基础设施**
- 容器化: Docker + Docker Compose
- CI/CD: GitHub Actions
- 监控: Prometheus + Grafana (计划中)

### 系统架构图

```
┌─────────────────────────────────────────────────────────────┐
│                       客户端层                               │
│  浏览器插件 (React)  │  Web管理后台 (Next.js)               │
└────────────────────────┬────────────────────────────────────┘
                         │ REST API
┌────────────────────────┼────────────────────────────────────┐
│                    后端服务层                                │
│  API服务(Go) │ 爬虫服务(Go) │ NLP服务(Python)               │
└────────────────────────┼────────────────────────────────────┘
                         │
┌────────────────────────┼────────────────────────────────────┐
│                    数据层                                    │
│  PostgreSQL  │  Redis  │  Elasticsearch (可选)              │
└─────────────────────────────────────────────────────────────┘
```

---

## 🚀 快速开始

### 环境要求

- Docker 20.10+（或 Docker Desktop）
- Docker Compose 2.0+
- Go 1.21+ (本地开发)
- Python 3.11+ (本地开发)
- Node.js 18+ (前端开发)

### 安装 Docker (如未安装)

**macOS:**
```bash
brew install --cask docker
# 或下载: https://www.docker.com/products/docker-desktop
```

**Ubuntu/Debian:**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

**Windows:**
下载并安装 [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)

### 一键启动开发环境

```bash
# 1. 克隆项目
git clone https://github.com/yourusername/unifocus.git
cd unifocus

# 2. 启动所有服务
docker-compose up -d

# 3. 查看服务状态
docker-compose ps

# 4. 查看日志
docker-compose logs -f api
```

### 服务访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| API 服务 | http://localhost:8080 | Go 后端 API |
| NLP 服务 | http://localhost:8000 | Python NLP 服务 |
| PostgreSQL | localhost:5432 | 数据库 (用户: unifocus, 密码: unifocus_dev_password) |
| Redis | localhost:6379 | 缓存 |
| pgAdmin | http://localhost:5050 | 数据库管理工具 (admin@unifocus.com / admin) |

### 健康检查

```bash
# API 服务
curl http://localhost:8080/health

# NLP 服务
curl http://localhost:8000/health
```

---

## 📁 项目结构

```
unifocus/
├── backend/              # Go 后端服务
│   ├── cmd/             # 应用入口
│   ├── internal/        # 内部代码
│   │   ├── api/         # API 层
│   │   ├── domain/      # 领域模型
│   │   ├── repository/  # 数据访问层
│   │   ├── service/     # 业务逻辑层
│   │   └── crawler/     # 爬虫模块
│   ├── pkg/             # 可复用包
│   ├── migrations/      # 数据库迁移
│   └── configs/         # 配置文件
│
├── nlp-service/         # Python NLP 服务
│   ├── app/
│   │   ├── api/         # API 路由
│   │   ├── services/    # NLP 服务
│   │   └── models/      # 数据模型
│   └── requirements.txt
│
├── extension/           # 浏览器插件 (待开发)
├── web/                 # Web 管理后台 (待开发)
├── infrastructure/      # 基础设施配置
├── docs/                # 项目文档
└── docker-compose.yml   # Docker 编排配置
```

---

## 💻 本地开发

### 后端开发 (Go)

```bash
cd backend

# 安装依赖
go mod download

# 运行 API 服务
go run cmd/api/main.go

# 运行测试
go test ./...

# 构建
go build -o bin/api cmd/api/main.go
```

### NLP 服务开发 (Python)

```bash
cd nlp-service

# 创建虚拟环境
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate

# 安装依赖
pip install -r requirements.txt

# 运行服务
uvicorn app.main:app --reload --port 8000
```

### 数据库迁移

```bash
# 手动执行迁移 (连接到 PostgreSQL)
psql -h localhost -U unifocus -d unifocus_dev < backend/migrations/001_init_schema.up.sql

# 或使用 Docker
docker exec -i unifocus_postgres psql -U unifocus -d unifocus_dev < backend/migrations/001_init_schema.up.sql
```

---

## 📚 文档

- [技术蓝图](readme/technical_blueprint.md) - 完整技术实现方案
- [数据库设计](docs/database.md) - 数据库表结构说明
- [API 文档](docs/api/README.md) - RESTful API 接口文档
- [开发指南](docs/development.md) - 开发规范与最佳实践
- [部署指南](docs/deployment.md) - 生产环境部署说明

---

## 🗓️ 开发路线图

### ✅ 已完成
- [x] 项目结构搭建
- [x] 数据库 Schema 设计
- [x] Docker 开发环境配置
- [x] 后端基础框架 (Go)
- [x] NLP 服务基础框架 (Python)

### 🚧 进行中 (MVP 阶段 - 8周)
- [ ] 用户认证系统 (JWT)
- [ ] 机会数据采集 (爬虫引擎)
- [ ] NLP 文本提取与实体识别
- [ ] 基础 API 开发
- [ ] 浏览器插件开发

### 📋 计划中 (阶段二 - 12周)
- [ ] 机会可达性评分系统
- [ ] 端侧协同采集 (公众号/小红书)
- [ ] 日历视图与 DDL 管理
- [ ] 标签宇宙可视化
- [ ] 竞赛级别自动识别

---

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

## 👥 团队

- **项目负责人**: [Your Name]
- **后端开发**: [Backend Team]
- **前端开发**: [Frontend Team]
- **AI/NLP**: [AI Team]

---

## 📧 联系我们

- 项目主页: [GitHub](https://github.com/yourusername/unifocus)
- 问题反馈: [Issues](https://github.com/yourusername/unifocus/issues)
- 邮箱: unifocus@example.com

---

<div align="center">

**Built with ❤️ by UniFocus Team**

⭐ 如果这个项目对你有帮助，请给我们一个 Star！

</div>
