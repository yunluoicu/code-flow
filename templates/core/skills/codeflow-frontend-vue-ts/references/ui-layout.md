# UI Layout

## 基本布局

复杂管理后台和 Dashboard 推荐：

- Sidebar
- Topbar
- Main Content
- Context Panel，可选

## 页面结构

页面应分为：

1. Page Header：标题、说明、核心操作。
2. Summary Cards：状态摘要。
3. Primary Content：表格、列表、详情、图表。
4. Secondary Context：风险、最近活动、筛选、说明。

## 卡片使用原则

不要默认把一切都包成 card。只有以下情况才使用 card：

- 内容是独立实体。
- 内容之间需要比较。
- 内容有明确交互边界。
- 内容是一个摘要/状态模块。

不要 card 套 card。优先用 spacing、typography、divider 建立层级。

## 必备状态

每个列表、详情、搜索、表单都需要：

- loading
- empty
- error
- success
- disabled
- hover
- active
- selected
- readonly
