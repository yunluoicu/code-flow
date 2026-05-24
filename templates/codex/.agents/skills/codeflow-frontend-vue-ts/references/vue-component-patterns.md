# Vue Component Patterns

## 推荐分层

- Page：页面级容器。
- Feature Section：业务区块。
- Business Component：业务组件。
- Base Component：基础组件。
- Composable：复用逻辑。
- Type / Schema：类型和接口。

## Composition API

要求：

1. 复杂逻辑放 composables。
2. API 请求封装，不直接散落在页面。
3. computed 拆成具名变量。
4. watch 必须有明确目的。
5. onMounted 不堆复杂流程。

## TypeScript

必须：

- 避免 any。
- API 返回结构显式类型化。
- 表单数据、筛选条件、列表项、详情对象必须有类型。
- 状态枚举使用 union type 或 const object。

示例：

~~~ts
export type RequirementStatus = 'draft' | 'planned' | 'executing' | 'review' | 'done' | 'archived'

export interface RequirementItem {
  id: string
  title: string
  status: RequirementStatus
  changeId?: string
  updatedAt: string
}
~~~
