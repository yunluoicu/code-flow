# Concurrency

并发实现必须明确：

1. 谁拥有 goroutine 生命周期。
2. 如何 cancel。
3. 如何处理 panic。
4. 如何处理 error。
5. 如何退出。
6. 是否需要 worker pool。

禁止：

- 裸 `go func` 无退出机制。
- goroutine 内吞错误。
- 没有 context 的无限循环。
