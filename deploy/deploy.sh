#!/bin/bash

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 设置Go环境变量
export PATH=$PATH:/usr/local/go/bin
export GOROOT=/usr/local/go
export GOPATH=/home/yxs/go

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: Go 未安装或未在PATH中找到${NC}"
    echo -e "请先安装Go或确保Go的安装路径在PATH中"
    exit 1
fi

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}请使用root权限运行此脚本${NC}"
    echo -e "使用方法: sudo ./deploy.sh"
    exit 1
fi

# 设置工作目录
WORK_DIR="/home/yxs/code/mingda_cloud_service"
SERVICE_NAME="mingda-cloud"

echo -e "${YELLOW}开始部署 Mingda Cloud Service...${NC}"

# 1. 编译程序
echo -e "\n${YELLOW}1. 编译程序${NC}"
cd $WORK_DIR
go build -o mingda_cloud_service cmd/server/main.go
if [ $? -ne 0 ]; then
    echo -e "${RED}编译失败${NC}"
    exit 1
fi
echo -e "${GREEN}编译成功${NC}"

# 2. 创建必要的目录
echo -e "\n${YELLOW}2. 创建必要的目录${NC}"
mkdir -p $WORK_DIR/logs
mkdir -p $WORK_DIR/uploads/images
chmod 755 $WORK_DIR/uploads/images
echo -e "${GREEN}目录创建完成${NC}"

# 3. 复制systemd服务文件
echo -e "\n${YELLOW}3. 安装systemd服务${NC}"
cp $WORK_DIR/deploy/mingda-cloud.service /etc/systemd/system/
chmod 644 /etc/systemd/system/mingda-cloud.service

# 4. 重新加载systemd配置
echo -e "\n${YELLOW}4. 重新加载systemd配置${NC}"
systemctl daemon-reload

# 5. 启动服务
echo -e "\n${YELLOW}5. 启动服务${NC}"
systemctl stop $SERVICE_NAME 2>/dev/null
systemctl start $SERVICE_NAME
if [ $? -ne 0 ]; then
    echo -e "${RED}服务启动失败，请检查日志${NC}"
    exit 1
fi

# 6. 设置开机自启
echo -e "\n${YELLOW}6. 设置开机自启${NC}"
systemctl enable $SERVICE_NAME

# 7. 检查服务状态
echo -e "\n${YELLOW}7. 检查服务状态${NC}"
systemctl status $SERVICE_NAME

echo -e "\n${GREEN}部署完成！${NC}"
echo -e "可以使用以下命令管理服务："
echo -e "  启动服务: ${YELLOW}sudo systemctl start $SERVICE_NAME${NC}"
echo -e "  停止服务: ${YELLOW}sudo systemctl stop $SERVICE_NAME${NC}"
echo -e "  重启服务: ${YELLOW}sudo systemctl restart $SERVICE_NAME${NC}"
echo -e "  查看状态: ${YELLOW}sudo systemctl status $SERVICE_NAME${NC}"
echo -e "  查看日志: ${YELLOW}sudo journalctl -u $SERVICE_NAME -f${NC}" 