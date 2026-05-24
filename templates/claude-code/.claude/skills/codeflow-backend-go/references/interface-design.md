# Interface Design

## 原则

- interface 定义在使用方。
- 小接口优先。
- 不要为了 mock 提前抽大接口。
- 接受 interface，返回 struct。

示例：

~~~go
type SessionStore interface {
    GetByID(ctx context.Context, id SessionID) (*Session, error)
    Save(ctx context.Context, s *Session) error
}
~~~
