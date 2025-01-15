#!/bin/bash

# 服务器地址
BASE_URL="http://localhost:8080/api/v1"

# 测试设备信息
DEVICE_SN="M1A2401A0100001"
DEVICE_MODEL="MD-400D"
DEVICE_SECRET="test_device_secret"
TIMESTAMP=$(date +%s)

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 1. 注册设备
echo -e "\n${GREEN}1. 注册设备${NC}"
register_response=$(curl -s -X POST "${BASE_URL}/devices/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"sn\": \"${DEVICE_SN}\",
    \"model\": \"${DEVICE_MODEL}\"
  }")
echo $register_response

# 2. 生成签名并认证
echo -e "\n${GREEN}2. 设备认证${NC}"
# 签名格式: sha256(sn + secret + timestamp)
SIGN=$(echo -n "${DEVICE_SN}${DEVICE_SECRET}${TIMESTAMP}" | sha256sum | cut -d' ' -f1)
auth_response=$(curl -s -X POST "${BASE_URL}/devices/auth" \
  -H "Content-Type: application/json" \
  -d "{
    \"sn\": \"${DEVICE_SN}\",
    \"sign\": \"${SIGN}\",
    \"timestamp\": ${TIMESTAMP}
  }")
echo $auth_response

# 提取token
TOKEN=$(echo $auth_response | jq -r '.data.token')
if [ "$TOKEN" == "null" ]; then
    echo -e "${RED}获取token失败${NC}"
    exit 1
fi

# 3. 使用token访问健康检查接口
echo -e "\n${GREEN}3. 访问健康检查接口${NC}"
health_response=$(curl -s -X GET "${BASE_URL}/health" \
  -H "Authorization: Bearer ${TOKEN}")
echo $health_response

# 4. 刷新token
echo -e "\n${GREEN}4. 刷新token${NC}"
refresh_response=$(curl -s -X POST "${BASE_URL}/devices/refresh" \
  -H "Authorization: Bearer ${TOKEN}")
echo $refresh_response 