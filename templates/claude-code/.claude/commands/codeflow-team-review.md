# /codeflow-team-review

创建 Claude Code Agent Team 进行并行 Review。

要求：
1. 先提醒用户 Agent Teams 是实验性能力且 token 成本更高。
2. 创建前必须用户确认。
3. Spawn teammates:
   - security-reviewer
   - performance-reviewer
   - test-reviewer
   - maintainability-reviewer
4. 不自动修改代码。
5. 输出 Must Fix / Should Fix / Test Gap / Risk / Evidence。
6. Team lead 综合最终 Review Result。
