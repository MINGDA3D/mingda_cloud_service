#!/bin/bash

# 服务器地址
BASE_URL="http://localhost:8080/api/v1"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 设备信息
DEVICE_SN="M1P2401A0100005"
DEVICE_MODEL="MD-400D"

# 检查是否提供了token参数
if [ -z "$1" ]; then
    echo -e "${RED}错误: 请提供token参数${NC}"
    echo -e "使用方法: $0 <token>"
    exit 1
fi

TOKEN="$1"

echo -e "${GREEN}开始上报设备信息${NC}"
echo "----------------------------------------"

echo -e "请求URL: ${BASE_URL}/device/info"
echo -e "Authorization: Bearer ${TOKEN}"

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

response=$(curl -v -s -X POST "${BASE_URL}/device/info" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${TOKEN}" \
    -d @- << EOF
{
    "device_info": {
        "device_sn": "${DEVICE_SN}",
        "device_model": "${DEVICE_MODEL}",
        "hardware_version": "V1.0"
    },
    "software_versions": {
        "klipper": "v0.11.0",
        "klipper_screen": "v1.0.0",
        "moonraker": "v0.8.0",
        "mainsail": "v2.5.0",
        "crowsnest": "v1.0.0",
        "firmware": {
            "mainboard": "v1.2.3",
            "printhead": "v1.0.1",
            "leveling": "v1.1.0"
        }
    }
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
    echo -e "\n${GREEN}设备信息上报成功${NC}"
else
    echo -e "\n${RED}设备信息上报失败${NC}"
    echo -e "${RED}错误代码: $code${NC}"
    echo -e "${RED}错误信息: $(echo "$json_response" | jq -r '.message')${NC}"
    echo -e "${RED}错误详情: $(echo "$json_response" | jq -r '.data // empty')${NC}"
    exit 1
fi 