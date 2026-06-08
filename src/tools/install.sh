#!/usr/bin/env bash
set -euo pipefail

MARKER_START="# >>> vocabmaster >>>"
MARKER_END="# <<< vocabmaster <<<"
CLEANUP_FILES=()

cleanup() {
	if ((${#CLEANUP_FILES[@]})); then
		rm -f "${CLEANUP_FILES[@]}"
	fi
}
trap cleanup EXIT

usage() {
	echo "Usage: src/tools/install.sh install|uninstall" >&2
}

resolve_bindir() {
	local bindir="${BINDIR:-}"
	if [[ -z "$bindir" ]]; then
		local gopath
		gopath="$(go env GOPATH)"
		if [[ -z "$gopath" ]]; then
			echo "无法解析 GOPATH；请使用 make install BINDIR=/path/to/bin" >&2
			exit 1
		fi
		gopath="${gopath%%:*}"
		bindir="$gopath/bin"
	fi

	if [[ "$bindir" == *:* ]]; then
		echo "安装失败: BINDIR 不能包含 ':'" >&2
		exit 1
	fi

	printf '%s\n' "$bindir"
}

choose_shell_rc() {
	if [[ -n "${SHELL_RC:-}" ]]; then
		printf '%s\n' "$SHELL_RC"
		return
	fi

	case "${SHELL:-}" in
		*/zsh)
			printf '%s\n' "$HOME/.zshrc"
			;;
		*/bash)
			if [[ "$(uname)" == "Darwin" ]]; then
				printf '%s\n' "$HOME/.bash_profile"
			else
				printf '%s\n' "$HOME/.bashrc"
			fi
			;;
		*)
			printf '%s\n' "$HOME/.profile"
			;;
	esac
}

shell_quote() {
	printf '%q' "$1"
}

resolve_rc_target() {
	local rc="$1"
	if [[ ! -L "$rc" ]]; then
		printf '%s\n' "$rc"
		return
	fi

	local target
	target="$(readlink "$rc")"
	if [[ "$target" == /* ]]; then
		printf '%s\n' "$target"
	else
		printf '%s/%s\n' "$(cd "$(dirname "$rc")" && pwd)" "$target"
	fi
}

temp_in_dir() {
	local dir="$1"
	local tmp
	tmp="$(mktemp "$dir/.vocabmaster.XXXXXX")"
	CLEANUP_FILES+=("$tmp")
	printf '%s\n' "$tmp"
}

strip_managed_block() {
	local src="$1"
	local dst="$2"

	awk -v start="$MARKER_START" -v end="$MARKER_END" '
		function flush_pending(    i) {
			for (i = 1; i <= pending_count; i++) print pending[i]
			delete pending
			pending_count = 0
		}
		$0 == start {
			if (in_block) {
				print "警告: 检测到不配对的 vocabmaster marker，已保留原内容并跳过该块清理" > "/dev/stderr"
				flush_pending()
			}
			in_block = 1
			pending_count = 1
			pending[pending_count] = $0
			next
		}
		$0 == end && in_block {
			delete pending
			pending_count = 0
			in_block = 0
			next
		}
		in_block {
			pending[++pending_count] = $0
			next
		}
		{ print }
		END {
			if (in_block) {
				print "警告: 检测到不配对的 vocabmaster marker，已保留原内容并跳过该块清理" > "/dev/stderr"
				flush_pending()
			}
		}
	' "$src" > "$dst"
}

write_shell_env() {
	local bindir="$1"
	if [[ "${UPDATE_SHELL_RC:-1}" == "0" ]]; then
		if [[ ":$PATH:" != *":$bindir:"* ]]; then
			echo "提示: $bindir 不在 PATH 中；已跳过 shell 配置写入，当前 shell 需手动加入 PATH 后使用 vm"
		fi
		return
	fi

	local display_rc
	display_rc="$(choose_shell_rc)"
	mkdir -p "$(dirname "$display_rc")"
	local rc
	rc="$(resolve_rc_target "$display_rc")"
	mkdir -p "$(dirname "$rc")"
	touch "$rc"

	local rc_dir
	rc_dir="$(dirname "$rc")"
	local tmp
	tmp="$(temp_in_dir "$rc_dir")"
	strip_managed_block "$rc" "$tmp"
	local next
	next="$(temp_in_dir "$rc_dir")"
	cp -p "$rc" "$next"

	local quoted_bindir
	local quoted_binary
	quoted_bindir="$(shell_quote "$bindir")"
	quoted_binary="$(shell_quote "$bindir/vocabmaster")"

	{
		cat "$tmp"
		if [[ -s "$tmp" ]]; then
			printf '\n'
		fi
		printf '%s\n' "$MARKER_START"
		printf '__vocabmaster_bindir=%s\n' "$quoted_bindir"
		printf 'case ":$PATH:" in\n'
		printf '  *:"$__vocabmaster_bindir":*) ;;\n'
		printf '  *) export PATH="$__vocabmaster_bindir:$PATH" ;;\n'
		printf 'esac\n'
		printf 'alias vm=%s\n' "$quoted_binary"
		printf 'unset __vocabmaster_bindir\n'
		printf '%s\n' "$MARKER_END"
	} > "$next"
	mv -f "$next" "$rc"
	rm -f "$tmp"

	echo "已写入 shell 环境配置 $display_rc"
	echo "请运行 source $display_rc 或重开终端后使用 vm"
}

remove_shell_env() {
	if [[ "${UPDATE_SHELL_RC:-1}" == "0" ]]; then
		return
	fi

	local display_rc
	display_rc="$(choose_shell_rc)"
	if [[ ! -e "$display_rc" && ! -L "$display_rc" ]]; then
		return
	fi
	local rc
	rc="$(resolve_rc_target "$display_rc")"
	if [[ ! -f "$rc" ]]; then
		return
	fi

	local tmp
	tmp="$(temp_in_dir "$(dirname "$rc")")"
	strip_managed_block "$rc" "$tmp"
	local next
	next="$(temp_in_dir "$(dirname "$rc")")"
	cp -p "$rc" "$next"
	cat "$tmp" > "$next"
	mv -f "$next" "$rc"
	rm -f "$tmp"

	echo "已移除 shell 环境配置 $display_rc"
}

install_cmd() {
	local build_bin="${BUILD_BIN:-build/vocabmaster}"
	local bindir
	bindir="$(resolve_bindir)"

	if [[ ! -f "$build_bin" ]]; then
		echo "安装失败: 找不到 $build_bin；请先运行 make build" >&2
		exit 1
	fi

	mkdir -p "$bindir"
	if [[ -e "$bindir/vm" || -L "$bindir/vm" ]]; then
		local target
		target="$(readlink "$bindir/vm" 2>/dev/null || true)"
		if [[ "$target" != "vocabmaster" && "$target" != "$bindir/vocabmaster" ]]; then
			echo "安装失败: $bindir/vm 已存在，请先移除或使用 BINDIR 指定其他目录" >&2
			exit 1
		fi
	fi

	install -m 0755 "$build_bin" "$bindir/vocabmaster"
	ln -sfn "vocabmaster" "$bindir/vm"
	write_shell_env "$bindir"

	echo "已安装到 $bindir/vocabmaster"
	echo "已创建 vm 命令 $bindir/vm"
}

uninstall_cmd() {
	local bindir
	bindir="$(resolve_bindir)"

	rm -f "$bindir/vocabmaster"
	if [[ -L "$bindir/vm" ]]; then
		local target
		target="$(readlink "$bindir/vm" 2>/dev/null || true)"
		if [[ "$target" == "vocabmaster" || "$target" == "$bindir/vocabmaster" ]]; then
			rm -f "$bindir/vm"
		else
			echo "跳过 $bindir/vm：它不是 vocabmaster 创建的链接"
		fi
	elif [[ -e "$bindir/vm" ]]; then
		echo "跳过 $bindir/vm：它不是 vocabmaster 创建的链接"
	fi

	remove_shell_env
	echo "已卸载 $bindir/vocabmaster 和 vocabmaster 管理的 vm 入口"
}

main() {
	local command="${1:-}"
	case "$command" in
		install)
			install_cmd
			;;
		uninstall)
			uninstall_cmd
			;;
		*)
			usage
			exit 2
			;;
	esac
}

main "$@"
