# 交付说明（2026-03-28）

本次交付针对 UI 可用性、构建稳定性与报告交互做了生产向增强，重点如下：

## 已完成项

1. **Windows 构建脚本增强**
   - `scripts/build_windows.ps1` 现在会按 **Windows Kit 8.1 -> 10 -> 11** 的顺序探测并尝试构建。
   - 同时对 `signtool.exe` 也做了同顺序回退查找。
   - 支持通过环境变量 `WINDOWS_KITS_81_PATH` / `WINDOWS_KITS_10_PATH` / `WINDOWS_KITS_11_PATH` 显式指定 SDK 根目录。

2. **仓库地址自动推导**
   - 粘贴 MR/PR 链接后，前端会自动根据当前配置推导仓库地址。
   - 仓库地址仍然支持手工改写；若手工填写，则优先采用手工值。

3. **历史记录修复**
   - 修复了历史列表字段映射错误。
   - 现在可稳定显示：总问题数、高/中/低分布、报告打开入口。

4. **主界面报告交互增强**
   - 默认折叠显示问题摘要（3 行）。
   - 点击展开后可查看：详细描述、影响分析、证据、修复建议。
   - 若模型给出 `recommendationCode`，会在界面与 HTML/Markdown 报告中展示建议代码片段。

5. **查看报告能力**
   - 新增“查看报告”和“打开目录”操作。
   - 优先打开 HTML 报告，不存在时自动回退打开报告目录。

6. **UI 美化**
   - 主界面升级为蓝紫 / 蓝绿系玻璃质感风格。
   - 提升信息层级、按钮反馈、状态展示与结果可读性。

## 主要变更文件

- `scripts/build_windows.ps1`
- `app.go`
- `internal/uiapi/dto.go`
- `internal/model/types.go`
- `internal/findings/normalize.go`
- `internal/review/prompts.go`
- `internal/review/orchestrator.go`
- `internal/report/render.go`
- `internal/report/templates/report.html.tmpl`
- `frontend/src/App.vue`
- `frontend/src/components/ReviewForm.vue`
- `frontend/src/components/ProgressPanel.vue`
- `frontend/src/components/HistoryList.vue`
- `frontend/src/components/ReportViewer.vue`
- `frontend/src/bridge.ts`
- `frontend/src/types.ts`
- `frontend/src/styles.css`

## 说明

由于当前交付环境无法联网安装 Go / Node 依赖，因此本次主要完成的是**源码级修复与可交付打包**。
建议在你本地或 CI 机器执行以下验证：

```bash
# 前端
cd frontend
npm install
npm run build

# 后端
cd ..
go test ./...

# Windows 构建（PowerShell）
./scripts/build_windows.ps1
```
