# Go Style

## 原则

- 简单优先。
- 类型表达业务含义。
- 小接口。
- 不提前抽象。
- 接受 interface，返回 concrete type。
- 包名短小、清晰、无重复。

## Type-first

优先使用领域类型表达业务含义：

~~~go
type UserID string
type PlatformID int64
type SessionID string
~~~

## 让非法状态不可表示

~~~go
type HandoffReason string

const (
    HandoffReasonKeyword HandoffReason = "keyword"
    HandoffReasonImage   HandoffReason = "image"
    HandoffReasonSystem  HandoffReason = "system"
)
~~~
