#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INSTALL_SH="$ROOT_DIR/scripts/install.sh"

fail() {
	echo "install_test: $*" >&2
	exit 1
}

assert_file() {
	[[ -e "$1" ]] || fail "expected $1 to exist"
}

assert_not_exists() {
	[[ ! -e "$1" && ! -L "$1" ]] || fail "expected $1 to be absent"
}

assert_contains() {
	grep -Fq "$2" "$1" || fail "expected $1 to contain $2"
}

assert_not_contains() {
	! grep -Fq "$2" "$1" || fail "expected $1 not to contain $2"
}

make_stub_binary() {
	local path="$1"
	mkdir -p "$(dirname "$path")"
	{
		printf '#!/usr/bin/env bash\n'
		printf 'echo vocabmaster-stub\n'
	} > "$path"
	chmod +x "$path"
	cat > "$(dirname "$path")/build-info.json" <<'JSON'
{
  "name": "vocabmaster",
  "package": "github.com/vocabmaster/vocabmaster",
  "version": "test",
  "commit": "0000000000000000000000000000000000000000",
  "shortCommit": "000000000000",
  "branch": "test",
  "dirty": false,
  "builtAt": "2000-01-01T00:00:00Z"
}
JSON
}

run_install() {
	local build_bin="$1"
	local bindir="$2"
	local shell_rc="$3"
	BUILD_BIN="$build_bin" BINDIR="$bindir" SHELL_RC="$shell_rc" UPDATE_SHELL_RC=1 "$INSTALL_SH" install >/dev/null
}

run_uninstall() {
	local bindir="$1"
	local shell_rc="$2"
	BINDIR="$bindir" SHELL_RC="$shell_rc" UPDATE_SHELL_RC=1 "$INSTALL_SH" uninstall >/dev/null
}

test_basic_install_uninstall() {
	local dir build_bin bindir rc
	dir="$(mktemp -d)"
	build_bin="$dir/build/vocabmaster"
	bindir="$dir/bin"
	rc="$dir/.zshrc"
	make_stub_binary "$build_bin"

	run_install "$build_bin" "$bindir" "$rc"
	assert_file "$bindir/vocabmaster"
	assert_file "$bindir/vocabmaster.build-info.json"
	[[ -L "$bindir/vm" ]] || fail "expected vm symlink"
	assert_contains "$rc" "# >>> vocabmaster >>>"
	"$bindir/vm" | grep -Fq "vocabmaster-stub" || fail "vm symlink did not run stub"

	run_install "$build_bin" "$bindir" "$rc"
	[[ "$(grep -Fc '# >>> vocabmaster >>>' "$rc")" == "1" ]] || fail "managed block duplicated"

	run_uninstall "$bindir" "$rc"
	assert_not_exists "$bindir/vocabmaster"
	assert_not_exists "$bindir/vocabmaster.build-info.json"
	assert_not_exists "$bindir/vm"
	assert_not_contains "$rc" "# >>> vocabmaster >>>"
	rm -rf "$dir"
}

test_missing_build_metadata_rejected() {
	local dir build_bin bindir rc
	dir="$(mktemp -d)"
	build_bin="$dir/build/vocabmaster"
	bindir="$dir/bin"
	rc="$dir/.zshrc"
	make_stub_binary "$build_bin"
	rm -f "$(dirname "$build_bin")/build-info.json"

	if BUILD_BIN="$build_bin" BINDIR="$bindir" SHELL_RC="$rc" "$INSTALL_SH" install >/dev/null 2>&1; then
		fail "install should reject missing build metadata"
	fi
	assert_not_exists "$bindir/vocabmaster"
	assert_not_exists "$bindir/vocabmaster.build-info.json"
	rm -rf "$dir"
}

test_unmatched_marker_preserves_user_lines() {
	local dir build_bin bindir rc
	dir="$(mktemp -d)"
	build_bin="$dir/build/vocabmaster"
	bindir="$dir/bin"
	rc="$dir/.zshrc"
	make_stub_binary "$build_bin"
	{
		printf 'before=1\n'
		printf '# >>> vocabmaster >>>\n'
		printf 'old managed line\n'
		printf 'USER_LATER=1\n'
	} > "$rc"

	run_install "$build_bin" "$bindir" "$rc" 2>/dev/null
	assert_contains "$rc" "USER_LATER=1"
	run_install "$build_bin" "$bindir" "$rc" 2>/dev/null
	assert_contains "$rc" "USER_LATER=1"
	rm -rf "$dir"
}

test_foreign_vm_is_not_removed() {
	local dir build_bin bindir rc
	dir="$(mktemp -d)"
	build_bin="$dir/build/vocabmaster"
	bindir="$dir/bin"
	rc="$dir/.zshrc"
	make_stub_binary "$build_bin"
	mkdir -p "$bindir"
	printf 'foreign\n' > "$bindir/vm"

	if BUILD_BIN="$build_bin" BINDIR="$bindir" SHELL_RC="$rc" "$INSTALL_SH" install >/dev/null 2>&1; then
		fail "install should reject foreign vm"
	fi
	assert_contains "$bindir/vm" "foreign"

	rm -f "$bindir/vm"
	run_install "$build_bin" "$bindir" "$rc"
	ln -sfn /tmp/foreign-vm "$bindir/vm"
	run_uninstall "$bindir" "$rc"
	[[ -L "$bindir/vm" ]] || fail "foreign vm symlink should remain"
	rm -rf "$dir"
}

test_update_shell_rc_zero() {
	local dir build_bin bindir rc
	dir="$(mktemp -d)"
	build_bin="$dir/build/vocabmaster"
	bindir="$dir/bin"
	rc="$dir/.zshrc"
	make_stub_binary "$build_bin"

	BUILD_BIN="$build_bin" BINDIR="$bindir" SHELL_RC="$rc" UPDATE_SHELL_RC=0 "$INSTALL_SH" install >/dev/null
	assert_file "$bindir/vm"
	assert_not_exists "$rc"
	rm -rf "$dir"
}

test_shell_rc_symlink_is_preserved() {
	local dir build_bin bindir rc target
	dir="$(mktemp -d)"
	build_bin="$dir/build/vocabmaster"
	bindir="$dir/bin"
	rc="$dir/.zshrc"
	target="$dir/dotfiles/zshrc"
	make_stub_binary "$build_bin"
	mkdir -p "$(dirname "$target")"
	printf 'existing=1\n' > "$target"
	ln -s "$target" "$rc"

	run_install "$build_bin" "$bindir" "$rc"
	[[ -L "$rc" ]] || fail "shell rc symlink was replaced"
	assert_contains "$target" "# >>> vocabmaster >>>"

	run_uninstall "$bindir" "$rc"
	[[ -L "$rc" ]] || fail "shell rc symlink was replaced by uninstall"
	assert_not_contains "$target" "# >>> vocabmaster >>>"
	rm -rf "$dir"
}

test_basic_install_uninstall
test_missing_build_metadata_rejected
test_unmatched_marker_preserves_user_lines
test_foreign_vm_is_not_removed
test_update_shell_rc_zero
test_shell_rc_symlink_is_preserved
echo "scripts/install_test.sh ok"
