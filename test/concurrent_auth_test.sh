#!/bin/bash

# 服务器地址
BASE_URL="http://localhost:8080/api/v1"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 设备数量
DEVICE_COUNT=10

# 存储测试结果
SUCCESS_COUNT=0
FAIL_COUNT=0
declare -A RESULTS

# 测试单个设备的完整流程
test_device() {
    local device_index=$1
    local start_time=$(date +%s%N)
    
    # 生成设备信息
    local DEVICE_SN="M1P2004A1${device_index}"
    local DEVICE_MODEL="MD-400D"
    local TIMESTAMP=$(date +%s)
    local result=""

    # 1. 注册设备
    local register_response=$(curl -s -X POST "${BASE_URL}/devices/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"sn\": \"${DEVICE_SN}\",
            \"model\": \"${DEVICE_MODEL}\"
        }")

    # 提取设备密钥
    local DEVICE_SECRET=$(echo $register_response | jq -r '.data.secret')
    if [ "$DEVICE_SECRET" == "null" ]; then
        result="注册失败"
        RESULTS[$device_index]="$result"
        return
    fi

    # 2. 设备认证
    local SIGN=$(echo -n "${DEVICE_SN}${DEVICE_SECRET}${TIMESTAMP}" | sha256sum | cut -d' ' -f1)
    local auth_response=$(curl -s -X POST "${BASE_URL}/devices/auth" \
        -H "Content-Type: application/json" \
        -d "{
            \"sn\": \"${DEVICE_SN}\",
            \"sign\": \"${SIGN}\",
            \"timestamp\": ${TIMESTAMP}
        }")

    # 提取token
    local TOKEN=$(echo $auth_response | jq -r '.data.token')
    if [ "$TOKEN" == "null" ]; then
        result="认证失败"
        RESULTS[$device_index]="$result"
        return
    fi

    # 3. 访问健康检查接口
    local health_response=$(curl -s -X GET "${BASE_URL}/health" \
        -H "Authorization: Bearer ${TOKEN}")

    # 4. 刷新token
    local refresh_response=$(curl -s -X POST "${BASE_URL}/devices/refresh" \
        -H "Authorization: Bearer ${TOKEN}")

    local end_time=$(date +%s%N)
    local duration=$(( (end_time - start_time) / 1000000 )) # 转换为毫秒

    if [[ $refresh_response == *"code\":200"* ]]; then
        result="成功 (${duration}ms)"
        ((SUCCESS_COUNT++))
    else
        result="失败 (${duration}ms): $(echo $refresh_response | jq -r '.message')"
        ((FAIL_COUNT++))
    fi

    RESULTS[$device_index]="$result"
}

# 清理之前的测试结果
echo -e "\n${GREEN}开始并发测试 - ${DEVICE_COUNT}个设备${NC}"
echo "----------------------------------------"

# 并发执行测试
for i in $(seq 1 $DEVICE_COUNT); do
    test_device $i &
done

# 等待所有测试完成
wait

# 输出测试结果
echo -e "\n${GREEN}测试结果:${NC}"
echo "----------------------------------------"
for i in $(seq 1 $DEVICE_COUNT); do
    if [[ ${RESULTS[$i]} == 成功* ]]; then
        echo -e "设备 $i: ${GREEN}${RESULTS[$i]}${NC}"
    else
        echo -e "设备 $i: ${RED}${RESULTS[$i]}${NC}"
    fi
done
echo "----------------------------------------"
echo -e "${GREEN}成功: $SUCCESS_COUNT${NC}"
echo -e "${RED}失败: $FAIL_COUNT${NC}"
echo -e "总计: $DEVICE_COUNT" 