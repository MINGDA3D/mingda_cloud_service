#!/bin/bash

# 设置服务地址
BASE_URL="http://61.144.188.241:8081/api/v1"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 检查参数
if [ $# -lt 1 ]; then
    echo -e "${RED}使用方法: $0 <token> [image_path]${NC}"
    echo -e "示例: $0 \"your-token\" \"test/test_image.jpg\""
    exit 1
fi

TOKEN=$1
IMAGE_PATH=${2:-"test/test_image.jpg"}
TASK_ID="PT$(date +%Y%m%d%H%M%S)"

# 检查图片文件是否存在
if [ ! -f "$IMAGE_PATH" ]; then
    echo -e "${RED}错误: 图片文件不存在: $IMAGE_PATH${NC}"
    echo -e "请先准备测试图片，例如:"
    echo -e "mkdir -p test"
    echo -e "cp your_test_image.jpg test/test_image.jpg"
    exit 1
fi

# 上传图片
echo -e "\n${YELLOW}1. 上传图片${NC}"
echo -e "请求URL: ${BASE_URL}/device/print/image"
echo -e "Authorization: Bearer ${TOKEN}"
echo -e "TaskID: ${TASK_ID}"
echo -e "图片路径: ${IMAGE_PATH}"

response=$(curl -s -X POST "${BASE_URL}/device/print/image" \
    -H "Authorization: Bearer ${TOKEN}" \
    -F "file=@${IMAGE_PATH}" \
    -F "task_id=${TASK_ID}")

echo -e "\n响应数据:"
echo $response | jq '.'

# 检查上传是否成功
if [ "$(echo $response | jq -r '.code')" != "200" ]; then
    echo -e "${RED}图片上传失败${NC}"
    exit 1
fi

echo -e "${GREEN}图片上传成功${NC}"

# 查询图片状态
echo -e "\n${YELLOW}2. 查询图片状态${NC}"
echo -e "请求URL: ${BASE_URL}/device/print/images?task_id=${TASK_ID}"

response=$(curl -s -X GET "${BASE_URL}/device/print/images?task_id=${TASK_ID}" \
    -H "Authorization: Bearer ${TOKEN}")

echo -e "\n响应数据:"
echo $response | jq '.'

# 提取图片URL并替换域名
image_url=$(echo $response | jq -r '.data[0].ImageURL')
image_url=${image_url/localhost/61.144.188.241}
echo -e "图片URL: ${image_url}"

# 等待5秒
echo -e "\n${YELLOW}等待5秒后模拟AI回调...${NC}"
sleep 5

# 模拟AI回调
echo -e "\n${YELLOW}3. 模拟AI回调${NC}"
echo -e "请求URL: ${BASE_URL}/ai/callback"

callback_data='{
    "task_id": "'${TASK_ID}'",
    "status": "success",
    "result": {
        "has_defect": true,
        "defect_type": "stringing",
        "confidence": 0.95
    }
}'

echo -e "\n请求数据:"
echo $callback_data | jq '.'

response=$(curl -s -X POST "${BASE_URL}/ai/callback" \
    -H "Content-Type: application/json" \
    -d "$callback_data")

echo -e "\n响应数据:"
echo $response | jq '.'

# 检查回调是否成功
if [ "$(echo $response | jq -r '.code')" != "200" ]; then
    echo -e "${RED}AI回调失败${NC}"
    exit 1
fi

echo -e "${GREEN}AI回调成功${NC}"

# 最终查询结果
echo -e "\n${YELLOW}4. 查询最终结果${NC}"
echo -e "请求URL: ${BASE_URL}/device/print/images?task_id=${TASK_ID}"

response=$(curl -s -X GET "${BASE_URL}/device/print/images?task_id=${TASK_ID}" \
    -H "Authorization: Bearer ${TOKEN}")

echo -e "\n响应数据:"
echo $response | jq '.'

# 验证检测状态
status=$(echo $response | jq -r '.data[0].Status')
if [ "$status" == "2" ]; then
    echo -e "${GREEN}测试完成: 图片已完成AI检测${NC}"
else
    echo -e "${RED}测试失败: 图片检测状态异常 (status=$status)${NC}"
fi

# 显示检测结果
echo -e "\n${YELLOW}检测结果:${NC}"
echo -e "是否存在缺陷: $(echo $response | jq -r '.data[0].HasDefect')"
echo -e "缺陷类型: $(echo $response | jq -r '.data[0].DefectType')"
echo -e "置信度: $(echo $response | jq -r '.data[0].Confidence')"