# Go Security

后端安全检查：

1. 不打印 secret / token。
2. SQL / ES query / Mongo filter 防注入。
3. HTTP handler 参数校验。
4. 权限校验不放在外层临时判断。
5. goroutine / channel / context 避免泄漏。
6. 日志脱敏。
7. 错误返回不暴露内部实现。
