#!/bin/bash

# 服务器地址
BASE_URL="http://localhost:8080/api/v1"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 测试设备信息
DEVICE_SN="M1P2401A0100001"
DEVICE_MODEL="MD-400D"

# 临时文件
TOKEN_FILE=$(mktemp)

# 清理函数
cleanup() {
    rm -f "$TOKEN_FILE"
}
trap cleanup EXIT

# 注册设备并获取密钥
register_device() {
    echo -e "\n${GREEN}1. 注册设备${NC}"
    response=$(curl -s -X POST "${BASE_URL}/devices/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"sn\": \"${DEVICE_SN}\",
            \"model\": \"${DEVICE_MODEL}\"
        }")
    
    # 提取设备密钥
    DEVICE_SECRET=$(echo $response | jq -r '.data.secret')
    if [ "$DEVICE_SECRET" == "null" ]; then
        echo -e "${RED}设备注册失败: $(echo $response | jq -r '.message')${NC}"
        exit 1
    fi
    echo -e "${GREEN}设备注册成功，密钥: $DEVICE_SECRET${NC}"
}

# 设备认证并获取token
authenticate_device() {
    echo -e "\n${GREEN}2. 设备认证${NC}"
    TIMESTAMP=$(date +%s)
    SIGN=$(echo -n "${DEVICE_SN}${DEVICE_SECRET}${TIMESTAMP}" | sha256sum | cut -d' ' -f1)
    
    response=$(curl -s -X POST "${BASE_URL}/devices/auth" \
        -H "Content-Type: application/json" \
        -d "{
            \"sn\": \"${DEVICE_SN}\",
            \"sign\": \"${SIGN}\",
            \"timestamp\": ${TIMESTAMP}
        }")

    # 提取token
    TOKEN=$(echo $response | jq -r '.data.token')
    if [ "$TOKEN" == "null" ]; then
        echo -e "${RED}设备认证失败: $(echo $response | jq -r '.message')${NC}"
        exit 1
    fi
    echo "$TOKEN" > "$TOKEN_FILE"
    echo -e "${GREEN}设备认证成功，获取到token${NC}"
}

# 上报设备信息
report_device_info() {
    echo -e "\n${GREEN}3. 上报设备信息${NC}"
    TOKEN=$(cat "$TOKEN_FILE")
    
    echo -e "请求URL: ${BASE_URL}/device/info"
    echo -e "Authorization: Bearer ${TOKEN}"
    
    response=$(curl -v -s -X POST "${BASE_URL}/device/info" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TOKEN}" \
        -d "{
            \"device_info\": {
                \"device_sn\": \"${DEVICE_SN}\",
                \"device_model\": \"${DEVICE_MODEL}\",
                \"hardware_version\": \"V1.0\"
            },
            \"software_versions\": {
                \"klipper\": \"v0.11.0\",
                \"klipper_screen\": \"v1.0.0\",
                \"moonraker\": \"v0.8.0\",
                \"mainsail\": \"v2.5.0\",
                \"crowsnest\": \"v1.0.0\",
                \"firmware\": {
                    \"mainboard\": \"v1.2.3\",
                    \"printhead\": \"v1.0.1\",
                    \"leveling\": \"v1.1.0\"
                }
            }
        }" 2>&1)

    echo -e "\n请求数据:"
    echo "{
        \"device_info\": {
            \"device_sn\": \"${DEVICE_SN}\",
            \"device_model\": \"${DEVICE_MODEL}\",
            \"hardware_version\": \"V1.0\"
        },
        \"software_versions\": {
            \"klipper\": \"v0.11.0\",
            \"klipper_screen\": \"v1.0.0\",
            \"moonraker\": \"v0.8.0\",
            \"mainsail\": \"v2.5.0\",
            \"crowsnest\": \"v1.0.0\",
            \"firmware\": {
                \"mainboard\": \"v1.2.3\",
                \"printhead\": \"v1.0.1\",
                \"leveling\": \"v1.1.0\"
            }
        }
    }" | jq '.'

    echo -e "\n完整响应内容:"
    echo "$response"
    
    # 提取JSON响应部分
    json_response=$(echo "$response" | grep -A 1000 "{" | grep -B 1000 "}")
    
    if [ -n "$json_response" ]; then
        echo -e "\nJSON响应解析:"
        echo "$json_response" | jq '.'
        
        if [ "$(echo $json_response | jq -r '.code')" == "200" ]; then
            echo -e "${GREEN}设备信息上报成功${NC}"
        else
            echo -e "${RED}设备信息上报失败${NC}"
            echo -e "${RED}错误代码: $(echo $json_response | jq -r '.code')${NC}"
            echo -e "${RED}错误信息: $(echo $json_response | jq -r '.message')${NC}"
            echo -e "${RED}错误详情: $(echo $json_response | jq -r '.data // empty')${NC}"
            exit 1
        fi
    else
        echo -e "${RED}响应不是有效的JSON格式${NC}"
        exit 1
    fi
}

# 测试不同的错误情况
test_error_cases() {
    echo -e "\n${GREEN}4. 测试错误情况${NC}"
    TOKEN=$(cat "$TOKEN_FILE")

    echo -e "\n${YELLOW}4.1 测试缺少必填字段${NC}"
    response=$(curl -s -X POST "${BASE_URL}/device/info" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TOKEN}" \
        -d "{
            \"device_info\": {
                \"device_sn\": \"${DEVICE_SN}\"
            }
        }")
    echo -e "响应: $(echo $response | jq '.')"

    echo -e "\n${YELLOW}4.2 测试SN不匹配${NC}"
    response=$(curl -s -X POST "${BASE_URL}/device/info" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${TOKEN}" \
        -d "{
            \"device_info\": {
                \"device_sn\": \"M1P2401A0100002\",
                \"device_model\": \"${DEVICE_MODEL}\",
                \"hardware_version\": \"V1.0\"
            },
            \"software_versions\": {
                \"klipper\": \"v0.11.0\",
                \"klipper_screen\": \"v1.0.0\",
                \"firmware\": {
                    \"mainboard\": \"v1.2.3\",
                    \"printhead\": \"v1.0.1\"
                }
            }
        }")
    echo -e "响应: $(echo $response | jq '.')"
}

# 主测试流程
main() {
    echo -e "${GREEN}开始测试设备信息采集功能${NC}"
    echo "----------------------------------------"
    
    # 1. 注册设备
    register_device
    
    # 2. 设备认证
    authenticate_device
    
    # 3. 上报设备信息
    report_device_info
    
    # 4. 测试错误情况
    test_error_cases
    
    echo -e "\n${GREEN}测试完成${NC}"
    echo "----------------------------------------"
}

# 执行测试
main 