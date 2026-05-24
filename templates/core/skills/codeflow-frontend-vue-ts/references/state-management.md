# State Management

## 规则

1. 局部状态留在组件内。
2. 跨页面共享状态使用 store。
3. API 数据、筛选条件、UI 状态要分开。
4. loading / error / empty 不要混在业务数据里。
5. store 不直接承担复杂 UI 逻辑。

## 必备状态模型

~~~ts
interface AsyncState<T> {
  loading: boolean
  error: string | null
  data: T | null
}
~~~

## 表格/列表状态

- filters
- pagination
- sorting
- selectedRowKeys
- loading
- error
- empty
