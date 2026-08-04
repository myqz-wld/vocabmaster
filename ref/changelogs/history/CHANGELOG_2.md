# CHANGELOG_2: install path + vm alias

## 概要

修正 `make install` 的安装目录：默认安装到 `$(go env GOPATH)/bin`，不再写入 `/usr/local/bin` 或依赖 `sudo`。安装时同时创建 `vm` 短命令，用户可直接用 `vm study`。

## 变更内容

- `Makefile install` 默认解析 `$GOPATH/bin`，支持 `make install BINDIR=/path/to/bin` 覆盖安装目录。
- `install` 复制 `build/vocabmaster` 到目标目录，并创建 `vm -> vocabmaster` 符号链接；若目标目录已有非本项目的 `vm`，安装会中止并提示用户处理冲突。
- `uninstall` 同步删除 `vocabmaster` 和 `vm` 两个入口。
- `Makefile` 仅在 gvm 已安装 `go.mod` 声明版本时执行 `gvm use`，否则使用当前 `go`，避免本机只有 gvm `system` 时 `make build/install` 直接失败。
- `README.md` 安装说明补充默认路径、`BINDIR` 覆盖方式，并将使用示例切换到 `vm`。

## 验证

- `make -n build`
- `make -n install BINDIR=/tmp/vocabmaster-install-check`
- `make build`
- `make install BINDIR=<tmpdir>`
- `<tmpdir>/vocabmaster --help`
- `<tmpdir>/vm --help`
- `make uninstall BINDIR=<tmpdir>`
