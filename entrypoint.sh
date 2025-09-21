#!/bin/bash

# 设置脚本在遇到错误时退出
set -e

# 获取模式参数，默认为 development
MODE=${1:-development}

echo "启动模式: $MODE"

if [ "$MODE" = "production" ]; then
    echo "生产模式: 构建并运行二进制文件"
    
    # 构建生产版本
    go build -o scheduler-go .
    
    # 运行生产版本
    ./scheduler-go
else
    echo "开发模式: 使用 go run 直接运行"
    
    # 开发模式直接运行
    go run main.go
fi
