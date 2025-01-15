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

# 创建临时文件存储结果
RESULT_FILE=$(mktemp)
SUCCESS_FILE=$(mktemp)
FAIL_FILE=$(mktemp)

# 清理函数
cleanup() {
    rm -f "$RESULT_FILE" "$SUCCESS_FILE" "$FAIL_FILE"
}
trap cleanup EXIT

# 生成15位设备SN
generate_sn() {
    local index=$1
    # M1P + YYMM + 5位序列号 + 2位随机数
    local year_month=$(date +%y%m)  # 当前年月，如2401
    local seq=$(printf "%05d" $index)  # 5位序列号，如00001
    local random=$(printf "%04d" $(( RANDOM % 10000 )))  # 2位随机数
    # echo "M1P${year_month}${seq}${random}"  # 固定15位：M1P + 4 + 5 + 2 = 15位
    echo "M1P2204A101${random}"  # 固定15位：M1P + 4 + 5 + 2 = 15位
}

# 测试单个设备的完整流程
test_device() {
    local device_index=$1
    local start_time=$(date +%s%N)
    
    # 生成设备信息
    local DEVICE_SN=$(generate_sn $device_index)
    local DEVICE_MODEL="MD-400D"
    local TIMESTAMP=$(date +%s)

    echo -e "${YELLOW}设备 $device_index 使用 SN: $DEVICE_SN${NC}" >&2

    # 1. 注册设备
    local register_response=$(curl -s -X POST "${BASE_URL}/devices/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"sn\": \"${DEVICE_SN}\",
            \"model\": \"${DEVICE_MODEL}\"
        }")

    # 提取设备密钥和错误信息
    local DEVICE_SECRET=$(echo $register_response | jq -r '.data.secret')
    if [ "$DEVICE_SECRET" == "null" ]; then
        local error_msg=$(echo $register_response | jq -r '.message')
        echo "$device_index:注册失败: $error_msg (响应: $register_response)" >> "$RESULT_FILE"
        echo "1" >> "$FAIL_FILE"
        return
    fi

    echo -e "${YELLOW}设备 $device_index 注册成功，密钥: $DEVICE_SECRET${NC}" >&2

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
        local error_msg=$(echo $auth_response | jq -r '.message')
        echo "$device_index:认证失败: $error_msg (响应: $auth_response)" >> "$RESULT_FILE"
        echo "1" >> "$FAIL_FILE"
        return
    fi

    echo -e "${YELLOW}设备 $device_index 认证成功${NC}" >&2

    # 3. 访问健康检查接口
    local health_response=$(curl -s -X GET "${BASE_URL}/health" \
        -H "Authorization: Bearer ${TOKEN}")

    # 4. 刷新token
    local refresh_response=$(curl -s -X POST "${BASE_URL}/devices/refresh" \
        -H "Authorization: Bearer ${TOKEN}")

    local end_time=$(date +%s%N)
    local duration=$(( (end_time - start_time) / 1000000 )) # 转换为毫秒

    if [[ $refresh_response == *"code\":200"* ]]; then
        echo "$device_index:成功 (${duration}ms)" >> "$RESULT_FILE"
        echo "1" >> "$SUCCESS_FILE"
        echo -e "${YELLOW}设备 $device_index 完成所有流程${NC}" >&2
    else
        local error_msg=$(echo $refresh_response | jq -r '.message')
        echo "$device_index:失败 (${duration}ms): $error_msg (响应: $refresh_response)" >> "$RESULT_FILE"
        echo "1" >> "$FAIL_FILE"
        echo -e "${RED}设备 $device_index 刷新失败${NC}" >&2
    fi
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

# 计算结果
SUCCESS_COUNT=$(wc -l < "$SUCCESS_FILE")
FAIL_COUNT=$(wc -l < "$FAIL_FILE")

# 输出测试结果
echo -e "\n${GREEN}测试结果:${NC}"
echo "----------------------------------------"
while IFS=: read -r device_index result; do
    if [[ $result == 成功* ]]; then
        echo -e "设备 $device_index: ${GREEN}$result${NC}"
    else
        echo -e "设备 $device_index: ${RED}$result${NC}"
    fi
done < <(sort -n -t: -k1 "$RESULT_FILE")

echo "----------------------------------------"
echo -e "${GREEN}成功: $SUCCESS_COUNT${NC}"
echo -e "${RED}失败: $FAIL_COUNT${NC}"
echo -e "总计: $DEVICE_COUNT" 