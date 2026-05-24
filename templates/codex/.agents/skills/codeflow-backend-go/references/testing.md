# Go Testing

必须优先考虑：

- table-driven tests。
- 边界场景。
- 错误场景。
- 并发场景。
- 数据为空。
- 外部依赖失败。

示例：

~~~go
func TestPolicyRouter_NoHandoffWhenKnowledgeMiss(t *testing.T) {}
~~~

推荐执行：

~~~bash
go test ./...
go test -race ./...
go vet ./...
~~~
