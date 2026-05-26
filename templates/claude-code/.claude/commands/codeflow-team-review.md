# /codeflow-team-review

创建 Claude Code Agent Team 进行并行 Review。Spawn security / performance / test / maintainability 四个 reviewer，各自审查后由 Team lead 综合输出最终 Review Result。

使用方式：

```text
/codeflow-team-review <审查范围描述>
```

适用场景：重要 PR 需要多维度审查、核心链路变更。前置条件：仅 Claude Code 支持；需设置 CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1。token 成本更高，创建前会提醒用户确认。

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
