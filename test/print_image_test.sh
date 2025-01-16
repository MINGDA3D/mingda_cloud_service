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
    echo -e "使用方法: $0 <token> [image_path]"
    exit 1
fi

TOKEN="$1"

# 检查是否提供了图片路径，如果没有则使用默认测试图片
IMAGE_PATH=${2:-"test/test_image.jpg"}
if [ ! -f "$IMAGE_PATH" ]; then
    echo -e "${RED}错误: 图片文件不存在: $IMAGE_PATH${NC}"
    exit 1
fi

# 生成测试用的任务ID
TASK_ID="PT$(date +%Y%m%d%H%M%S)"

# 上传打印图片
upload_print_image() {
    echo -e "\n${YELLOW}1. 上传打印图片${NC}"
    echo -e "请求URL: ${BASE_URL}/device/print/image"
    echo -e "Authorization: Bearer ${TOKEN}"
    echo -e "任务ID: ${TASK_ID}"
    echo -e "图片路径: ${IMAGE_PATH}"
    
    response=$(curl -v -s -X POST "${BASE_URL}/device/print/image" \
        -H "Authorization: Bearer ${TOKEN}" \
        -F "task_id=${TASK_ID}" \
        -F "image=@${IMAGE_PATH}")
    
    echo -e "\n响应内容:"
    echo "$response" | jq '.'
    
    # 检查响应状态
    code=$(echo "$response" | jq -r '.code')
    if [ "$code" = "200" ]; then
        echo -e "\n${GREEN}图片上传成功${NC}"
        # 保存图片URL用于后续验证
        IMAGE_URL=$(echo "$response" | jq -r '.data.image_url')
        echo -e "图片访问地址: ${IMAGE_URL}"
    else
        echo -e "\n${RED}图片上传失败${NC}"
        echo -e "${RED}错误代码: $code${NC}"
        echo -e "${RED}错误信息: $(echo "$response" | jq -r '.message')${NC}"
        echo -e "${RED}错误详情: $(echo "$response" | jq -r '.data // empty')${NC}"
        exit 1
    fi
}

# 获取打印图片列表
get_print_images() {
    echo -e "\n${YELLOW}2. 获取打印图片列表${NC}"
    
    # 2.1 获取所有图片
    echo -e "\n${YELLOW}2.1 获取所有图片${NC}"
    echo -e "请求URL: ${BASE_URL}/device/print/images"
    
    response=$(curl -s -X GET "${BASE_URL}/device/print/images" \
        -H "Authorization: Bearer ${TOKEN}")
    
    echo -e "\n响应内容:"
    echo "$response" | jq '.'
    
    # 2.2 获取指定任务的图片
    echo -e "\n${YELLOW}2.2 获取指定任务的图片${NC}"
    echo -e "请求URL: ${BASE_URL}/device/print/images?task_id=${TASK_ID}"
    
    response=$(curl -s -X GET "${BASE_URL}/device/print/images?task_id=${TASK_ID}" \
        -H "Authorization: Bearer ${TOKEN}")
    
    echo -e "\n响应内容:"
    echo "$response" | jq '.'
}

# 验证图片是否可访问
verify_image_access() {
    if [ -n "$IMAGE_URL" ]; then
        echo -e "\n${YELLOW}3. 验证图片是否可访问${NC}"
        echo -e "图片URL: ${IMAGE_URL}"
        
        http_code=$(curl -s -o /dev/null -w "%{http_code}" "$IMAGE_URL")
        
        if [ "$http_code" = "200" ]; then
            echo -e "${GREEN}图片可以正常访问${NC}"
        else
            echo -e "${RED}图片访问失败，HTTP状态码: $http_code${NC}"
        fi
    fi
}

# 执行测试
echo -e "${YELLOW}开始测试打印图片功能...${NC}"
echo "----------------------------------------"

# 执行测试步骤
upload_print_image
get_print_images
verify_image_access

echo -e "\n${GREEN}测试完成${NC}"
echo "----------------------------------------" 