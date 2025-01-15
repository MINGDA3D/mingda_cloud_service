#!/bin/bash

# 设置API基础URL
BASE_URL="http://localhost:8080/api/v1"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 测试设备信息
SN="MD3D2501130001"
MODEL="MD-1000"

echo -e "${GREEN}开始测试设备认证流程...${NC}\n"

# 1. 测试设备注册
echo "1. 测试设备注册..."
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/devices/register" \
  -H "Content-Type: application/json" \
  -d '{
    "sn": "'"${SN}"'",
    "model": "'"${MODEL}"'"
  }')

echo "注册响应: $REGISTER_RESPONSE"
echo

# 提取设备密钥
SECRET=$(echo $REGISTER_RESPONSE | grep -o '"secret":"[^"]*' | cut -d'"' -f4)

if [ -z "$SECRET" ]; then
    echo -e "${RED}获取设备密钥失败${NC}"
    exit 1
fi

echo -e "${GREEN}设备注册成功，获取到密钥${NC}"
echo

# 2. 测试设备认证
echo "2. 测试设备认证..."

# 获取当前时间戳
TIMESTAMP=$(date +%s)

# 生成签名 (sha256(sn + secret + timestamp))
SIGN=$(echo -n "${SN}${SECRET}${TIMESTAMP}" | sha256sum | cut -d' ' -f1)

echo "使用以下参数进行认证："
echo "SN: $SN"
echo "TIMESTAMP: $TIMESTAMP"
echo "SIGN: $SIGN"
echo

# 发送认证请求
AUTH_RESPONSE=$(curl -s -X POST "${BASE_URL}/devices/auth" \
  -H "Content-Type: application/json" \
  -d '{
    "sn": "'"${SN}"'",
    "timestamp": '"${TIMESTAMP}"',
    "sign": "'"${SIGN}"'"
  }')

echo "认证响应: $AUTH_RESPONSE"
echo

# 提取token
TOKEN=$(echo $AUTH_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}获取访问令牌失败${NC}"
    exit 1
fi

echo -e "${GREEN}设备认证成功，获取到令牌${NC}"
echo

# 保存认证信息供后续使用
echo "# 设备认证信息" > auth_info.txt
echo "SN=$SN" >> auth_info.txt
echo "SECRET=$SECRET" >> auth_info.txt
echo "TOKEN=$TOKEN" >> auth_info.txt

echo -e "${GREEN}认证信息已保存到 auth_info.txt${NC}"
echo
echo -e "${GREEN}测试完成${NC}" 