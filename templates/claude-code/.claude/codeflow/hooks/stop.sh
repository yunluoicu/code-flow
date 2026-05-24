#!/usr/bin/env bash
set -e
echo "## CodeFlow 完成前检查提醒"
echo "请确认：是否中文回复；是否执行 TDD 或替代验证；是否执行 /review；是否询问 verification / finishing；是否存在未经确认的 Git 写操作。"
echo "完成回复必须包含：修改内容、修改文件、测试结果、未验证内容、风险和待确认点。"
