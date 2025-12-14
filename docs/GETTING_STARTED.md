# 开发快速入门

本文档帮助你快速搭建 UniFocus 开发环境并开始开发。

## 📋 前置要求

### 必须安装
- **Docker Desktop** (推荐) 或 Docker + Docker Compose
- **Git**

### 可选安装（用于本地开发）
- **Go 1.21+** - 后端开发
- **Python 3.11+** - NLP服务开发
- **Node.js 18+** - 前端开发

---

## 🚀 快速启动（5分钟）

### 1. 克隆项目

```bash
git clone https://github.com/yourusername/unifocus.git
cd unifocus
```

### 2. 安装 Docker（如未安装）

**macOS:**
```bash
# 使用 Homebrew 安装
brew install --cask docker

# 或手动下载
# https://www.docker.com/products/docker-desktop
```

**Windows:**
下载并安装 [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)

**Linux (Ubuntu/Debian):**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
# 重新登录以使权限生效
```

### 3. 启动开发环境

```bash
# 方式1: 使用 Makefile（推荐）
make up

# 方式2: 使用 docker-compose
docker-compose up -d
```

等待几分钟，Docker 会自动：
- 下载镜像（PostgreSQL, Redis, Go, Python）
- 构建后端和NLP服务
- 启动所有容器
- 初始化数据库

### 4. 验证服务

```bash
# 查看服务状态
make ps
# 或
docker-compose ps

# 测试 API 服务
curl http://localhost:8080/health

# 测试 NLP 服务
curl http://localhost:8000/health

# 或使用 Makefile
make test-all
```

### 5. 访问服务

| 服务 | 地址 | 说明 |
|------|------|------|
| API 文档 | http://localhost:8080/health | Go 后端 API |
| NLP 文档 | http://localhost:8000/docs | FastAPI 交互式文档 |
| pgAdmin | http://localhost:5050 | 数据库管理（admin@unifocus.com / admin） |

---

## 💻 本地开发（无 Docker）

如果你想在本地直接运行服务而不使用 Docker：

### 后端开发

```bash
# 1. 安装 Go 依赖
cd backend
go mod download

# 2. 启动 PostgreSQL 和 Redis
# 方式1: 使用 Docker 只启动数据库
docker-compose up -d postgres redis

# 方式2: 手动安装并启动
# macOS: brew install postgresql redis
# Ubuntu: sudo apt install postgresql redis-server

# 3. 执行数据库迁移
make db-migrate

# 4. 运行 API 服务
go run cmd/api/main.go
# 或
make dev-backend
```

### NLP 服务开发

```bash
# 1. 创建虚拟环境
cd nlp-service
python -m venv venv

# 2. 激活虚拟环境
# macOS/Linux:
source venv/bin/activate
# Windows:
venv\Scripts\activate

# 3. 安装依赖
pip install -r requirements.txt

# 4. 下载 Spacy 中文模型（首次运行）
python -m spacy download zh_core_web_sm

# 5. 启动服务
uvicorn app.main:app --reload --port 8000
# 或
make dev-nlp
```

---

## 🗃️ 数据库操作

### 初始化数据库

```bash
# 执行迁移
make db-migrate
```

### 重置数据库

```bash
# 清空并重新创建所有表
make db-reset
```

### 进入数据库 Shell

```bash
# 交互式 SQL Shell
make db-shell

# 或手动连接
docker exec -it unifocus_postgres psql -U unifocus -d unifocus_dev
```

### 查看表结构

```sql
-- 列出所有表
\dt

-- 查看表结构
\d users
\d opportunities

-- 查询示例
SELECT * FROM users;
SELECT * FROM competition_level_rules;
```

---

## 🔧 常用命令

### Makefile 命令（推荐）

```bash
make help           # 显示所有可用命令
make up             # 启动所有服务
make down           # 停止所有服务
make logs           # 查看所有日志
make logs-api       # 查看 API 日志
make restart        # 重启所有服务
make clean          # 清理所有容器和数据卷
make test-all       # 测试所有服务健康状态
```

### Docker Compose 命令

```bash
docker-compose up -d              # 后台启动
docker-compose down               # 停止服务
docker-compose logs -f api        # 查看 API 日志
docker-compose restart api        # 重启 API 服务
docker-compose exec api sh        # 进入 API 容器
docker-compose ps                 # 查看服务状态
```

---

## 📝 开发工作流

### 1. 创建新功能分支

```bash
git checkout -b feature/your-feature-name
```

### 2. 修改代码

**后端 (Go):**
- 修改 `backend/` 下的代码
- Docker 开发模式会自动重启服务

**NLP 服务 (Python):**
- 修改 `nlp-service/` 下的代码
- uvicorn 的 `--reload` 会自动重载

### 3. 测试修改

```bash
# 测试 API
curl http://localhost:8080/api/v1/your-endpoint

# 查看日志
make logs-api
```

### 4. 提交代码

```bash
git add .
git commit -m "feat: add your feature description"
git push origin feature/your-feature-name
```

---

## 🐛 常见问题

### Q: Docker 启动失败，端口被占用

```bash
# 查看端口占用
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# 停止占用进程或修改 docker-compose.yml 中的端口映射
```

### Q: 数据库连接失败

```bash
# 检查 PostgreSQL 是否正常运行
docker-compose ps postgres

# 查看 PostgreSQL 日志
docker-compose logs postgres

# 重启 PostgreSQL
docker-compose restart postgres
```

### Q: Go 依赖下载缓慢

```bash
# 配置 Go 代理（中国大陆）
go env -w GOPROXY=https://goproxy.cn,direct

# 或使用阿里云镜像
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

### Q: Python 依赖安装失败

```bash
# 使用清华镜像源
pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# 或临时使用
pip install --index-url https://pypi.tuna.tsinghua.edu.cn/simple -r requirements.txt
```

### Q: Docker 镜像拉取慢

配置 Docker 镜像加速（中国大陆）：

**Docker Desktop (macOS/Windows):**
1. 打开 Docker Desktop → Settings → Docker Engine
2. 添加镜像加速器：
```json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://registry.docker-cn.com"
  ]
}
```

**Linux:**
编辑 `/etc/docker/daemon.json`：
```json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn"
  ]
}
```
然后重启 Docker：
```bash
sudo systemctl restart docker
```

---

## 📚 下一步

- 阅读 [技术蓝图](../readme/technical_blueprint.md) 了解系统架构
- 查看 [API 文档](api/README.md) 了解接口设计
- 参考 [数据库设计](database.md) 了解数据模型

---

## 🆘 需要帮助？

- 提交 [Issue](https://github.com/yourusername/unifocus/issues)
- 查看 [FAQ](FAQ.md)
- 联系团队：unifocus@example.com
