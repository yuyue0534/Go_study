#!/bin/bash

echo "======================================"
echo "    Go电商平台 - 项目完整性验证"
echo "======================================"
echo ""

errors=0

# 检查Go文件
echo "【Go源文件检查】"
files=(
    "internal/models/models.go"
    "internal/database/database.go"
    "internal/middleware/middleware.go"
    "internal/handlers/handlers.go"
    "internal/handlers/admin.go"
    "cmd/server/main.go"
)

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        size=$(wc -l < "$file")
        echo "✓ $file ($size 行)"
    else
        echo "✗ $file - 文件缺失"
        ((errors++))
    fi
done

echo ""
echo "【HTML模板检查】"
html_files=(
    "web/templates/index.html"
    "web/templates/login.html"
    "web/templates/register.html"
    "web/templates/products.html"
    "web/templates/product_detail.html"
    "web/templates/cart.html"
    "web/templates/checkout.html"
    "web/templates/orders.html"
    "web/templates/profile.html"
    "web/templates/seller.html"
    "web/templates/admin.html"
)

for file in "${html_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✓ ${file##*/}"
    else
        echo "✗ ${file##*/} - 文件缺失"
        ((errors++))
    fi
done

echo ""
echo "【静态资源检查】"
if [ -f "web/static/css/style.css" ]; then
    echo "✓ style.css"
else
    echo "✗ style.css - 文件缺失"
    ((errors++))
fi

if [ -f "web/static/js/common.js" ]; then
    echo "✓ common.js"
else
    echo "✗ common.js - 文件缺失"
    ((errors++))
fi

echo ""
echo "【配置文件检查】"
if [ -f "go.mod" ]; then
    echo "✓ go.mod"
else
    echo "✗ go.mod - 文件缺失"
    ((errors++))
fi

if [ -f "run.sh" ]; then
    echo "✓ run.sh"
else
    echo "✗ run.sh - 文件缺失"
    ((errors++))
fi

if [ -f "run.bat" ]; then
    echo "✓ run.bat"
else
    echo "✗ run.bat - 文件缺失"
    ((errors++))
fi

echo ""
echo "【文档文件检查】"
if [ -f "README.md" ]; then
    echo "✓ README.md"
else
    echo "✗ README.md - 文件缺失"
    ((errors++))
fi

if [ -f "QUICK_START.txt" ]; then
    echo "✓ QUICK_START.txt"
else
    echo "✗ QUICK_START.txt - 文件缺失"
    ((errors++))
fi

if [ -f "STRUCTURE.txt" ]; then
    echo "✓ STRUCTURE.txt"
else
    echo "✗ STRUCTURE.txt - 文件缺失"
    ((errors++))
fi

echo ""
echo "======================================"
if [ $errors -eq 0 ]; then
    echo "✓ 验证通过！所有文件完整。"
    echo ""
    echo "项目统计："
    echo "  - Go源文件：6个"
    echo "  - HTML模板：11个"
    echo "  - 静态资源：2个"
    echo "  - 配置文件：3个"
    echo "  - 文档文件：3个"
    echo "  - 总计：25个文件"
    echo ""
    echo "下一步："
    echo "  ./run.sh        # 启动项目"
    echo "  或"
    echo "  go run cmd/server/main.go"
else
    echo "✗ 发现 $errors 个错误"
    exit 1
fi
echo "======================================"
