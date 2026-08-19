---
plan_id: "build-dir-migration-20260526"
created_at: "2026-05-26"
worktree_path: ".claude/worktrees/build-dir-migration-20260526"
status: "completed"
base_commit: "f2755b2e6e1e7a11d724479a7fc81678bcbd3f3e"
base_branch: "main"
final_commit: "e7570672a36ffeb801e8812f8c14d691123dcf79"
completed_at: "2026-05-26"
---

# Plan: vocabmaster 项目 build 产物迁移到 build/ 统一根出口

## 总目标

把 vocabmaster (Go CLI) 编译产物从**项目根** `./vocabmaster` 迁到 `./build/vocabmaster`:

- `Makefile build:` 输出路径 `vocabmaster` → `build/vocabmaster`
- `Makefile install:` 源路径同步更新
- `Makefile run:` 执行路径同步更新
- `Makefile clean:` 改成 `rm -rf build/`
- `.gitignore` 加 `build/` 删 `vocabmaster` 单数 entry
- README 内 `# 可执行文件在 ./vocabmaster` narrative 更新

**Why**:对齐 agent-deck 项目同款 build-dir-migration-20260526(commit `6a6903e`)+ cross-language `build/` canonical(应用打包 CLAUDE §src/build)+ 项目根更干净(binary 不与源码混在一起)+ `.gitignore` 整 build/ 比 single binary file ignore 更通用(将来若加 dev tool / 其他二进制不需改 .gitignore)。

**§Step 0 RFC 决策点 1:Go 项目是否适用 build/ canonical?** **决策 ✅ 适用,proceed with migration**(不 abort):
- Go 社区**没有**强制「binary 必须落 root」规则;`build/<name>` / `bin/<name>` / `_output/bin/<name>` 在 Go 生态都有先例(etcd 用 `bin/`,kubernetes 用 `_output/bin/`,vocabmaster 当前用根)
- `build/<name>` 与 `go install` 路径(`$GOPATH/bin/<name>`)不冲突 — `go install` 是用户安装路径,`build/<name>` 是 `make build` 产物路径,正交
- 迁移成本极低(3 文件:Makefile / .gitignore / README;Go 源码 0 binary path 引用)
- 跨语言 portfolio 一致性收益真实(user 多个 personal project 共用 build/ canonical)

**如何应用**(给下一会话):cold-start `Bash: cat <plan-abs-path>` 全文 → frontmatter 取 worktree_path → `EnterWorktree(path: <worktree_path>)` → 按 §下一会话第一步 接力

## 不变量

1. **build artifact 落 `build/vocabmaster`** — Makefile `build:` target 显式 `go build -o build/vocabmaster .`,不依赖 Go default 输出路径(go build 无 `-o` 时输出 `./vocabmaster` 与本 plan 目标冲突)
2. **`.gitignore` 整 build/ 忽略 + 删 `vocabmaster` 单数 entry** — `build/vocabmaster` 不入 git;老根 `./vocabmaster` 老 binary 一并物理删(避免 .gitignore 改后老根 binary 变成 tracked file)
3. **不向后兼容**(hard cutover) — 已装 `/usr/local/bin/vocabmaster` 用户重跑 `make install` 重装(install 路径不变,只是源 binary 路径变);老根 `./vocabmaster` 由 `make clean` 一并删(本 plan §Phase A.4 已扩 `clean:` target 为 `rm -rf build/ vocabmaster` 双删,与 reviewer-codex R1 MED-2x finding 对齐)
4. **不留兼容旧根 `./vocabmaster` fallback / 描述 / migration helper**(user 硬指令 + user CLAUDE §提示词资产维护 约束 2「当前事实」)
5. **`make build` + `./build/vocabmaster --help` smoke test pass** 是收口前置条件
6. **每个 fix 必有同步实测** — 改 Makefile 同步跑 `make build` verify 产物 actually 落 `build/vocabmaster`
7. **doc 内 `./vocabmaster` 提及全替换**(README L32 narrative;cobra 命令行示范 `vocabmaster study` 是用户 PATH 上的 binary 不在 scope)
8. **base_commit `f2755b2 chore(cleanup)` 即 main 当前 HEAD**(LOW-2 claude R1 fix:消歧义,frontmatter `base_commit: f2755b2...` 与 main HEAD 完全相同 — 不是「在该 commit 之上」)
9. **扫描范围**:`Makefile + .gitignore + README.md + cmd/ + internal/ + tools/ + data/`(顶层 self-describe 文件 + Go 源码 + 工具脚本);**不扫**:`changelog/ reviews/ conventions/ plans/`(历史归档保持当时事实) + `.claude/ vendor/ node_modules/ build/`(临时 / 第三方)

## 设计决策(不再争论)

### D1: 子目录拆分粒度 — 平铺 `build/vocabmaster`(单 binary)

vocabmaster 当前是**单 binary** Go CLI app(`main.go` + `cmd/` cobra subcommands)。**不需** `build/main/ build/cli/` 子目录拆分。

未来若 vocabmaster 加 server / dev tool 等额外 binary,按 `build/<binary-name>` canonical 平铺加(如 `build/vocabmaster-server`);不需改 Makefile build target 结构 / 不需改 .gitignore(整 build/ 已 ignore)。

### D2: .gitignore = 整 build/ + 删 `vocabmaster` 单数 entry

`.gitignore`:
- 删 L1 `vocabmaster`(单 binary 文件 ignore,不再需要因为 binary 不在 root)
- 加 `build/`(整 build/ canonical,应用打包 CLAUDE §新项目工程地基 §.gitignore 必备条目)

**Why**:未来加新 binary `build/<other-name>` 不需同步改 .gitignore;比单文件 ignore 更通用 + canonical 对齐 agent-deck plan。

### D3: 不向后兼容(hard cutover)

老根 `./vocabmaster` binary 用户重跑 `make build` 自动落新位置 `./build/vocabmaster`;老根 binary 不会被自动删,Phase C 通过 `git status` 实测后让 Phase F.0 手动 `rm -f vocabmaster` 删干净。

已装 `/usr/local/bin/vocabmaster` 不动(install 路径不变,只是 install 源路径变)。

### D4: changelog X 接起算

收口时按 `ls changelog/CHANGELOG_*.md | max + 1` fail-fast 算 X。当前 `changelog/INDEX.md` 还是 skeleton(只有 placeholder `[CHANGELOG_1.md](CHANGELOG_1.md)` 没真文件),所以本 plan 是 CHANGELOG_1(真第一条 changelog)。

### D5: Step 1.5 deep-review × 1 轮(Go 单 binary 改造,1 轮足够)

vocabmaster 改造 scope 极小(3 文件,无逻辑变更,纯 build path 替换),deep-review × 1 轮 mixed kind(plan + impl 一起评)足够;不需 R2/R3。如果 R1 出现 HIGH-1 反驳轮触发再考虑 R2。

### D6: Go 源码 0 binary path 引用 — Phase B empty

grep 实测 `cmd/ internal/ tools/ data/` 内 `vocabmaster` 字符串引用全部是:
- `Use: "vocabmaster"`(cobra command name,与 binary path 无关)
- `github.com/vocabmaster/vocabmaster/...`(Go module import path)
- `dataDir = filepath.Join(home, ".vocabmaster")`(数据目录,不是 binary)

**结论**:Phase B(src 注释)为空 — vocabmaster 没有任何 Go 源码内 hardcode binary path 需改。

### D7: G-manual 收口路径(R1 HIGH-1 fix:changelog 落 worktree 内避 ff-merge divergence)

按 user 授权(详 §下一会话第一步 user 授权 callout)走 G-manual 收口,**核心 invariant**:**所有 commit 必须落在 worktree branch,然后单一 ff-merge 把 worktree(含 impl + changelog)merge 回 main**。R1 reviewer-claude HIGH-1 沙盒 reproduce + lead 现场 reproduce(`exit 128 fatal: Not possible to fast-forward`)证明:如果 H.1 ExitWorktree 后在 main repo 上跑 changelog commit,worktree branch 与 main 都从 base 各加 1 commit → diverge,H.4 ff-merge 必 fail。

正确序列:① commit Phase A-F impl(in worktree)→ ② cleanup `.deep-review-cache/`(worktree clean gate 前置)→ ③ 算 X + 写 changelog + commit(in worktree)→ ④ ExitWorktree(action: keep)→ ⑤ ff-merge worktree branch → main(纯 fast-forward)→ ⑥ rm 老根 `./vocabmaster` 物理删 → ⑦ mkdir -p plans/ + mv plan + frontmatter update + 写 plans/INDEX.md + commit(in main)→ ⑧ git worktree remove + branch -D → ⑨ shutdown_baton_teammates。

详 §Phase H 7 步分解。

## 影响面 spike

### A. Makefile build target(`Makefile`)

实测项:
- L6 `go build -o vocabmaster .` → `go build -o build/vocabmaster .` 是否让 `make build` 产物落 `build/vocabmaster`?(Phase F.0 实测 ✅)
- Go `go build -o build/<name> .` 父目录创建行为:**自动创建**父目录(Go 1.4+ 稳定行为,reviewer-claude R1 MED-1 + reviewer-codex R1 LOW-1x 双方独立实测确认 — `/tmp` 沙盒 + 当前 worktree 双场景都 pass,**不需** `mkdir -p build` prefix)

### B. Go 源码内 hardcode binary path

`spike pre-check` 已 grep 实测:`cmd/ internal/ tools/ data/` 内 0 处 hardcode binary path(详 §设计决策 D6)→ **Phase B empty**

### C. .gitignore(`vocabmaster` 单数 entry → `build/`)

- 删 L1 `vocabmaster`(单 binary ignore)
- 加 `build/`(canonical)
- 同步 `rm -f vocabmaster`(物理删老根 binary,避免变 tracked file)

### D. README.md narrative

`spike pre-check` 已 grep 实测:
- L32 `# 可执行文件在 ./vocabmaster` → `# 可执行文件在 ./build/vocabmaster`(narrative 唯一需改处)
- L39 起 `vocabmaster study` etc 是用户 PATH 上的 binary 调用示范,**不在 scope**(install 后用户从 PATH 跑,不是 `./vocabmaster`)

### E. tools/ 内 python 脚本

`tools/extract_words.py` — grep 实测无 binary path 引用(纯 Python 词典提取脚本),**不动**。

## 步骤 checklist

### Phase Spike: spike 实测(在 worktree 内跑)

- [x] **spike1 Makefile build target**(实测 ✅):改 `Makefile` `go build -o vocabmaster .` → `go build -o build/vocabmaster .` → `make build` 实测产物落 `build/vocabmaster`(`ls build/` 确认 + `./build/vocabmaster --help` smoke test 验证可执行)
  - **结论 ✅**:`/tmp/go-build-test` 干净 dir 实测 `go build -o build/test .` 自动创建 `build/` 父目录(Go 1.4+ 稳定行为,**不需** `mkdir -p build` prefix);worktree 内 `make clean && make build` 也 verify 同款行为。reviewer-claude R1 MED-1 + reviewer-codex R1 LOW-1x 双方独立实测确认。
- [x] **spike2 Makefile shell 引号 bug 修法**(spike inline 完成 ✅):原 `bash -lc '...'` 在 GNU Make 3.81 + /bin/sh 解析单引号嵌套混乱(`grep '^go '` + `awk '{...}'` 内单引号与外层 `bash -lc '...'` 冲突)。**修法**:outer 改双引号 `bash -lc "..."`,GOENV 内 `$$HOME / $$(...) / $$2` 经 make 转义后变 `$HOME / $(...) / $2`,在双引号内 `$HOME` 被 shell 提前展开,`$(...)` 命令替换正常,awk 单引号 body 不受外层 double-quote 影响(awk 拿到的是字面 string)— /tmp 沙盒 + worktree 内双实测 pass。reviewer-claude R1 INFO-1 + reviewer-codex R1 INFO-1x 双方 verified。

spike 结论 inline 到 §设计决策 + §已知踩坑;残留风险列表入 §已知踩坑。

### Phase A: Makefile 改造(已完成 ✅ — checklist for trace)

- [x] A.1 改 `Makefile:6` `go build -o vocabmaster .` → `go build -o build/vocabmaster .`(spike1 验证;Go 1.4+ 自动建父目录,**不需** `mkdir -p build` prefix)
- [x] A.2 改 `Makefile:12` `./vocabmaster` → `./build/vocabmaster`(run target)
- [x] A.3 改 `Makefile:15` `sudo cp vocabmaster /usr/local/bin/vocabmaster` → `sudo cp build/vocabmaster /usr/local/bin/vocabmaster`(install target)
- [x] A.4 改 `Makefile:23` `rm -f vocabmaster` → `rm -rf build/ vocabmaster`(clean target — 删整 build/ 同时删老根 vocabmaster,与 §不变量 3 promise 对齐 + reviewer-codex R1 MED-2x finding 修法)
- [x] A.5 改 `Makefile:6/9` `bash -lc '...'` → `bash -lc "..."`(scope-side-fix:Make 3.81 单引号嵌套解析 bug,详 spike2 结论)
- [x] A.6 `make build` verify 产物落 `build/vocabmaster`(16.5MB) + `./build/vocabmaster --help` smoke test ✅

### Phase B: Go 源码同步 — empty(详 §设计决策 D6)

- [x] B.1 grep `cmd/ internal/ tools/ data/` 内 `./vocabmaster` / `/usr/local/bin/vocabmaster` hardcode 实测 0 残留(实测 ✅ exit code 1)

### Phase C: .gitignore 改造(已完成 ✅)

- [x] C.1 改 `.gitignore:1` 删 `vocabmaster` 单数 entry + 加 `build/`
- [x] C.2 `git check-ignore -v build/test.txt` 实测 ✅ 命中 `.gitignore:1:build/`
- [x] C.3 老根 binary 物理删 — 推迟到 §Phase H ff-merge 后(主 repo 内执行,worktree 内无老根 binary 文件;详 §Phase H 序列)
- [x] C.4 `git status --short` 干净 — 仅 `M .gitignore + M Makefile + M README.md`,无 `?? build/`(整 build/ ignore)

### Phase D: README.md narrative 更新(已完成 ✅)

- [x] D.1 改 `README.md:32` `# 可执行文件在 ./vocabmaster` → `# 可执行文件在 ./build/vocabmaster`
- [x] D.2 grep 0 残留 `./vocabmaster`(scope = §不变量 9 列定的 path;`vocabmaster study` etc 用户 PATH 命令不在 scope)

### Phase E: 其他 doc(empty — vocabmaster 项目无 active CLAUDE.md / 工程文档需同步)

- vocabmaster 项目根**无** CLAUDE.md(详 §影响面 spike — 与 agent-deck 不同,vocabmaster 没有项目级 CLAUDE.md)
- `conventions/` `changelog/` `reviews/` 全是 skeleton 没真内容引用 binary path
- **Phase E empty**

### Phase F: 全套验证(已完成 ✅)

- [x] F.0 **§不变量 6 enforcement**:Phase A-D 任何后续补 fix 改完立即跑对应命令实测(改 Makefile 跑 `make build`;改 .gitignore 跑 `git check-ignore -v`;改 README 只 grep)
- [x] F.1 `make build` ✅ verify 产物落 `build/vocabmaster` (16.5MB)
- [x] F.2 `./build/vocabmaster --help` smoke test ✅ verify 可执行 + 输出 cobra help
- [x] F.3 `make clean` ✅ verify `build/ + 老根 vocabmaster` 整删
- [x] F.4 `make build` 重跑 ✅ verify clean 后能重 build(`build/` 不存在时 Go 自动创建并落产物)
- [x] F.5 grep 0 残留 ✅(R1 LOW-2x fix:**narrow scope 到 root `./vocabmaster`**;install/uninstall 内 `/usr/local/bin/vocabmaster` 是 expected install target 不算残留):
  - 精确 grep:`git grep -nE '\./vocabmaster' -- Makefile README.md .gitignore cmd internal tools data` → **0 输出**(exit 1)— 无残留根 `./vocabmaster` 引用
  - 宽松 grep(含 install path,sanity check):`git grep -nE '(\\./|/usr/local/bin/)vocabmaster' -- Makefile README.md .gitignore` 输出 `Makefile:15 sudo cp build/vocabmaster /usr/local/bin/vocabmaster` + `Makefile:16` echo + `Makefile:19` `sudo rm -f /usr/local/bin/vocabmaster` — 都是 expected install/uninstall target,不是残留
- [x] F.6 `git status --short` 干净 — `M .gitignore + M Makefile + M README.md` 三文件 + 无 `?? build/`(整 build/ ignore)

### Phase G: Step 1.5 deep-review × 1 轮 mixed kind(plan + impl 一起评)— **已完成 R1 ✅**

- [x] G.1 invoke deep-review SKILL `kind='mixed'`,scope = Phase A-D 改动文件 + 本 plan 文件(R1 invocation `bc82bb6e`,team `build-dir-vocabmaster-r1` id `66d26c83`)
- [x] G.2 R1 finding 双 reviewer reply 收齐:reviewer-claude 8 finding (1 HIGH/3 MED/2 LOW/2 INFO) + reviewer-codex 6 finding (0 HIGH/2 MED/3 LOW/1 INFO),共 10 ✅ 必修(1 双方独立 + lead 自检 + 1 lead 现场 reproduce HIGH-1) + 2 INFO verified;**0 反驳 + 0 反驳轮**
- [x] G.3 R1 fix 已落本 plan + Makefile clean target 扩;R2 评估**不需**(0 反驳 + 0 新 HIGH 边界条件 + lead 已现场 reproduce HIGH-1 ff-merge fail);→ 直接 Phase H 收口

### Phase H: 收口(G-manual 路径,user 已授权;R1 HIGH-1 重写 sequence avoid ff-merge divergence)

> **R1 HIGH-1 修法核心**:所有 commit 落 worktree branch(包括 changelog),然后单一 ff-merge worktree → main。反例(旧 sequence H.1→H.2/3→H.4)让 changelog 落 main / impl 落 worktree → 二者从 base 各加 1 commit diverge → `git merge --ff-only` 必报 `fatal: Not possible to fast-forward` exit 128(reviewer-claude /tmp 沙盒 reproduce + lead 现场 `/tmp/ff-test` reproduce 双 verify)。

- [ ] H.-1 **commit Phase A-G tracked impl 改动**(in worktree:Makefile + .gitignore + README.md)
  - **Plan 文件本身不在 H.-1 commit 范围**:`.claude/plans/` 已被 `.gitignore` 忽略,plan 文件入库职责在 H.5 mv plan → `plans/<id>.md` + commit(in main)
- [ ] H.0 **cleanup `.deep-review-cache/`** (in worktree):`rm -rf .deep-review-cache/`(R1 reviewer-codex LOW-3x finding:G.1 SKILL 创建 cache 阻塞 H.1 worktree clean gate;SKILL Step 7 也会 cleanup 但 G-manual 路径 lead 主动 cleanup 更可靠)
- [ ] H.1 worktree clean gate(`git status --short` 空,无 M / ?? / A 任何 entry,包括 `?? .deep-review-cache/` 已删)
- [ ] H.2.0 **算 X**(in worktree):`ls changelog/CHANGELOG_*.md 2>/dev/null | grep -oE '[0-9]+' | sort -n | tail -1 | (read max; echo $((${max:-0} + 1)))`;当前 `changelog/INDEX.md` 是 skeleton 无真文件 → X=1
- [ ] H.2/3 **写 `changelog/CHANGELOG_<X>.md` + sync `changelog/INDEX.md` + commit**(in worktree — **关键:必须在 worktree 内 commit,不要 ExitWorktree 后在 main 上 commit,否则下 H.4 ff-merge fail**)
- [ ] H.3 worktree final clean gate(再次 `git status --short` 空,确认 changelog commit 已落 worktree branch)
- [ ] H.4 ExitWorktree(action: "keep") — cwd 切回 main repo
- [ ] H.5 `git -C <main-repo> merge --ff-only worktree-build-dir-migration-20260526`(main 在 base,worktree branch 在 base + 2 commits 即 impl + changelog,纯 fast-forward 成功)
- [ ] H.5.5 **`rm -f vocabmaster`**(主 repo 物理删老根 17MB binary;.gitignore 已删 `vocabmaster` 单数 entry 后老根 binary 失去 ignore 兜底,必须物理删避免变 untracked file)
- [ ] H.6 **`mkdir -p plans/`**(主 repo 内,plans/ 当前不存在 — R1 reviewer-claude MED-3 finding `ls plans/: No such file or directory`;不先 mkdir 直接 `mv .claude/plans/<id>.md plans/<id>.md` 会变 rename 成 `plans` 单文件而非 move 进目录)→ `mv .claude/plans/build-dir-migration-20260526.md plans/build-dir-migration-20260526.md` → 编辑 frontmatter `status=completed + final_commit=<H.5 ff-merge 后 main HEAD SHA> + completed_at=<ISO date>` → 写 `plans/INDEX.md` **4 列 canonical** `| 文件 | 状态 | 关联 changelog | 概要 |`(对齐应用约定 archive_plan tool 4 列 canonical — R1 LOW-1 finding) → `git add plans/ && git commit -m "..."`(in main)
- [ ] H.7 `git worktree remove .claude/worktrees/build-dir-migration-20260526 && git branch -D worktree-build-dir-migration-20260526`
- [ ] H.8 **`shutdown_baton_teammates`**(**mandatory** — §设计决策 D5 deep-review × 1 轮 reviewer pair 必有 dormant 残留;G-manual 路径绕过 archive_plan tool baton-cleanup phase 1 → 必须手动调 escape hatch)
  - **R1 INFO-2 callout**:`shutdown_baton_teammates` 调用前 caller 必须仍是某 active team 的 lead;如所有 team 已 archive,本 tool 返 **error 而非 silent success**(应用约定 SSOT)。Phase H 序列中 caller 仍是 `build-dir-vocabmaster-r1` team 的 lead,正常调用应成功

### Phase I: Post-archive fs 真验证

- [ ] I.1 archive 文件真存在(`ls -la plans/build-dir-migration-20260526.md`)
- [ ] I.2 git commit 含 archive(`git log --oneline -3` 看到 H.6 commit)
- [ ] I.3 `plans/INDEX.md` append 行存在
- [ ] I.3.5 frontmatter `status=completed + final_commit + completed_at` 真写入
- [ ] I.4 `git log --oneline --follow plans/build-dir-migration-20260526.md` 看到 H.6 移动事件
- [ ] I.5 worktree dir 真删(`ls .claude/worktrees/` 应无 `build-dir-migration-20260526/`)+ branch 真删(`git branch -a | grep worktree-build` 0 输出)
- [ ] I.6 通知 user 收口完成 + 走 user 自验证步骤(可选:`make install` 重装 `/usr/local/bin/vocabmaster` 实测新 build 路径生效)

## 当前进度

- ✅ §Step 0 RFC 完成(决策:Go 项目适用 build/ canonical,proceed with migration)
- ✅ §Step 1 plan v1 写完(2026-05-26 本会话,base_commit `f2755b2`)
- ✅ §Step 0.5 spike + §Step 2 EnterWorktree(MCP enter_worktree + builtin EnterWorktree path: 双步,worktree `worktree-build-dir-migration-20260526` base 在 `f2755b2`)
- ✅ §Phase A-D 实施完成:Makefile build/run/install/clean 4 target 改 + Makefile shell 引号 bug inline 修(scope-side-fix)+ .gitignore 删 vocabmaster 加 build/ + README L32 narrative 更新
- ✅ §Phase F 全套验证完成:F.1 make build ✅ + F.2 smoke test ✅ + F.3 make clean ✅ + F.4 重 build ✅ + F.5 grep 0 残留 ✅ + F.6 git status clean ✅
- ✅ §Step 1.5 / Phase G deep-review R1 mixed kind 完成(2026-05-26 本会话,R1 invocation `bc82bb6e`,team `build-dir-vocabmaster-r1` id `66d26c83`):reviewer-claude 8 finding + reviewer-codex 6 finding 全 ✅ 必修(0 反驳 + 1 双方独立强冗余 + lead 全自检 + 1 HIGH lead 现场 reproduce);finding fix 已落本 plan + Makefile clean target 扩
- ✅ R2 评估**不需**(0 反驳 + 0 新 HIGH 边界条件 + R1 finding 全 ✅ 落地)→ 直接 Phase H 收口
- ⏳ **§Phase H 收口(重写 sequence)**:H.-1 commit Phase A-G impl → H.0 cleanup .deep-review-cache/ → H.1 worktree clean gate → H.2.0 算 X → H.2/3 写 changelog + commit(**in worktree**) → H.3 worktree final clean → H.4 ExitWorktree → H.5 ff-merge → H.5.5 rm 老根 vocabmaster → H.6 mkdir plans/ + mv plan + 4 列 INDEX + commit → H.7 worktree remove + branch -D → H.8 shutdown_baton_teammates
- ⏳ §Phase I post-archive 真验证

## 下一会话第一步(cold-start 接力指令)

> ⚠️ 本 plan 由首会话写出。新会话 cold-start 时按 §当前进度 接力 — **找最近一个 ⏳ entry 就是接力起点**(不要从头跑)。

> 📜 **2026-05-26 user 授权 callout**(context: plan 起头 user 指令明示授权):
> - **lead 全权决定 hand off 时机**(不需逐 phase 请求 user 确认)
> - **隐含 G-manual 路径授权**(若 archive_plan tool 撞 precheck fail 走 §Phase H 5 步手工归档兜底)
> - **user 离开期间允许全自动推进**
> - **若决策「Go 项目沿用根 `./vocabmaster` 不改 `build/`」(理由:Go 社区惯例 + 与 cross-language canonical 冲突权衡)允许 abort plan + 写 N/A 报告通知 user** — 本 plan §Step 0 RFC 已决策 **proceed**(不 abort),如下一会话 lead 推翻这个决策需在 plan 内显式 commit "abort RFC 决策推翻" reason 并走 abort 路径

### Cold-start 5 步(标准接力流程)

1. `Bash: cat .claude/plans/build-dir-migration-20260526.md`(全文)
2. 读 §当前进度,找最近一个 ⏳ entry — 就是接力起点
3. EnterWorktree(builtin) `path: .claude/worktrees/build-dir-migration-20260526`(避 v2.1.112 stale base bug,worktree 已存在不要再 git worktree add)
4. `git log --oneline -3` 自检 HEAD 含本 plan 的 commit 历史(或 base_commit `f2755b2`)
5. 按 §当前进度 ⏳ 起点对应 §Phase 章节实施,每完成一 Phase / Step 在本 plan 文件 `- [ ]` 打勾 + commit 进度

## 已知踩坑 / 风险

- **Go `go build -o build/vocabmaster .` 父目录创建**:**Go 1.4+ 自动创建父目录**(reviewer-claude R1 MED-1 + reviewer-codex R1 LOW-1x 双方独立沙盒实测确认 — `/tmp` 干净 dir 跑 `go build -o build/sub2/deep/binary .` 多层不存在父目录都自动建)。Makefile **不需** `mkdir -p build &&` prefix
- **老根 `./vocabmaster` binary 物理残留**:`.gitignore` 删 `vocabmaster` 单数 entry 后,主 repo 老根 binary 17MB(Apr 9 build)会变 untracked。**两路径**:① `make clean` 一并删(本 plan §Phase A.4 已扩 `clean: rm -rf build/ vocabmaster`)② Phase H.5.5 主 repo 内手工 `rm -f vocabmaster`
- **`make install` 后 `/usr/local/bin/vocabmaster` 不会自动更新**:本 plan 不动 install 路径(目标路径不变),用户需重跑 `make install` 让新 `build/vocabmaster` 同步到 `/usr/local/bin/`;不属于本 plan scope(install 是用户主动行为)
- **Step 1.5 deep-review 必有 dormant teammate 残留**:G-manual 路径绕过 archive_plan tool baton-cleanup phase 1 → Phase H.8 必须手动调 `shutdown_baton_teammates` escape hatch(R1 INFO-2:caller 必须仍是某 active team 的 lead,如全 archive tool 返 error 而非 silent success)
- **`.deep-review-cache/` 阻塞 worktree clean gate**(R1 LOW-3x):G.1 SKILL 创建 cache 在 worktree 内是 untracked,Phase H.0 必须 `rm -rf .deep-review-cache/` 否则 H.1 worktree clean gate fail
- **HIGH-1 ff-merge divergence 已规避**(R1 HIGH-1 修法落地):Phase H 序列把 changelog commit 落 worktree branch 而非 main,然后单一 ff-merge worktree → main 把 impl + changelog 一并 merge 进 main。反例 sequence(H.1 ExitWorktree → main 上 commit changelog → H.4 ff-merge)已 reviewer-claude 沙盒 + lead 现场 reproduce `exit 128 fatal: Not possible to fast-forward`
- **spike 失败回滚**:若 Phase A spike1 `make build` 失败 → revert Makefile 改回 `go build -o vocabmaster .` + 在 plan §已知踩坑 标注「Go build -o build/<name> 行为与预期不符」+ plan 重写 D1 设计决策(目前看高度不可能,Go 单 binary 输出极简单)

## 关联

- **触发**:user 指令「agent-deck 项目刚完成 build-dir-migration,你来对 vocabmaster 做适配同款改造」
- **上游 plan(参考样本)**:agent-deck `build-dir-migration-20260526`(final commit `6a6903e`)
  - agent-deck plan 文件:`../agent-deck/ref/plans/build-dir-migration-20260526.md`
  - agent-deck changelog:`../agent-deck/ref/changelogs/CHANGELOG_154.md`
- **changelog 关联**:本 plan 完成后写 `changelog/CHANGELOG_1.md`(本仓库第一条 changelog;`changelog/INDEX.md` 当前 skeleton)
