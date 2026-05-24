# API / Service / Repository

## Handler

- 只负责参数解析、鉴权上下文、调用 service、返回响应。
- 不写复杂业务。

## Service

- 承载业务流程。
- 不依赖 HTTP 细节。
- 处理事务边界和业务状态。

## Repository

- 封装数据库访问。
- 查询条件显式。
- update 字段显式。
- 事务边界清楚。
