# AIGuard 操作手册

> 版本：0.2.0 | 最后更新：2026-03-29

---

## 目录

1. [产品简介](#1-产品简介)
2. [环境要求](#2-环境要求)
3. [安装与构建](#3-安装与构建)
4. [配置说明](#4-配置说明)
5. [界面操作指南](#5-界面操作指南)
6. [审计流程详解](#6-审计流程详解)
7. [报告解读](#7-报告解读)
8. [自定义审计规则](#8-自定义审计规则)
9. [工作区管理](#9-工作区管理)
10. [常见问题与排查](#10-常见问题与排查)

---

## 1. 产品简介

AIGuard（AI代码监视）是一款面向团队内部使用的**本地桌面程序**，用于对 GitHub PR / GitLab MR / 本地仓库的代码改动进行 AI 辅助审计。

### 核心能力

- **多平台支持**：GitHub PR、GitLab MR、本地仓库三种模式
- **自动识别**：粘贴 MR/PR 链接即可自动推导仓库地址和分支
- **智能审计**：项目画像 → 规则预扫 → LLM 分片审计 → 裁决去重，四阶段流水线
- **多格式报告**：HTML（可视化）、Markdown（易读）、JSON（程序集成）
- **历史对比**：自动对比同分支上次审计，标注新增 / 已修复 / 仍存在
- **隐私安全**：全程本地运行，敏感信息（密钥、Token）在发送 LLM 前自动脱敏

### 技术架构

| 层级 | 技术栈 |
|------|--------|
| 桌面框架 | Wails 2.11 |
| 后端 | Go 1.22 |
| 前端 | Vue 3 + TypeScript + Vite |
| AI 接口 | OpenAI-compatible API |

---

## 2. 环境要求

### 运行时依赖

| 组件 | 版本要求 | 说明 |
|------|---------|------|
| Git | 最新版 | 仓库克隆、diff 生成 |
| OpenAI-compatible 模型服务 | - | 本地或远程均可 |

### 构建时额外依赖

| 组件 | 版本要求 |
|------|---------|
| Go | 1.21+（macOS 15 建议 1.23.3+） |
| Node.js | 18+ |
| Wails CLI | 最新版（`go install github.com/wailsapp/wails/v2/cmd/wails@latest`） |

### 平台特定依赖

- **Windows**：WebView2 Runtime
- **macOS**：Xcode Command Line Tools
- **Linux (Ubuntu 24.04+)**：`libwebkit2gtk-4.1-dev`、`libgtk-3-dev`、`pkg-config`

---

## 3. 安装与构建

### 3.1 快速开始（开发模式）

```bash
# 1. 安装前端依赖
cd frontend && npm install && cd ..

# 2. 准备配置文件
cp examples/config.production.yaml config.yaml
# 编辑 config.yaml，填入你的 LLM 服务地址和密钥

# 3. 启动开发模式
wails dev
```

### 3.2 构建生产二进制

```bash
# 通用构建
wails build
```

### 3.3 平台专用构建脚本

**Windows (PowerShell)：**

```powershell
./scripts/preflight.ps1          # 依赖检查
./scripts/build_windows.ps1      # 构建 EXE
./scripts/build_windows.ps1 -Nsis  # 带 NSIS 安装包
./scripts/package_release.ps1    # 打包发布
```

**macOS：**

```bash
bash ./scripts/preflight.sh       # 依赖检查
bash ./scripts/build_macos.sh     # 构建 universal app
bash ./scripts/package_release.sh # 打包发布
```

**Linux：**

```bash
bash ./scripts/preflight.sh             # 依赖检查
bash ./scripts/build_linux.sh           # 构建二进制
bash ./scripts/package_release_linux.sh # 打包发布
```

> 提示：发版前建议先运行 `wails doctor` 确认环境完整。

---

## 4. 配置说明

配置文件为 YAML 格式，默认路径为项目根目录的 `config.yaml`。

### 4.1 LLM 服务配置

```yaml
openai:
  base_url: "http://127.0.0.1:8000/v1"              # LLM 服务地址
  api_key: "your-api-key"                            # API 密钥
  default_model: "Qwen/Qwen3-Coder-30B-A3B-Instruct" # 模型名称
  proxy:
    enabled: false    # 是否启用代理
    url: ""           # 代理地址
```

### 4.2 运行时配置

```yaml
runtime:
  request_timeout_sec: 180     # 单次请求超时（秒）
  concurrency: 4               # 并发审计包数量
  max_retry: 2                 # 失败重试次数
  safe_input_tokens: 160000    # 安全输入 token 上限
  reserved_output_tokens: 12000 # 预留输出 token
  log_level: "info"            # 日志级别（debug/info/warn/error）
```

### 4.3 审计配置

```yaml
review:
  workspace_dir: "./workspace"         # 工作区目录
  diff_strategy: "merge_base"          # Diff 策略
  max_changed_files: 200               # 最大变更文件数
  max_hunks_per_file: 40               # 每文件最大 hunks
  export_formats: ["html", "md", "json"] # 导出格式
  enable_project_brief: true           # 启用项目画像
  enable_prescan: true                 # 启用规则预扫
  redact_secrets_before_llm: true      # LLM 发送前脱敏
  code_browse_base_url: ""             # 代码在线浏览基址（用于报告中的跳转链接）
```

### 4.4 规则配置

```yaml
rules:
  custom_rule_file: "./examples/AIGUARD.md"  # 自定义审计规则文件
  ignore:                                     # 忽略的文件/目录（glob 模式）
    - "node_modules/**"
    - "dist/**"
    - "vendor/**"
    - "*.min.js"
    - "*.lock"
```

### 4.5 Git 配置

```yaml
git:
  preferred_protocol: "ssh"   # 优选协议：ssh 或 https
  gitlab:
    ssh:
      host: ""                # 自定义 SSH 主机（留空则使用 MR 链接中的主机）
      port: ""                # 自定义 SSH 端口
      user: "git"
    https:
      scheme: "https"
      host: ""
      port: ""
  github:
    ssh:
      host: ""
      port: ""
      user: "git"
    https:
      scheme: "https"
      host: ""
      port: ""
```

### 4.6 环境变量覆盖

所有配置均可通过环境变量覆盖（优先级高于配置文件）：

| 环境变量 | 对应配置 |
|---------|---------|
| `AIGUARD_OPENAI_BASE_URL` 或 `OPENAI_BASE_URL` | openai.base_url |
| `AIGUARD_OPENAI_API_KEY` 或 `OPENAI_API_KEY` | openai.api_key |
| `AIGUARD_OPENAI_DEFAULT_MODEL` 或 `OPENAI_DEFAULT_MODEL` | openai.default_model |
| `AIGUARD_WORKSPACE_DIR` | review.workspace_dir |
| `AIGUARD_LOG_LEVEL` | runtime.log_level |
| `AIGUARD_CONCURRENCY` | runtime.concurrency |

---

## 5. 界面操作指南

### 5.1 界面布局

应用界面分为**左右两栏**：

```
+---------------------------+-------------------------------+
|        监视配置            |         执行状态               |
|  (MR/PR 链接、仓库地址、   |   (进度条、阶段、问题统计)      |
|   分支、配置路径、操作按钮)  |                               |
+---------------------------+-------------------------------+
|        历史记录            |         运行日志               |
|  (历史审计报告列表)         |   (实时日志输出)               |
+---------------------------+-------------------------------+
|                           |         审计报告               |
|                           |   (问题卡片、健康度评分)        |
+---------------------------+-------------------------------+
```

### 5.2 操作步骤

#### 第一步：填写审计信息

1. **MR/PR 链接**（推荐）：粘贴 GitLab MR 或 GitHub PR 的完整 URL
   - 示例：`https://gitlab.example.com/group/project/-/merge_requests/123`
   - 示例：`https://github.com/org/repo/pull/42`
   - 仓库地址会自动推导填入

2. **仓库地址**：自动推导后仍可手动编辑覆盖
   - 支持 SSH 格式：`git@gitlab.example.com:group/project.git`
   - 支持 HTTPS 格式：`https://gitlab.example.com/group/project.git`

3. **本地仓库路径**（可选）：如果代码已在本地，可直接填写本地路径

4. **配置文件路径**：默认为 `config.yaml`

5. **工作区路径**：默认为 `./workspace`

#### 第二步：拉取代码

点击 **「拉取代码」** 按钮：
- 自动克隆/拉取仓库到工作区
- 自动检测并推荐源分支和目标分支
- 分支列表填充到下拉框

> 如果分支自动识别不准确，可手动修改。

#### 第三步：开始监视

确认分支无误后，点击 **「开始监视」**：
- 右侧「执行状态」面板实时显示进度
- 审计完成后自动生成报告

#### 第四步：查看报告

审计完成后：
- 执行状态面板右侧出现 **「查看报告」** 和 **「打开目录」** 按钮
- 「查看报告」：在浏览器中打开 HTML 报告
- 「打开目录」：在文件管理器中打开报告所在目录
- 历史记录列表中也可以打开之前的报告

### 5.3 功能按钮说明

| 按钮 | 功能 |
|------|------|
| 拉取代码 | 克隆/拉取仓库，识别分支 |
| 开始监视 | 启动完整审计流程 |
| 取消任务 | 中止正在运行的审计任务 |
| 查看日志 | 展开/收起实时运行日志 |
| 清理缓存 | 清除工作区中的仓库缓存、报告和日志 |
| 刷新历史 | 重新加载历史报告列表 |

---

## 6. 审计流程详解

完整审计流程包含以下阶段：

```
初始化 (5%) → 同步代码 (15%) → 项目画像 (30%) → 规则预扫 (45%)
    → 构建审计包 (60%) → AI审计 (68%-84%) → 结果裁决 (88%)
    → 生成报告 (95%) → 完成 (100%)
```

### 阶段说明

| 阶段 | 说明 |
|------|------|
| **初始化** | 准备工作区、加载配置 |
| **同步代码** | 准备 Git worktree，计算 merge-base diff |
| **项目画像** | 分析项目技术栈、模块边界、敏感区域、关键入口 |
| **规则预扫** | 确定性规则扫描（不依赖 LLM） |
| **构建审计包** | 按文件切分 diff + 上下文，生成 Review Pack |
| **AI审计** | 并发调用 LLM 对每个 Pack 进行代码审计 |
| **结果裁决** | 合并、去重、统一严重级别、跨文件整合 |
| **生成报告** | 输出 HTML / Markdown / JSON 三种格式 |

### 审计维度

审计覆盖以下维度：

- **安全**：SQL 注入、XSS、CSRF、敏感信息泄露、权限绕过
- **性能**：资源未关闭、事务边界、缓存缺失、死锁风险
- **健壮性**：异常处理、边界条件、空值风险、并发问题
- **规范**：命名、复杂度、重复代码、SOLID 原则
- **框架**：React/Vue/Spring/Nest/FastAPI 等惯用法
- **架构**：跨文件代码重复、逻辑矛盾、设计缺陷、接口契约违反

---

## 7. 报告解读

### 7.1 报告格式

每次审计生成三种格式的报告：

| 格式 | 文件名 | 用途 |
|------|--------|------|
| HTML | `report.html` | 浏览器中查看，可视化交互 |
| Markdown | `report.md` | 文本阅读，易于分享 |
| JSON | `report.json` | 程序集成，自动化处理 |

### 7.2 报告内容结构

#### 统计摘要

显示四个指标卡片：

| 指标 | 含义 |
|------|------|
| 总问题数 | 所有发现问题的总数 |
| 高（高危+严重） | 需立即修复的高优先级问题 |
| 中（一般） | 建议修复但不紧急 |
| 低（建议） | 改进建议，可选修复 |

#### 健康度评分

五个维度的百分制评分：

- **安全性**：代码安全防护水平
- **性能**：运行效率与资源使用
- **健壮性**：异常处理与容错能力
- **可维护性**：代码结构与可读性
- **框架实践**：框架最佳实践遵循度

#### 问题卡片

每个问题包含：
- **严重级别**：高危 / 严重 / 一般 / 建议
- **分类**：安全 / 性能 / 健壮性 / 规范 / 框架 / 架构
- **置信度**：high / medium / low
- **文件位置**：精确到文件名和行号
- **问题描述**：对问题的详细说明
- **影响分析**：该问题可能带来的影响
- **证据**：问题存在的具体证据
- **修复建议**：文字说明 + 可选的代码片段

#### 历史对比

如果同仓库同分支存在上次审计记录，报告会自动对比并标注：
- **新增**：本次新出现的问题
- **已修复**：上次存在但本次已消失的问题
- **仍存在**：上次和本次都存在的问题

### 7.3 工件目录

报告目录下还包含 `artifacts/` 子目录，存放审计过程中的中间产物：

```
artifacts/
├── diff.json              # 完整 diff 信息
├── project_brief.json     # 项目画像结果
├── review_packs.json      # 发送给 LLM 的审计包
└── prescan_findings.json  # 规则预扫发现
```

---

## 8. 自定义审计规则

通过 `AIGUARD.md` 文件可以定制审计行为，让 AI 更了解你的项目。

### 配置方式

在 `config.yaml` 中指定规则文件路径：

```yaml
rules:
  custom_rule_file: "./AIGUARD.md"
```

### 规则文件示例

```markdown
## 项目用途
本项目是在线支付系统，处理用户交易和资金流转。关键链路包括：下单、支付、退款、对账。

## 高敏感目录
- src/auth — 认证模块
- src/payment — 支付模块
- internal/security — 安全模块

## 禁止用法
- 禁止 SQL 字符串拼接，必须使用参数化查询
- 禁止关闭 TLS 校验
- 禁止明文打印 token / password / secret
- 禁止在事务中执行外部 HTTP 调用
- 禁止硬编码密钥

## 框架约定
- Go HTTP 客户端必须设置 Timeout
- 所有导出接口必须做权限校验和审计日志

## 严重级别升级规则
- 若命中权限、认证、支付、导出链路，问题等级上调一级
- 若命中注入、越权、敏感信息泄露，至少判定为严重
```

---

## 9. 工作区管理

### 9.1 目录结构

工作区默认位于 `./workspace`，自动管理以下结构：

```
workspace/
├── repos/       # 仓库缓存（clone 的 bare repo）
├── worktrees/   # Git worktree 隔离环境
├── cache/       # 中间缓存
├── reports/     # 审计报告（按 taskId 组织）
│   └── {taskId}/
│       ├── report.html
│       ├── report.md
│       ├── report.json
│       └── artifacts/
└── logs/        # 运行日志
```

### 9.2 缓存清理

点击界面上的 **「清理缓存」** 按钮会清除：
- 仓库克隆缓存（下次需重新拉取）
- 所有历史报告
- 日志文件

> 注意：清理前会检查是否有正在运行的任务，运行中不允许清理。

---

## 10. 常见问题与排查

### Q: 仓库地址推导为空？

**原因**：MR/PR 链接格式不符合 GitHub/GitLab 标准模式。

**解决**：
- 当前版本即使未匹配 MR/PR 模式，也会尝试从 URL 中拼接地址
- 如果拼接结果不正确，可手动编辑仓库地址字段

### Q: LLM 接口联通性检查失败？

**排查步骤**：
1. 确认 `config.yaml` 中 `openai.base_url` 地址可达
2. 确认 `openai.api_key` 正确
3. 确认模型服务已启动
4. 如需代理，检查 `openai.proxy` 配置

### Q: 拉取代码失败？

**排查步骤**：
1. 检查 Git SSH 密钥是否配置（`ssh -T git@your-host`）
2. 如使用 HTTPS，确认凭据已缓存
3. 检查 `config.yaml` 中 `git` 配置的自定义主机/端口是否正确
4. 查看日志面板获取详细错误信息

### Q: 审计结果为空（没有发现问题）？

**可能原因**：
- 本次 diff 为空或仅包含二进制文件
- 变更文件被 `rules.ignore` 规则排除
- 模型确实认为没有问题（正常情况）

### Q: 审计速度慢？

**优化建议**：
- 增大 `runtime.concurrency`（默认 4，可调到 8）
- 减小 `review.max_changed_files` 限制审计范围
- 使用更快的模型或本地部署模型
- 关闭 `review.enable_project_brief` 可跳过画像阶段

### Q: 报告中的代码跳转链接不可用？

**解决**：在 `config.yaml` 中设置 `review.code_browse_base_url`，例如：

```yaml
review:
  code_browse_base_url: "https://gitlab.example.com/group/project/-/blob"
```

---

> 如有更多问题，请查看日志面板中的详细输出，或检查 `workspace/logs/` 目录下的日志文件。
