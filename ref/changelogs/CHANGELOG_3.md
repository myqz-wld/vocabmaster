# CHANGELOG_3: install writes vm shell environment

## Summary

`make install` 不再只创建 `vm` symlink 和打印 PATH 提示；安装会写入受控 shell 配置块，把安装目录加入 PATH，并添加 `vm` alias。首次安装后用户 source shell 配置或重开终端即可直接运行 `vm`。本轮 deep-review 同步修复了搜索排序与日期统计两个既有问题。

## Changes

- `Makefile install/uninstall` 改为调用 `tools/install.sh`，保留默认 `$(go env GOPATH)/bin` 和 `BINDIR=/path/to/bin` 覆盖。
- `tools/install.sh install` 安装 `vocabmaster`、创建 `vm` symlink，并向 shell 配置写入 `# >>> vocabmaster >>>` 管理块。
- shell 配置管理块只在 start/end marker 成对时删除；发现不配对 marker 时保留原内容并告警，避免吞掉用户后续 rc 内容。
- shell 配置落盘改为同目录临时文件 + `mv` 原子替换，避免中断时截断用户 rc；若 rc 是 symlink，则写入 symlink 目标并保留链接本身。
- `tools/install.sh uninstall` 删除 `vocabmaster`、仅删除本项目管理的 `vm` symlink，并移除管理块。
- `Makefile` 通过导出环境变量向安装脚本传参，避免含单引号的 `BINDIR` / `SHELL_RC` 破坏 `bash -lc`。
- `library.Search` 改为 exact / prefix / other 各桶内按 ID 排序后拼接，恢复 prefix 优先级。
- `store.GetReviewCountOnDate` / streak 日期窗口改用本地日历日午夜，避免 `Truncate(24h)` 在 Asia/Shanghai 等时区偏到 08:00。
- `README.md` 补充首次安装后 `source ~/.zshrc` 或重开终端，以及 `UPDATE_SHELL_RC=0` 跳过 shell 配置写入。
- `tools/install_test.sh` 覆盖安装/卸载、重复安装、marker 损坏、foreign `vm`、跳过 shell rc 写入。
- `tools/install_test.sh` 显式隔离 `UPDATE_SHELL_RC`，避免调用者环境变量影响测试结果。

## Verification

- `bash -n tools/install.sh`
- `bash -n tools/install_test.sh`
- `make test`
- Isolated `make install BINDIR=<tmp>/bin SHELL_RC=<tmp>/.zshrc`
- `<tmp>/bin/vm --help`
- `zsh -lc "source <tmp>/.zshrc; command -v vm; vm --help"`
- Repeated isolated install keeps one managed shell block
- Isolated `make uninstall BINDIR=<tmp>/bin SHELL_RC=<tmp>/.zshrc`
- Isolated `make install` / `make uninstall` with single quote in `BINDIR` and `SHELL_RC`
- `tools/install_test.sh` symlink rc fixture preserves the symlink across install/uninstall
- `UPDATE_SHELL_RC=0 tools/install_test.sh`
