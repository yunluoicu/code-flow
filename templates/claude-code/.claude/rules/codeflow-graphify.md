# codeflow-graphify

Graphify 知识图谱集成约束。Graphify 为可选增强能力，不自动执行索引或更新操作。

规则：
- 无 graphify-out/graph.json 时提醒用户可执行 /graphify .
- 检测到代码文件新于图谱时提醒 /graphify . --update
- Hooks 层拦截 Graphify 写入/更新操作，需用户确认后才放行
