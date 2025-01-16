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
TASK_ID=""

# 上报打印任务状态
report_print_status() {
    echo -e "\n${YELLOW}1. 上报打印任务状态${NC}"
    
    # 1.1 开始打印
    echo -e "\n${YELLOW}1.1 上报开始打印状态${NC}"
    echo -e "请求URL: ${BASE_URL}/device/print/status"
    echo -e "Authorization: Bearer ${TOKEN}"
    
    START_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    TASK_ID="PT$(date +%Y%m%d%H%M%S)"
    
    echo -e "\n请求数据:"
    echo "{
        \"task_id\": \"${TASK_ID}\",
        \"file_name\": \"test_model.gcode\",
        \"status\": \"printing\",
        \"start_time\": \"${START_TIME}\",
        \"progress\": 0,
        \"duration\": 0,
        \"filament_used\": 0,
        \"layers_completed\": 0
    }" | jq '.'
    
    response=$(curl -v -s -X POST "${BASE_URL}/device/print/status" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TOKEN}" \
        -d "{
            \"task_id\": \"${TASK_ID}\",
            \"file_name\": \"test_model.gcode\",
            \"status\": \"printing\",
            \"start_time\": \"${START_TIME}\",
            \"progress\": 0,
            \"duration\": 0,
            \"filament_used\": 0,
            \"layers_completed\": 0
        }")
    
    echo -e "\n响应: $response"
    
    # 1.2 更新打印进度
    echo -e "\n${YELLOW}1.2 更新打印进度${NC}"
    echo -e "请求数据:"
    echo "{
        \"task_id\": \"${TASK_ID}\",
        \"file_name\": \"test_model.gcode\",
        \"status\": \"printing\",
        \"progress\": 45.5,
        \"duration\": 1800,
        \"filament_used\": 1250.5,
        \"layers_completed\": 150
    }" | jq '.'
    
    response=$(curl -v -s -X POST "${BASE_URL}/device/print/status" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TOKEN}" \
        -d "{
            \"task_id\": \"${TASK_ID}\",
            \"file_name\": \"test_model.gcode\",
            \"status\": \"printing\",
            \"progress\": 45.5,
            \"duration\": 1800,
            \"filament_used\": 1250.5,
            \"layers_completed\": 150
        }")
    
    echo -e "\n响应: $response"
    
    # 1.3 完成打印
    echo -e "\n${YELLOW}1.3 上报打印完成状态${NC}"
    END_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    echo -e "请求数据:"
    echo "{
        \"task_id\": \"${TASK_ID}\",
        \"file_name\": \"test_model.gcode\",
        \"status\": \"completed\",
        \"end_time\": \"${END_TIME}\",
        \"progress\": 100,
        \"duration\": 3600,
        \"filament_used\": 2500.8,
        \"layers_completed\": 300
    }" | jq '.'
    
    response=$(curl -v -s -X POST "${BASE_URL}/device/print/status" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TOKEN}" \
        -d "{
            \"task_id\": \"${TASK_ID}\",
            \"file_name\": \"test_model.gcode\",
            \"status\": \"completed\",
            \"end_time\": \"${END_TIME}\",
            \"progress\": 100,
            \"duration\": 3600,
            \"filament_used\": 2500.8,
            \"layers_completed\": 300
        }")
    
    echo -e "\n响应: $response"
}

# 查询打印任务列表
get_print_tasks() {
    echo -e "\n${YELLOW}2. 查询打印任务列表${NC}"
    
    # 2.1 查询所有任务
    echo -e "\n${YELLOW}2.1 查询所有任务${NC}"
    echo -e "请求URL: ${BASE_URL}/device/print/tasks"
    
    response=$(curl -s -X GET "${BASE_URL}/device/print/tasks" \
        -H "Authorization: Bearer ${TOKEN}")
    
    echo -e "响应: $response"
    
    # 2.2 查询已完成的任务
    echo -e "\n${YELLOW}2.2 查询已完成的任务${NC}"
    echo -e "请求URL: ${BASE_URL}/device/print/tasks?status=completed"
    
    response=$(curl -s -X GET "${BASE_URL}/device/print/tasks?status=completed" \
        -H "Authorization: Bearer ${TOKEN}")
    
    echo -e "响应: $response"
}

# 查询任务历史
get_task_history() {
    if [ -z "${TASK_ID}" ]; then
        echo -e "${RED}没有找到可查询的任务${NC}"
        return
    fi
    
    echo -e "\n${YELLOW}3. 查询任务历史${NC}"
    echo -e "请求URL: ${BASE_URL}/device/print/task/${TASK_ID}/history"
    
    response=$(curl -s -X GET "${BASE_URL}/device/print/task/${TASK_ID}/history" \
        -H "Authorization: Bearer ${TOKEN}")
    
    echo -e "响应: $response"
}

# 执行测试
echo -e "${YELLOW}开始测试打印任务功能...${NC}"
echo "----------------------------------------"

# 执行测试步骤
report_print_status
get_print_tasks
get_task_history

echo -e "\n${GREEN}测试完成${NC}"
echo "----------------------------------------" 