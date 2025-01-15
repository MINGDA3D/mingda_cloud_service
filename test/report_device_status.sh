#!/bin/bash

# 服务器地址
BASE_URL="http://localhost:8080/api/v1"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 检查是否提供了token参数
if [ -z "$1" ]; then
    echo -e "${RED}错误: 请提供token参数${NC}"
    echo -e "使用方法: $0 <token>"
    exit 1
fi

TOKEN="$1"

echo -e "${GREEN}开始上报设备状态${NC}"
echo "----------------------------------------"

echo -e "请求URL: ${BASE_URL}/device/status"
echo -e "Authorization: Bearer ${TOKEN}"

# 模拟设备状态数据
echo -e "\n请求数据:"
echo "{
    \"storage_total\": 32768,
    \"storage_used\": 15360,
    \"storage_free\": 17408,
    \"cpu_usage\": 35.6,
    \"cpu_temperature\": 45.8,
    \"memory_total\": 4096,
    \"memory_used\": 2048,
    \"memory_free\": 2048
}" | jq '.'

# 发送请求
response=$(curl -v -s -X POST "${BASE_URL}/device/status" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${TOKEN}" \
    -d @- << EOF
{
    "storage_total": 32768,
    "storage_used": 15360,
    "storage_free": 17408,
    "cpu_usage": 35.6,
    "cpu_temperature": 45.8,
    "memory_total": 4096,
    "memory_used": 2048,
    "memory_free": 2048
}
EOF
)

echo -e "\n完整响应内容:"
echo "$response"

# 提取 JSON 响应部分
json_response=$(echo "$response" | grep -A 1000 "{" | grep -B 1000 "}")

# 检查响应状态
code=$(echo "$json_response" | jq -r '.code')
if [ "$code" = "200" ]; then
    echo -e "\n${GREEN}设备状态上报成功${NC}"
    
    # 查询数据库记录
    echo -e "\n数据库记录:"
    echo -e "${GREEN}1. 设备状态记录:${NC}"
    mysql -u mingda -pmingda3D250113 md_device_db -e "SELECT * FROM md_device_status ORDER BY id DESC LIMIT 1\G"
    
    echo -e "\n${GREEN}2. 设备在线状态:${NC}"
    mysql -u mingda -pmingda3D250113 md_device_db -e "SELECT * FROM md_device_online ORDER BY id DESC LIMIT 1\G"
else
    echo -e "\n${RED}设备状态上报失败${NC}"
    echo -e "${RED}错误代码: $code${NC}"
    echo -e "${RED}错误信息: $(echo "$json_response" | jq -r '.message')${NC}"
    echo -e "${RED}错误详情: $(echo "$json_response" | jq -r '.data // empty')${NC}"
    exit 1
fi 