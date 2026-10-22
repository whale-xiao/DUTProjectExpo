# 黑盒测试

脚本化黑盒测试，只通过 HTTP 接口与页面行为验证，不接触内部实现。对应第 8 周「黑盒测试」交付物。

## 前置条件

被测系统需要先跑起来（二选一）：

```bash
# 方式 A：容器化（推荐，与测试报告里的环境一致）
docker compose up -d --build

# 方式 B：本地三终端
cd backend && go run .                 # :8080
cd frontend-showcase && npm run dev    # :5173
cd frontend-admin && npm run dev       # :5174
```

方式 A 下脚本默认访问 `http://127.0.0.1`（nginx :80），方式 B 下需把 `BASE_URL` 指到后端：

```bash
BASE_URL=http://127.0.0.1:8080 uv run --with requests python tests/functional_test.py
```

## 目录

```
tests/
├── functional_test.py      功能测试（51 条用例，8 个模块）
├── perf_test.py            性能测试（7 个场景，固定并发）
├── gen_test_doc.py         由结果 JSON 生成《黑盒测试报告》docx
└── results/
    ├── functional_results.json
    └── perf_results.json
```

## 运行

```bash
# 功能测试：逐条断言，控制台打印 PASS/FAIL，结果写入 results/
uv run --with requests python tests/functional_test.py

# 性能测试：每个场景先预热 20 次再计时，统计 RPS / 平均 / P50 / P95 / P99
uv run --with requests python tests/perf_test.py

# 由结果生成测试报告（默认输出到 docs/黑盒测试报告.docx）
uv run --with python-docx python tests/gen_test_doc.py
```

可用环境变量覆盖：

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `BASE_URL` | `http://127.0.0.1` | 被测地址（容器化用 nginx 的 80） |
| `ADMIN_USER` | `admin` | 后台账号 |
| `ADMIN_PASSWORD` | `admin123456` | 后台密码 |

## 测试数据约定

写操作用例统一使用「【测试】」前缀 + 随机后缀命名，用例结束自动删除，不会污染演示数据。
但**删除接口是软删除**，同名资源无法重建（见下方 D-04），因此脚本每次运行都使用新的随机后缀；如果反复运行导致库里积累软删行，可物理清理：

```sql
DELETE FROM projects WHERE deleted_at IS NOT NULL;
DELETE FROM categories WHERE deleted_at IS NOT NULL;
DELETE FROM tags WHERE deleted_at IS NOT NULL;
```

## 当前基线（2026-10-13）

功能 51 条：**通过 47，未通过 4**。这 4 条是真实缺陷，不是脚本问题：

| 编号 | 用例 | 现象 |
| --- | --- | --- |
| D-01 | FT-34 | `image` 内容块 `url` 为空仍可保存 |
| D-02 | FT-35 | `heading` / `paragraph` 块 `text` 为空仍可保存 |
| D-03 | FT-36 | 非法 `block_type` 返回 50000（应为 40000） |
| D-04 | FT-51 | 软删除后同名分类/标签无法重建，返回 50000 |

修复后重跑本套脚本即可回归，预期变成 51/51。缺陷根因与修复建议见 `docs/黑盒测试报告.docx` 第 6 章。
