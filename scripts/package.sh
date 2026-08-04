#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="$ROOT_DIR/build"
DIST_DIR="$ROOT_DIR/dist"
BINARY="$BUILD_DIR/vocabmaster"
BUILD_INFO="$BUILD_DIR/build-info.json"

if [[ ! -x "$BINARY" || ! -f "$BUILD_INFO" ]]; then
	echo "打包失败: 缺少 build/vocabmaster 或 build/build-info.json；请先运行 make build" >&2
	exit 1
fi

json_string() {
	local key="$1"
	awk -F '"' -v key="$key" '$2 == key { print $4; exit }' "$BUILD_INFO"
}

version="$(json_string version)"
short_commit="$(json_string shortCommit)"
release_id="${version:-$short_commit}"
if [[ -z "$release_id" ]]; then
	release_id="unknown"
fi

platform="$(go env GOOS)-$(go env GOARCH)"
package_name="vocabmaster-${release_id}-${platform}"
archive="$DIST_DIR/${package_name}.tar.gz"
checksum="${archive}.sha256"
stage_root="$(mktemp -d "${TMPDIR:-/tmp}/vocabmaster-pack.XXXXXX")"
trap 'rm -rf "$stage_root"' EXIT

package_dir="$stage_root/$package_name"
mkdir -p "$package_dir" "$DIST_DIR"
install -m 0755 "$BINARY" "$package_dir/vocabmaster"
install -m 0644 "$BUILD_INFO" "$package_dir/build-info.json"
install -m 0644 "$ROOT_DIR/README.md" "$package_dir/README.md"
ln -s vocabmaster "$package_dir/vm"

tar -C "$stage_root" -czf "$archive" "$package_name"
if command -v shasum >/dev/null 2>&1; then
	(
		cd "$DIST_DIR"
		shasum -a 256 "$(basename "$archive")" > "$(basename "$checksum")"
	)
else
	(
		cd "$DIST_DIR"
		sha256sum "$(basename "$archive")" > "$(basename "$checksum")"
	)
fi

echo "已生成 $archive"
echo "已生成 $checksum"
