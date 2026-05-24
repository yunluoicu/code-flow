# Context

## 规则

1. IO / RPC / DB / Cache / MQ / HTTP 必须传 `context.Context`。
2. context 不存业务可选参数。
3. 不把 context 存进 struct。
4. 后台任务要考虑 timeout / cancel。
5. goroutine 必须可退出。
