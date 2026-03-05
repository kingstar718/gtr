# Todo 模块实现Plan

## 📋 **整体架构**

```
gtr/
├── cmd/
│   └── todo/
│       ├── add.go          # 添加任务命令
│       ├── list.go         # 列表任务命令
│       ├── edit.go         # 编辑任务命令
│       ├── delete.go       # 删除任务命令
│       └── interactive.go  # 交互菜单模式
├── internal/
│   └── todo/
│       ├── model.go        # Todo数据结构
│       ├── storage.go      # 数据存储（JSON/SQLite）
│       ├── service.go      # 业务逻辑
│       └── validator.go    # 数据验证
└── main.go
```

---

## 🎯 **核心功能清单**

### **1. 数据模型** 

```go
Task结构体：
- ID (自动生成)
- Title (标题)
- Description (描述)
- Priority (优先级: High/Medium/Low)
- Status (状态: Pending/InProgress/Done)
- DueDate (截止日期)
- Tags (标签列表)
```

### **2. 命令行模式**（直接输入参数）

```bash
gtr todo add "完成报告" --priority high --due 2025-03-10
gtr todo list [--status done] [--priority high] [--tag work]
gtr todo edit <id> --status done
gtr todo delete <id>
gtr todo view <id>
```

### **3. 交互式模式**（无参数启动）

```
gtr todo  # 启动菜单
→ 主菜单：
  1. 查看所有任务
  2. 添加新任务
  3. 编辑任务
  4. 删除任务
  5. 搜索任务
  6. 退出

→ 添加任务流程：
  - 输入标题
  - 选择优先级 (survey Select)
  - 输入描述 (survey Input)
  - 选择标签 (survey MultiSelect)
  - 设置截止日期 (survey Confirm)
```

### **4. 数据存储**

```
存储位置：~/.gtr/todos.json
格式：
{
  "tasks": [
    {
      "id": "uuid-xxxx",
      "title": "完成报告",
      "description": "完成季度报告",
      "priority": "high",
      "status": "pending",
      "dueDate": "2025-03-10",
      "tags": ["work", "urgent"]
    }
  ]
}
```

### **5. 列表展示**

```
ID    | Title      | Priority | Status    | DueDate    | Tags
------|------------|----------|-----------|------------|--------
1     | 完成报告   | High     | Pending   | 2025-03-10 | work
2     | 学习Go     | Medium   | InProgress| 2025-03-15| learn
```

### **6. 搜索和过滤**

```bash
gtr todo list --status done              # 按状态
gtr todo list --priority high            # 按优先级
gtr todo list --tag work                 # 按标签
gtr todo list --due-before 2025-03-10    # 按截止日期
gtr todo search "关键词"                  # 全文搜索
```

---

## 📦 **技术方案选择**

| 模块 | 方案 | 原因 |
|------|------|------|
| CLI框架 | Cobra | 已有，成熟 |
| 交互界面 | survey | 功能全，API简洁 |
| 存储 | JSON | 轻量，无依赖 |
| ID生成 | UUID | 跨平台一致性 |
| 时间处理 | time包 | 标准库足够 |
| 彩色输出 | fatih/color | 轻量，易用 |

---

## 🔨 **实现步骤**（可选顺序）

### **第一阶段（基础功能）：**
- [x] 创建Todo数据结构和JSON存储
- [x] 实现add命令（命令行+交互）
- [x] 实现list命令，支持基础显示
- [x] 实现delete命令

### **第二阶段（增强功能）：**
- [x] 使用bubble tea重构TUI界面
- [x] 实现edit命令（支持优先级和状态快速切换）
- [x] 添加过滤和搜索功能（按状态、优先级过滤）
- [x] 彩色输出和表格展示（已完成，使用lipgloss）
- [x] 状态管理（Pending→InProgress→Done）

### **第三阶段（优化）：**
- [ ] 数据备份/导出功能
- [ ] 快捷键和别名
- [ ] 时间提醒（可选）
- [ ] 同步功能（可选）

---

## ⚙️ **技术细节**

### **优先级考虑：**

```
High   → 红色显示
Medium → 黄色显示
Low    → 绿色显示
```

### **状态流转：**

```
Pending → InProgress → Done
  ↑                      ↓
  └──────────────────────┘ (可撤销)
```

### **数据校验：**

- 标题不能为空
- 截止日期不能早于当前日期
- 优先级必须是预定义值

---

## 📊 **预期工作量**

| 阶段 | 代码行数 | 时间 |
|------|---------|------|
| 基础功能 | ~800行 | 中等 |
| 增强功能 | ~500行 | 中等 |
| 优化 | ~300行 | 小 |
| 测试 | ~400行 | 中等 |

---

## 🎁 **额外功能建议**（可选）

- 导出为CSV/Markdown
- 定期任务（重复）
- 子任务支持
- 统计仪表板（完成率、优先级分布等）
- 快捷命令（`gtr t` 快速添加）

---

## 📝 **修改记录**

| 日期 | 修改内容 |
|------|---------|
| 2025-03-05 | 初始版本 |
| 2025-03-05 | 完成第一阶段基础功能实现 |
| 2025-03-05 | 升级Go到1.24，集成bubble tea TUI框架 |
| 2025-03-05 | 完成第二阶段增强功能（TUI、编辑、过滤、状态管理） |
| 2025-03-05 | 用bubbles v1组件重构TUI，完全弃用lipgloss直接使用 |
| 2025-03-05 | 用bubbles Table组件替换List组件，移除ID列，添加表格标题 |
