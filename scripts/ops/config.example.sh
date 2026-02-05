#!/bin/bash
# ============================================
# Edit Business 部署配置
# ============================================
#
# 此文件包含敏感信息，不应提交到 Git
# 复制 config.example.sh 并修改为实际配置值
# ============================================

# 系统配置
SYSTEM_NAME="edit-business"
BINARY_NAME="edit-business"
DOMAIN="edit.crazyaigc.com"

# SSH 配置（使用 ~/.ssh/config 中定义的别名）
SERVER="shanghai-tencent"

# 远程目录
REMOTE_DIR="/var/www/${SYSTEM_NAME}"
