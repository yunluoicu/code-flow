# CodeFlow 2.0 Web Dashboard 设计规范

## 定位

CodeFlow Web Dashboard 是正式的本地开发工作台，不是临时 HTML 页面。

目标风格：

```text
Developer Workspace / AI Control Center
```

参考气质：Linear、Raycast、Atlassian、OpenDesign、本地优先工具台。

## 技术实现

```text
Go API + 嵌入式前端 Dashboard
```

运行方式保持不变：

```bash
codeflow web --port 4399
```

前端源码位于：

```text
web/
```

Go 内置发布版静态资源位于：

```text
internal/core/webui/dist/
```

## 页面布局

```text
左侧 Sidebar
顶部 Header
中间主内容区
项目卡片 / KPI / 表格 / Timeline / Search
```

## 必备页面

```text
项目总览
项目列表
项目详情
需求管理
迭代管理
OpenSpec
Superpowers
Review / Check
Graphify
全文搜索
设置
```

## 视觉要求

- 不能使用裸 HTML 表格拼页面
- 必须有统一侧边栏和顶部搜索
- 必须有 KPI 卡片
- 必须有状态色：success / warning / danger / neutral
- 必须有空状态、错误状态和建议命令
- Graphify 缺失或过期必须明显提示
- Review / Check 必须突出 Must Fix / Test Gap / Risk
