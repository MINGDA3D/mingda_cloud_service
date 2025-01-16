#!/bin/bash

# 服务器地址
BASE_URL="http://localhost:8080/api/v1"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 检查是否提供了token参数
if [ -z "$1" ]; then
    echo -e "${RED}错误: 请提供token参数${NC}"
    echo -e "使用方法: $0 <token>"
    exit 1
fi

TOKEN="$1"

# 上报设备告警
report_device_alarm() {
    echo -e "\n${YELLOW}1. 上报设备告警${NC}"
    
    # 存储空间不足告警
    echo -e "\n${YELLOW}1.1 上报存储空间不足告警${NC}"
    echo -e "请求URL: ${BASE_URL}/device/alarm"
    echo -e "Authorization: Bearer ${TOKEN}"
    
    echo -e "\n请求数据:"
    echo "{
        \"alarm_type\": 1,
        \"alarm_level\": 2,
        \"alarm_value\": 95.5,
        \"alarm_desc\": \"存储空间使用率超过95%\"
    }" | jq '.'
    
    response=$(curl -v -s -X POST "${BASE_URL}/device/alarm" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TOKEN}" \
        -d "{
            \"alarm_type\": 1,
            \"alarm_level\": 2,
            \"alarm_value\": 95.5,
            \"alarm_desc\": \"存储空间使用率超过95%\"
        }")
    
    echo -e "\n响应: $response"
    
    # CPU温度过高告警
    echo -e "\n${YELLOW}1.2 上报CPU温度过高告警${NC}"
    echo -e "请求数据:"
    echo "{
        \"alarm_type\": 2,
        \"alarm_level\": 3,
        \"alarm_value\": 85.6,
        \"alarm_desc\": \"CPU温度超过85度\"
    }" | jq '.'
    
    response=$(curl -v -s -X POST "${BASE_URL}/device/alarm" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TOKEN}" \
        -d "{
            \"alarm_type\": 2,
            \"alarm_level\": 3,
            \"alarm_value\": 85.6,
            \"alarm_desc\": \"CPU温度超过85度\"
        }")
    
    echo -e "\n响应: $response"
}

# 查询设备告警
get_device_alarms() {
    echo -e "\n${YELLOW}2. 查询设备告警${NC}"
    
    # 查询所有告警
    echo -e "\n${YELLOW}2.1 查询所有告警${NC}"
    echo -e "请求URL: ${BASE_URL}/device/alarms"
    
    response=$(curl -s -X GET "${BASE_URL}/device/alarms" \
        -H "Authorization: Bearer ${TOKEN}")
    
    echo -e "响应: $response"
    
    # 查询未处理的告警
    echo -e "\n${YELLOW}2.2 查询未处理的告警${NC}"
    echo -e "请求URL: ${BASE_URL}/device/alarms?status=0"
    
    response=$(curl -s -X GET "${BASE_URL}/device/alarms?status=0" \
        -H "Authorization: Bearer ${TOKEN}")
    
    echo -e "响应: $response"
    
    # 保存第一个告警的ID用于后续测试
    ALARM_ID=$(echo $response | jq -r '.data[0].id')
}

# 处理告警
resolve_alarm() {
    if [ "$ALARM_ID" == "" ] || [ "$ALARM_ID" == "null" ]; then
        echo -e "${RED}没有找到可处理的告警${NC}"
        return
    fi
    
    echo -e "\n${YELLOW}3. 处理告警${NC}"
    echo -e "请求URL: ${BASE_URL}/device/alarm/${ALARM_ID}/resolve"
    
    echo -e "\n请求数据:"
    echo "{
        \"resolve_desc\": \"已清理临时文件，释放存储空间\"
    }" | jq '.'
    
    response=$(curl -s -X POST "${BASE_URL}/device/alarm/${ALARM_ID}/resolve" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TOKEN}" \
        -d "{
            \"resolve_desc\": \"已清理临时文件，释放存储空间\"
        }")
    
    echo -e "响应: $response"
}

# 忽略告警
ignore_alarm() {
    if [ "$ALARM_ID" == "" ] || [ "$ALARM_ID" == "null" ]; then
        echo -e "${RED}没有找到可忽略的告警${NC}"
        return
    fi
    
    echo -e "\n${YELLOW}4. 忽略告警${NC}"
    echo -e "请求URL: ${BASE_URL}/device/alarm/${ALARM_ID}/ignore"
    
    response=$(curl -s -X POST "${BASE_URL}/device/alarm/${ALARM_ID}/ignore" \
        -H "Authorization: Bearer ${TOKEN}")
    
    echo -e "响应: $response"
}

# 执行测试
echo -e "${YELLOW}开始测试设备告警功能...${NC}"
echo "----------------------------------------"

# 执行测试步骤
report_device_alarm
get_device_alarms
resolve_alarm
ignore_alarm

echo -e "\n${GREEN}测试完成${NC}"
echo "----------------------------------------" 