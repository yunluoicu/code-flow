# Error Handling

## 规则

1. 错误向上返回，不吞错误。
2. 使用 `fmt.Errorf("...: %w", err)` 包装。
3. 对外返回前做错误分类。
4. 需要判断的错误定义 sentinel 或类型错误。
5. 日志和错误不要重复制造噪声。

示例：

~~~go
if err != nil {
    return fmt.Errorf("load session failed: %w", err)
}
~~~
