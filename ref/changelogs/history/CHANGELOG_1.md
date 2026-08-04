# CHANGELOG_1: build-dir-migration-20260526 plan 收口

## 概要

把 vocabmaster (Go CLI) 编译产物从根 `./vocabmaster` 迁到 `./build/vocabmaster`,对齐 agent-deck 项目 build-dir-migration-20260526(commit `6a6903e`)+ cross-language `build/` canonical(应用打包 CLAUDE.md §src/build)。impl 硬切 `build/`(不留任何根 binary fallback / migration helper),同步改 1 Makefile + 1 .gitignore + 1 README;Go 源码 0 hardcode binary path 引用(`Phase B empty`)。

deep-review × 1 轮 mixed kind(reviewer-claude Opus 4.7 + reviewer-codex gpt-5.5 xhigh 异构对抗)共 14 finding fix loop 收敛:R1 reviewer-claude 8 finding(1 HIGH/3 MED/2 LOW/2 INFO)+ reviewer-codex 6 finding(0 HIGH/2 MED/3 LOW/1 INFO),共 10 ✅ 必修(1 双方独立 + lead 全自检 + 1 HIGH lead 现场 reproduce ff-merge divergence)+ 2 INFO verified。**0 反驳 + 0 新 HIGH 边界条件** — code 改造本质完全正确,所有 finding 都是 plan §当前进度 stale / §不变量边界 / §Phase H sequence ff-merge divergence / §已知踩坑 事实错误 / cosmetic 修法。

Phase F 全套验证:`make build` ✅ (build/vocabmaster 16.5MB) + `./build/vocabmaster --help` smoke ✅ + `make clean` ✅ (双删 build/ + 老根 vocabmaster) + `make build` 重 build ✅ + grep 0 残留 ✅ + `git status --short` 仅 3 文件 M。

## 变更内容

### Phase A — Makefile 改造(`Makefile`)

- `build:` target `go build -o vocabmaster .` → `go build -o build/vocabmaster .`(Go 1.4+ 自动建父目录,**不需** `mkdir -p build` prefix — reviewer-claude R1 MED-1 + reviewer-codex R1 LOW-1x 双方独立沙盒实测确认)
- `run:` target `./vocabmaster` → `./build/vocabmaster`
- `install:` target `sudo cp vocabmaster /usr/local/bin/...` → `sudo cp build/vocabmaster /usr/local/bin/...`(install 目标路径 `/usr/local/bin/vocabmaster` 不变,只是源 binary 路径变)
- `clean:` target `rm -f vocabmaster` → `rm -rf build/ vocabmaster`(双删 build/ + 老根 binary,与 .gitignore 删 vocabmaster 单数 entry 配套,避免老根 binary 变 untracked file — reviewer-codex R1 MED-2x finding 修法)
- **scope-side-fix**:`bash -lc '...'` → `bash -lc "..."`(L6/L9)— 原始 Makefile 在 GNU Make 3.81 + /bin/sh 解析单引号嵌套混乱(`grep '^go '` + `awk '{...}'` 内单引号与外层 `bash -lc '...'` 冲突),`make build` 直接挂。outer 改双引号修好(spike2 实测;reviewer-claude R1 INFO-1 + reviewer-codex R1 INFO-1x 双方 verified)。**Why scope-side-fix**:不修阻断 Phase F.1 `make build` 验证 + 用户跟 plan 收口后跑 `make build` 仍会撞同款 fail;1-char 修法 trivial 风险极低

### Phase B — Go 源码同步(`cmd/ internal/ tools/ data/`)

**empty** — grep 实测 0 处 hardcode binary path 引用(`cmd/root.go` 内 `vocabmaster` 字符串是 cobra app name + data dir `.vocabmaster` 不是 binary path;Go import path `github.com/vocabmaster/vocabmaster/...` 也不是 binary path)。Go 项目特殊性:src 不引用 binary 输出路径(由 Makefile 单一控制)。

### Phase C — `.gitignore` 改造

- L1 `vocabmaster` 单数 entry → `build/` 整目录 ignore(对齐应用打包 CLAUDE §新项目工程地基 §.gitignore canonical)
- `git check-ignore -v build/vocabmaster` ✅ 命中 `.gitignore:1:build/`
- 主 repo 老根 `./vocabmaster` 17MB(Apr 9 build)物理删走 Phase H.5.5(本 changelog 之前在 worktree 内 commit 时主 repo 老根 binary 仍在)+ `make clean` 一并删兜底

### Phase D — `README.md` narrative 更新

- L32 `# 可执行文件在 ./vocabmaster` → `# 可执行文件在 ./build/vocabmaster`(唯一 narrative 处;cobra 命令行示范 `vocabmaster study` etc 是 `make install` 后用户 PATH 上的 binary 不在 scope)

### Phase G — Step 1.5 deep-review × 1 轮 mixed kind fix loop

- **R1 (kind='mixed')**:reviewer-claude 8 finding(1 HIGH/3 MED/2 LOW/2 INFO)+ reviewer-codex 6 finding(0 HIGH/2 MED/3 LOW/1 INFO),共 10 ✅ 必修(1 双方独立 + lead 全自检 + 1 HIGH lead 现场 reproduce)+ 2 INFO verified;**0 反驳 + 0 新 HIGH 边界条件**
- **HIGH-1 (claude)**:Phase H.1-H.4 ff-merge sequence broken — reviewer 沙盒 reproduce + lead 现场 `/tmp/ff-test` reproduce `exit 128 fatal: Not possible to fast-forward`(main / worktree 都从 base 加 1 commit diverge)。**修法**:重写 §Phase H sequence,changelog commit 落 worktree branch(不落 main)→ 单一 ff-merge worktree(含 impl + changelog)→ main,fast-forward 成功
- **MED-1 (claude) + LOW-1x (codex) 双方独立**:Go `mkdir -p build` 不需要(plan §已知踩坑 + §Phase A.1 stale 描述)— 双方独立 `/tmp` 沙盒 + worktree 双场景实测 `go build -o build/sub2/deep/binary .` 多层不存在父目录都自动建。**修法**:删 §已知踩坑 第 1 条 mkdir 断言 + §影响面 spike A 改成 spike 结论 + §Phase A.1 「如需 mkdir -p build」括注删
- **MED-1x (codex)**:plan §当前进度 ⏳ entry stale(Phase A-F 已完成 plan 还说待跑)— update §当前进度 反映实际进度
- **MED-2 (claude)**:§D7 vs §Phase H 收口顺序自相矛盾 — 与 HIGH-1 同步重写对齐
- **MED-3 (claude)**:§Phase H.5 缺 `mkdir -p plans/`(`ls plans/: No such file or directory`,直接 `mv plan plans/<id>.md` 会变 rename 而非 move 进目录)— H.6 显式拆细 mkdir + mv + frontmatter + INDEX + commit
- **MED-2x (codex)**:`make clean` 不删老根 vocabmaster + plan §不变量 3 promise(下次 make clean 帮删)矛盾 — 扩 `clean: rm -rf build/ vocabmaster` 双删
- **LOW-1 (claude)**:plans/INDEX.md skeleton 模板未给(app archive_plan tool 4 列 `| 文件 | 状态 | 关联 changelog | 概要 |` vs 当前 changelog/INDEX.md 2 列) — H.6 明文写死用 4 列 canonical
- **LOW-2 (claude)**:§不变量 8 措辞含糊(「base 在其上」歧义)— 改成「base_commit `f2755b2` 即 main 当前 HEAD」
- **LOW-2x (codex)**:§F.5 grep gate false positive(install/uninstall 内 `/usr/local/bin/vocabmaster` 是 expected) — narrow scope 到 root `./vocabmaster`
- **LOW-3x (codex)**:`.deep-review-cache/` untracked 阻塞 H.1 worktree clean gate — 加 H.0 cleanup `.deep-review-cache/` step
- **INFO-1 (claude) + INFO-1x (codex)**:Makefile 双引号修复 verified ✅ — 不修
- **INFO-2 (claude)**:Phase H.8 `shutdown_baton_teammates` 建议补 lead callout — 加一句 "caller 必须仍是某 active team 的 lead,如全 archive tool 返 error 而非 silent success"

### Phase Spike — 2 spike 实测全 ✅

- spike1 Makefile build target:`go build -o build/vocabmaster .` Go 1.4+ 自动建父目录 ✅
- spike2 Makefile shell 引号 bug 修法:`bash -lc '...'` → `bash -lc "..."` /tmp 沙盒 + worktree 双实测 pass ✅

## 改动文件统计

- **3 impl**:`Makefile`(build/run/install/clean 4 target + shell 引号修)+ `.gitignore`(删 vocabmaster 加 build/)+ `README.md`(L32 narrative)
- **1 plan 主体**:`.claude/plans/build-dir-migration-20260526.md`(全程在 worktree 内被 `.gitignore` 忽略,Phase H.6 mv 入 `plans/build-dir-migration-20260526.md`)
- **1 changelog**:`changelog/CHANGELOG_1.md`(本文)
- **commits**:`8a2bde4 feat(build-dir): migrate build artifacts to build/ canonical (Phase A-D)`(Phase A-D 落地)+ 本 changelog commit + Phase H archive commit
- **build artifact**(被 `.gitignore` 整 build/ ignore 不入 git):`build/vocabmaster` (16.5MB)

详 [`plans/build-dir-migration-20260526.md`](../../plans/history/build-dir-migration-20260526.md)
