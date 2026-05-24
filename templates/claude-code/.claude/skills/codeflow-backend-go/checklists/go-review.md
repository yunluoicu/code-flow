# Go Review Checklist

- [ ] 符合 gofmt / go vet。
- [ ] 没有不必要抽象。
- [ ] context 贯穿 IO 边界。
- [ ] 错误正确 wrap。
- [ ] 日志结构化。
- [ ] handler/service/repository 分层清楚。
- [ ] interface 定义在使用方。
- [ ] 有 table-driven tests。
- [ ] 有异常路径测试。
- [ ] 有并发/数据一致性风险检查。
- [ ] 无资源泄漏。
- [ ] 无全局状态污染。
- [ ] 新增配置有说明。
- [ ] 发布/迁移影响已说明。
