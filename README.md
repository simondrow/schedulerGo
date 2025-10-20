# Go Task Scheduling API Server //Go 任务调度 API 服务器

This is a task scheduling API server based on the **Gin framework** and **MongoDB**, implemented in Go language.
这是一个基于 **Gin 框架** 和 **MongoDB** 的任务调度 API 服务器，使用 Go 语言实现。

## 技术栈

- **框架**: [Gin](https://gin-gonic.com/) - 高性能 Go Web 框架
- **数据库**: [MongoDB](https://www.mongodb.com/) - NoSQL 文档数据库
- **语言**: Go 1.21+
- **依赖管理**: Go Modules

## 功能特性

- RESTful API 设计
- MongoDB 数据库集成
- 用户任务管理
- 按日期的任务查询和更新（1-7代表周日到周六）
- 日志记录中间件
- CORS 支持
- 数据自动初始化

## 快速开始

### 环境要求

- Go 1.21 或更高版本
- MongoDB (默认连接: `mongodb://localhost:27017/scheduler`)

### 运行方法

```bash
# 克隆项目
git clone <repository-url>
cd goscheduler

# 安装依赖
go mod download

# 运行服务器
go run main.go

# 或者使用启动脚本
./entrypoint.sh          # 开发模式
./entrypoint.sh production  # 生产模式
```

服务器将在 `http://localhost:3000` 启动。

## API 接口

### 获取所有用户数据
```bash
GET /api/users
```

**响应示例**:
```json
{
  "Waner": {
    "任务": {
      "1": ["7点半滴眼药水", "7点半练反转拍", "7点半朗读英语背单词"],
      "2": ["7点滴眼药水", "7点朗读英语背单词"]
    }
  }
}
```

### 更新用户任务
```bash
POST /api/users
Content-Type: application/json

{
  "name": "用户名",
  "tasks": {
    "1": ["任务1", "任务2"],
    "2": ["任务3", "任务4"]
  }
}
```

### 获取特定用户任务
```bash
GET /api/users/{name}
```

### 获取用户特定日期任务
```bash
GET /api/users/{name}/day/{day}
```
**参数说明**:
- `day`: 1-7（1=周日，7=周六）

### 更新用户特定日期任务
```bash
POST /api/users/{name}/day/{day}
Content-Type: application/json

{
  "tasks": ["新任务1", "新任务2"]
}
```

## 快速测试

启动服务器后，使用以下命令测试 API：

```bash
# 获取所有用户数据
curl http://localhost:3000/api/users

# 获取特定用户任务
curl http://localhost:3000/api/users/Waner

# 获取用户周一的任务
curl http://localhost:3000/api/users/Waner/day/1
```

## 环境变量

- `MONGODB_URI`: MongoDB 连接字符串（默认: `mongodb://localhost:27017/scheduler`）
- `PORT`: 服务器端口（默认: `3000`）

## 项目结构
goscheduler/
├── main.go # 主程序文件
├── go.mod # Go 模块文件
├── go.sum # 依赖校验文件
├── entrypoint.sh # 启动脚本
└── README.md # 项目说明
