# 故障排查指南

## 常见问题

### 1. 数据库连接失败

**错误信息:**
```
Failed to initialize database: failed to ping database: dial tcp [::1]:5432: connect: connection refused
```

**解决方案:**

#### 方案A: 使用Docker Compose（推荐）

```bash
# 1. 启动数据库和Redis
cd unifocus
docker-compose up -d postgres redis

# 2. 等待数据库就绪（约5-10秒）
docker-compose ps

# 3. 初始化数据库（首次运行）
./scripts/setup-db.sh

# 4. 启动后端服务
cd backend
go run cmd/api/main.go
```

#### 方案B: 本地安装PostgreSQL

如果你本地已安装PostgreSQL：

```bash
# 1. 启动PostgreSQL服务
# macOS (使用Homebrew)
brew services start postgresql@15

# 或 Linux
sudo systemctl start postgresql

# 2. 创建数据库和用户
psql postgres
CREATE USER unifocus WITH PASSWORD 'unifocus_dev_password';
CREATE DATABASE unifocus_dev OWNER unifocus;
\q

# 3. 执行迁移
psql -U unifocus -d unifocus_dev < backend/migrations/001_init_schema.up.sql
```

### 2. Redis连接失败

**错误信息:**
```
Failed to initialize redis: dial tcp [::1]:6379: connect: connection refused
```

**解决方案:**

```bash
# 使用Docker启动Redis
docker-compose up -d redis

# 或本地安装Redis
# macOS
brew services start redis

# Linux
sudo systemctl start redis
```

### 3. Docker镜像拉取失败（网络问题）

**错误信息:**
```
dial tcp: lookup registry-1.docker.io: no such host
failed to resolve reference "docker.io/library/postgres:15-alpine"
```

**解决方案:**

#### 方案A: 配置Docker镜像加速器（推荐，适用于中国大陆）

1. **打开Docker Desktop**
2. **Settings -> Docker Engine**
3. **添加以下配置:**

```json
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.baidubce.com"
  ]
}
```

4. **点击 Apply & Restart**
5. **重新运行启动脚本**

#### 方案B: 检查网络和DNS

```bash
# 测试DNS解析
nslookup registry-1.docker.io

# 测试网络连接
ping registry-1.docker.io

# 如果无法解析，尝试更换DNS
# macOS: System Preferences -> Network -> Advanced -> DNS
# 添加: 8.8.8.8 或 114.114.114.114
```

#### 方案C: 使用代理（如果有）

在Docker Desktop中配置代理：
- Settings -> Resources -> Proxies
- 配置HTTP/HTTPS代理

#### 方案D: 手动拉取镜像

```bash
# 使用镜像加速器手动拉取
docker pull docker.mirrors.ustc.edu.cn/library/postgres:15-alpine
docker pull docker.mirrors.ustc.edu.cn/library/redis:7-alpine

# 重新标记
docker tag docker.mirrors.ustc.edu.cn/library/postgres:15-alpine postgres:15-alpine
docker tag docker.mirrors.ustc.edu.cn/library/redis:7-alpine redis:7-alpine
```

### 4. Docker未安装或未运行

**错误信息:**
```
command not found: docker
```

**解决方案:**

1. 安装Docker Desktop: https://www.docker.com/products/docker-desktop
2. 启动Docker Desktop应用
3. 验证安装: `docker --version`

### 4. 端口被占用

**错误信息:**
```
bind: address already in use
```

**解决方案:**

```bash
# 检查端口占用
# macOS/Linux
lsof -i :8080  # 检查8080端口
lsof -i :5432  # 检查5432端口
lsof -i :6379  # 检查6379端口

# 停止占用端口的进程
kill -9 <PID>

# 或修改配置文件中的端口号
```

### 5. 数据库迁移失败

**解决方案:**

```bash
# 重新初始化数据库
docker-compose down -v  # 删除数据卷（⚠️ 会删除所有数据）
docker-compose up -d postgres
./scripts/setup-db.sh
```

### 6. Go模块依赖问题

**错误信息:**
```
go: cannot find module
```

**解决方案:**

```bash
cd backend
go mod tidy
go mod download
```

### 7. Python依赖问题

**错误信息:**
```
ModuleNotFoundError: No module named 'xxx'
```

**解决方案:**

```bash
cd nlp-service
pip install -r requirements.txt

# 或使用虚拟环境
python -m venv venv
source venv/bin/activate  # macOS/Linux
pip install -r requirements.txt
```

### 8. 前端依赖问题

**错误信息:**
```
npm ERR! code ELIFECYCLE
```

**解决方案:**

```bash
cd web
rm -rf node_modules package-lock.json
npm install
```

---

## 快速诊断命令

```bash
# 运行环境检查脚本
./scripts/check-docker.sh

# 检查所有服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f postgres
docker-compose logs -f redis
docker-compose logs -f api

# 检查数据库连接
docker exec -it unifocus_postgres psql -U unifocus -d unifocus_dev -c "SELECT 1;"

# 检查Redis连接
docker exec -it unifocus_redis redis-cli ping

# 测试API健康检查
curl http://localhost:8080/health

# 测试NLP服务
curl http://localhost:8000/health

# 测试网络连接
ping registry-1.docker.io
nslookup registry-1.docker.io
```

---

## 环境要求

- **Go**: 1.21+
- **Python**: 3.11+
- **Node.js**: 18+
- **Docker**: 20.10+ (推荐)
- **PostgreSQL**: 15+ (如果本地安装)
- **Redis**: 7+ (如果本地安装)

---

## 获取帮助

如果问题仍未解决：

1. 检查日志文件: `backend/logs/app.log`
2. 查看Docker日志: `docker-compose logs`
3. 确认配置文件: `backend/configs/config.dev.yaml`
4. 检查网络连接和防火墙设置

