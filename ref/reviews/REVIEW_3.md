# REVIEW_3: install vm environment + project deep review

## Trigger

用户要求 deep review 项目，并指出 `make install` 后直接运行 `vm` 仍报 command not found；本轮先修安装入口，再用 reviewer-claude + reviewer-codex 做两轮异构评审。

## Method

- Round 1 scope: `Makefile`、`tools/install.sh`、`README.md`、CLI 根命令、学习主流程、store、library、LLM、SM-2。
- Round 1 reviewers: reviewer-codex-vocabmaster / reviewer-claude-vocabmaster。
- Rebuttal: reviewer-codex 将 marker 损坏定为 HIGH；lead 发给 reviewer-claude 反驳，最终裁为 MEDIUM。
- Round 2 scope: R1 修复文件与新增测试。

## Triage

| Finding | 裁决 | 处置 |
|---|---|---|
| `tools/install.sh` 缺 end marker 时可能吞掉 rc 后续内容 | ✅ MEDIUM | 只删除成对管理块；不配对 marker 原样保留并告警 |
| `tools/install.sh` 非原子覆盖 shell rc | ✅ MEDIUM | 同目录 temp + `mv` 原子替换；保留 rc symlink |
| `internal/library.Search` prefix 档被全量 ID sort 抹掉 | ✅ MEDIUM | exact / prefix / other 桶内排序后按优先级拼接 |
| `internal/store` 日期统计用 `Truncate(24h)` | ✅ MEDIUM | 改为本地日历日午夜窗口 |
| `Makefile` 用户路径插入单引号 `bash -lc` | ✅ LOW | 改为 export 环境变量传给安装脚本 |
| `tools/install_test.sh` 继承外部 `UPDATE_SHELL_RC` | ✅ LOW | 测试 helper 显式设置 `UPDATE_SHELL_RC=1` |
| `README.md` 通用安装步骤硬编码 `source ~/.zshrc` | ✅ LOW | 改为按安装输出 source 对应配置文件，zsh 仅作示例 |
| fish/其它 shell 自动写入缺口 | ❓ LOW | README 说明自动写入支持 zsh/bash/POSIX profile，其它 shell 手动配置 |

## Fixes

- 新增 `tools/install.sh` 和 `tools/install_test.sh`。
- `make install` 安装 `vocabmaster` / `vm` 并写 shell 配置管理块。
- `make uninstall` 只删除本项目管理的 `vm` symlink，并移除配置块。
- `make test` 现在同时执行 Go tests 和 installer shell tests。
- 新增 Search priority 回归测试与 store 本地日期窗口回归测试。

## Verification

- `bash -n tools/install.sh && bash -n tools/install_test.sh`
- `make test`
- `UPDATE_SHELL_RC=0 tools/install_test.sh`
- 临时 `BINDIR` / `SHELL_RC` 安装、重复安装、卸载。
- 带单引号的 `BINDIR` / `SHELL_RC` 走 `make install` 和 `make uninstall`。
- 真实 `make install` 写入 `~/.zshrc` 后，`source ~/.zshrc; vm --help` 通过。

## Residual Info

- `review_history.reviewed_at` 仍以 RFC3339 text 存储并做字符串区间比较；用户跨 offset 或 DST 时，日期统计仍可能有边界问题。彻底修复需要评估是否改存 UTC epoch。
- 损坏 marker 会被保留以保护用户数据；重复安装会保留孤立损坏块并追加新的成对管理块。
