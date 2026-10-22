# -*- coding: utf-8 -*-
"""生成《黑盒测试报告》docx，数据来自 functional_results.json / perf_results.json。"""
import json, os, sys
from docx import Document
from docx.shared import Pt, Cm, RGBColor
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.enum.table import WD_TABLE_ALIGNMENT
from docx.oxml.ns import qn
from docx.oxml import OxmlElement

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.dirname(HERE)
RESULTS_DIR = os.path.join(HERE, "results")

F = json.load(open(os.path.join(RESULTS_DIR, "functional_results.json"), encoding="utf-8"))
P = json.load(open(os.path.join(RESULTS_DIR, "perf_results.json"), encoding="utf-8"))
OUT = sys.argv[1] if len(sys.argv) > 1 else os.path.join(ROOT, "docs", "黑盒测试报告.docx")

SONG, HEI, MONO, HFILL = "宋体", "黑体", "Consolas", "D9E2F3"
doc = Document()


def strip_theme_font(style):
    rf = style.element.rPr.find(qn("w:rFonts"))
    if rf is not None:
        for a in ("w:asciiTheme", "w:eastAsiaTheme", "w:hAnsiTheme", "w:cstheme"):
            rf.attrib.pop(qn(a), None)


n = doc.styles["Normal"]
n.font.name = SONG; n.font.size = Pt(10.5)
n.element.rPr.rFonts.set(qn("w:eastAsia"), SONG)
strip_theme_font(n)
n.paragraph_format.space_before = Pt(3); n.paragraph_format.space_after = Pt(3)
n.paragraph_format.line_spacing = 1.35

for lvl, sz in ((1, 16), (2, 14), (3, 12)):
    st = doc.styles[f"Heading {lvl}"]
    st.font.name = HEI; st.font.size = Pt(sz); st.font.bold = True
    st.font.color.rgb = RGBColor(0, 0, 0)
    st.element.rPr.rFonts.set(qn("w:eastAsia"), HEI)
    strip_theme_font(st)
    st.paragraph_format.space_before = Pt(14 if lvl == 1 else 10)
    st.paragraph_format.space_after = Pt(8)

for st in doc.styles:
    if st.style_id in ("Heading1Char", "Heading2Char", "Heading3Char"):
        st.font.color.rgb = RGBColor(0, 0, 0); st.font.name = HEI
        st.element.rPr.rFonts.set(qn("w:eastAsia"), HEI)
        strip_theme_font(st)
        c = st.element.rPr.find(qn("w:color"))
        if c is not None:
            c.attrib.pop(qn("w:themeColor"), None)

sec = doc.sections[0]
sec.page_width, sec.page_height = Cm(21), Cm(29.7)
sec.top_margin = sec.bottom_margin = Cm(2.5)
sec.left_margin = sec.right_margin = Cm(2.5)


def run(p, text, font=SONG, size=10.5, bold=False, color=None):
    r = p.add_run(text)
    r.font.name = font; r.font.size = Pt(size); r.font.bold = bold
    if color:
        r.font.color.rgb = color
    r._element.rPr.rFonts.set(qn("w:eastAsia"), font)
    return r


def para(text, align=None, indent=None, size=10.5, space=None, bold=False):
    p = doc.add_paragraph()
    if align is not None:
        p.alignment = align
    if indent:
        p.paragraph_format.first_line_indent = Pt(indent)
    if space is not None:
        p.paragraph_format.space_before = Pt(space)
    run(p, text, SONG, size, bold)
    return p


def head(level, text):
    p = doc.add_paragraph(style=f"Heading {level}")
    run(p, text, HEI, 16 if level == 1 else (14 if level == 2 else 12), True, RGBColor(0, 0, 0))
    return p


def shade(cell, color):
    tcPr = cell._tc.get_or_add_tcPr()
    e = OxmlElement("w:shd"); e.set(qn("w:val"), "clear"); e.set(qn("w:fill"), color)
    tcPr.append(e)


def table(header, rows, widths, size=8.5, mono_cols=()):
    t = doc.add_table(rows=len(rows) + 1, cols=len(widths))
    t.style = "Table Grid"; t.alignment = WD_TABLE_ALIGNMENT.CENTER; t.autofit = False
    for ci, h in enumerate(header):
        c = t.cell(0, ci); c.width = Pt(widths[ci] / 20.0); shade(c, HFILL)
        p = c.paragraphs[0]; p.alignment = WD_ALIGN_PARAGRAPH.CENTER
        run(p, h, SONG, size, True)
    for ri, row in enumerate(rows, start=1):
        for ci, v in enumerate(row):
            c = t.cell(ri, ci); c.width = Pt(widths[ci] / 20.0)
            p = c.paragraphs[0]
            p.paragraph_format.space_before = Pt(1); p.paragraph_format.space_after = Pt(1)
            run(p, str(v), MONO if ci in mono_cols else SONG, size)
    return t


def caption(t):
    para(t, align=WD_ALIGN_PARAGRAPH.CENTER, size=9)


def add_field(p, instr):
    for kind, text in (("begin", None), (None, instr), ("separate", None), ("end", None)):
        r = p.add_run(); r.font.name = SONG; r.font.size = Pt(9)
        r._element.rPr.rFonts.set(qn("w:eastAsia"), SONG)
        e = OxmlElement("w:fldChar" if kind else "w:instrText")
        if kind:
            e.set(qn("w:fldCharType"), kind)
        else:
            e.set(qn("xml:space"), "preserve"); e.text = text
        r._r.append(e)


sec.different_first_page_header_footer = True
hp = sec.header.paragraphs[0]; hp.alignment = WD_ALIGN_PARAGRAPH.CENTER
run(hp, "大连理工大学实训基地优秀项目展示系统 · 黑盒测试报告", SONG, 9, color=RGBColor(0x59, 0x59, 0x59))
pPr = hp._p.get_or_add_pPr(); bdr = OxmlElement("w:pBdr"); b = OxmlElement("w:bottom")
for k, v in (("w:val", "single"), ("w:sz", "4"), ("w:space", "1"), ("w:color", "999999")):
    b.set(qn(k), v)
bdr.append(b)
pPr.insert_element_before(bdr, "w:shd", "w:tabs", "w:suppressAutoHyphens", "w:kinsoku", "w:wordWrap",
                          "w:overflowPunct", "w:topLinePunct", "w:autoSpaceDE", "w:autoSpaceDN", "w:bidi",
                          "w:adjustRightInd", "w:snapToGrid", "w:spacing", "w:ind", "w:contextualSpacing",
                          "w:mirrorIndents", "w:suppressOverlap", "w:jc", "w:textDirection", "w:textAlignment",
                          "w:textboxTightWrap", "w:outlineLvl", "w:divId", "w:cnfStyle", "w:rPr", "w:sectPr")
fp = sec.footer.paragraphs[0]; fp.alignment = WD_ALIGN_PARAGRAPH.CENTER
run(fp, "第 ", SONG, 9); add_field(fp, "PAGE"); run(fp, " 页  共 ", SONG, 9)
add_field(fp, "NUMPAGES"); run(fp, " 页", SONG, 9)

# ═════════ 封面 ═════════
for line in ["项目名称：大连理工大学实训基地优秀项目展示系统", "文档类型：黑盒测试报告", "版本：1.0.0"]:
    para(line, align=WD_ALIGN_PARAGRAPH.CENTER, size=12, space=6)
for _ in range(4):
    doc.add_paragraph()
para("黑 盒 测 试 报 告", align=WD_ALIGN_PARAGRAPH.CENTER, size=26, bold=True)
for _ in range(4):
    doc.add_paragraph()
for line in ["测试执行日期：2026-10-13", "测试环境：Docker Compose 一键部署", "大连理工大学 · 创新实践基地"]:
    para(line, align=WD_ALIGN_PARAGRAPH.CENTER, size=12, space=6)
doc.add_page_break()

# ═════════ 目录 ═════════
head(1, "目　录")
for t_, lv in [("1  测试概述", 0), ("1.1  测试目的", 1), ("1.2  测试范围", 1), ("1.3  测试依据", 1),
               ("2  测试环境", 0), ("2.1  部署形态", 1), ("2.2  环境与版本", 1), ("2.3  测试数据", 1),
               ("3  测试方法与判定标准", 0), ("3.1  测试方法", 1), ("3.2  执行方式", 1), ("3.3  通过标准", 1),
               ("4  功能测试用例与执行结果", 0), ("4.1  结果汇总", 1), ("4.2  各模块用例明细", 1),
               ("5  性能测试", 0), ("5.1  测试场景与结果", 1), ("5.2  结果分析", 1),
               ("6  缺陷汇总与分析", 0), ("6.1  缺陷清单", 1), ("6.2  修复建议", 1),
               ("7  测试结论", 0)]:
    p = doc.add_paragraph()
    p.paragraph_format.left_indent = Pt(lv * 19)
    p.paragraph_format.space_before = Pt(1); p.paragraph_format.space_after = Pt(1)
    run(p, t_, HEI if lv == 0 else SONG, 10.5, lv == 0)
doc.add_page_break()

# ═════════ 1 概述 ═════════
head(1, "1  测试概述")
head(2, "1.1  测试目的")
para("验证「大连理工大学实训基地优秀项目展示系统」在实际部署形态下的功能正确性、接口健壮性与基本性能水平，确认系统满足需求分析与概要设计中约定的功能点，并为代码终版提交与答辩提供可追溯的测试证据。", indent=21)

head(2, "1.2  测试范围")
para("本轮采用黑盒方法，只依据接口契约与页面行为验证，不检视内部实现。范围如下：", indent=21)
table(["类别", "覆盖内容", "不在本轮范围"], [
    ["功能测试", "前台展示、搜索与联想、认证与鉴权、项目全生命周期、内容块编排、分类与标签、图片上传、异常与边界", "无"],
    ["性能测试", "首页、列表、详情、搜索、静态页面、静态图片、登录共 7 个场景的并发响应", "极限压力、稳定性长跑、容量规划"],
    ["未覆盖", "—", "浏览器兼容性、移动端适配、安全渗透（注入/越权/XSS 深度利用）、高可用与故障恢复"],
], [1500, 5000, 2572])
caption("表 1-1  测试范围")

head(2, "1.3  测试依据")
table(["序号", "依据", "说明"], [
    ["1", "需求分析说明书 V1.0", "功能点与业务规则来源"],
    ["2", "数据库设计 V1.0", "唯一约束、软删除、外键引用等约束来源"],
    ["3", "UML 设计 V1.0 §7 接口清单", "接口路径、方法与语义"],
    ["4", "统一响应约定", "code=0 成功；4xxxx 客户端错误；5xxxx 服务端错误"],
], [700, 3200, 5172])

# ═════════ 2 环境 ═════════
head(1, "2  测试环境")
head(2, "2.1  部署形态")
para("被测系统以 Docker Compose 一键部署，三个服务在同一内部网络：nginx（对外 80 / 8081）、Go 后端（8080，不对外发布）、MySQL 8.0（3306，不对外发布）。测试通过 nginx 访问，与真实用户路径一致。", indent=21)
table(["服务", "镜像 / 版本", "职责", "对外端口"], [
    ["nginx", "nginx:alpine（自建 dut-expo/nginx）", "托管两个前端静态资源，反代 /api 与 /uploads", "80、8081"],
    ["backend", "多阶段构建（golang:1.26-alpine → alpine:3.20）", "Go + Gin + GORM 业务服务", "不发布"],
    ["mysql", "mysql:8.0", "业务库 project_showcase，9 张表", "不发布"],
], [1300, 3300, 3172, 1300])
caption("表 2-1  服务清单")

head(2, "2.2  环境与版本")
table(["项目", "版本 / 说明"], [
    ["操作系统", "Windows 11（Docker Desktop 29.6.1，Compose v5.2.0）"],
    ["后端运行时", "Go 1.26.x，静态编译，CGO_ENABLED=0"],
    ["数据库", "MySQL 8.0，utf8mb4 / utf8mb4_general_ci"],
    ["前端", "Vue 3 + Vite 8，构建后由 nginx 托管"],
    ["测试客户端", "Python 3.11 + requests，脚本化直接调用 HTTP 接口并断言"],
    ["测试时间", F["time"]],
], [2400, 6672])

head(2, "2.3  测试数据")
para("使用后端启动时自动写入的种子数据：15 个已上架项目、4 个分类、30 个标签，其中 5 个项目处于推荐位。", indent=21)
para("涉及写操作的用例（新建项目、分类、标签）统一使用「【测试】」前缀 + 随机后缀命名，用例结束后删除；功能性验证完成后对残留的软删除记录做了物理清理，数据库已还原为 15 / 4 / 30 的种子状态。", indent=21)

# ═════════ 3 方法 ═════════
head(1, "3  测试方法与判定标准")
head(2, "3.1  测试方法")
for t_ in ["等价类划分：对每个输入域划分有效类与无效类，如项目名的空 / 正常 / 超长、图片的合法扩展名 / 非法扩展名 / 超限体积。",
           "边界值分析：分页 page_size 超出上限、文件体积接近 5MB、联想条数上限 6 条等。",
           "错误推测：结合设计文档中已知的薄弱点设计用例，如软删除与唯一约束的冲突、内容块字段缺失、非法枚举值。",
           "场景法：串起「新建草稿 → 前台不可见 → 整单保存 → 上架 → 前台可见 → 下架 → 删除」的完整业务链路。"]:
    para(t_, indent=21)

head(2, "3.2  执行方式")
para("所有功能用例由脚本化客户端执行：直接向 nginx 暴露的地址发起 HTTP 请求，对 HTTP 状态码、响应体 code、关键业务字段逐项断言，并记录实际返回值。用例可在任意一次部署后重复执行，结果写入 JSON 供本报告汇总。", indent=21)
para("性能用例采用固定并发、固定请求数的压测方式，每个场景先预热 20 次再计时，统计吞吐量、平均耗时与 P50 / P95 / P99 分位。", indent=21)

head(2, "3.3  通过标准")
table(["维度", "通过标准"], [
    ["功能", "断言全部成立：响应 code 与预期一致，关键业务字段正确，错误场景返回约定的错误码"],
    ["性能", "常规接口 P95 < 200ms，登录等含计算接口 P95 < 300ms，全部场景错误数为 0"],
    ["数据", "用例执行后不得遗留测试数据，演示数据保持可用"],
], [1600, 7472])

# ═════════ 4 功能结果 ═════════
head(1, "4  功能测试用例与执行结果")
head(2, "4.1  结果汇总")
total, passed, failed = F["total"], F["passed"], F["failed"]
table(["指标", "数值", "说明"], [
    ["用例总数", str(total), "覆盖 8 个功能模块"],
    ["通过", str(passed), "预期与实际一致"],
    ["未通过", str(failed), "均为真实缺陷，见第 6 章"],
    ["通过率", "%.1f%%" % (passed * 100.0 / total), "按用例条数计"],
], [2400, 1800, 4872])
caption("表 4-1  功能测试结果汇总")

head(2, "4.2  各模块用例明细")
groups = []
for r in F["results"]:
    if not groups or groups[-1][0] != r["module"]:
        groups.append((r["module"], []))
    groups[-1][1].append(r)

for gi, (mod, items) in enumerate(groups):
    okn = sum(1 for x in items if x["status"] == "通过")
    head(3, "%s（%d/%d 通过）" % (mod, okn, len(items)))
    rows = [[r["id"], r["title"], r["steps"], r["expected"], r["actual"], r["status"]] for r in items]
    table(["编号", "测试项", "测试步骤", "预期结果", "实际结果", "结论"], rows,
          [800, 1700, 2200, 1650, 2122, 600], size=8)
    caption("表 4-%d  %s用例执行结果" % (gi + 2, mod))

# ═════════ 5 性能 ═════════
doc.add_page_break()
head(1, "5  性能测试")
head(2, "5.1  测试场景与结果")
rows = []
for r in P["results"]:
    rows.append([r["name"], str(r["requests"]), str(r["concurrency"]),
                 "%.1f" % r["rps"], "%.1f" % r["avg_ms"], "%.1f" % r["p95_ms"],
                 "%.1f" % r["p99_ms"], str(r["errors"])])
table(["场景", "请求数", "并发", "RPS", "平均(ms)", "P95(ms)", "P99(ms)", "错误数"], rows,
      [2200, 900, 800, 1100, 1100, 1100, 1100, 772], size=8.5)
caption("表 5-1  性能测试结果（本机单实例，非生产规格）")

head(2, "5.2  结果分析")
para("全部 7 个场景错误数均为 0，说明在 200~500 RPS 的压力下系统未出现超时或异常，接口稳定性满足演示与课程验收需要。", indent=21)
para("读接口表现最好：搜索联想与首页聚合的 P95 均在 35ms 以内，得益于按主键与联合索引的简单查询以及响应体精简（0.5KB / 2.3KB）。", indent=21)
para("项目详情的平均耗时约 90ms，明显高于列表接口。原因在于详情接口除多表查询外，还会在同一事务内执行 view_count + 1 的写操作；这是设计上有意的取舍（浏览量字段冗余），并发写锁带来额外开销。若后续浏览量成为瓶颈，可按设计文档的二期方案改为 Redis 计数后定时落库。", indent=21)
para("静态图片的平均大小为 299KB，P95 约 131ms，主要耗时在图片传输而非处理；实际部署时建议开启 nginx 的 gzip（已对文本类开启）并为图片增加强缓存或接入 CDN。", indent=21)
para("登录接口 85.8ms 的平均耗时符合预期，主要成本来自 bcrypt 的故意慢哈希，属于安全设计而非性能缺陷。", indent=21)
para("需要说明：以上为本机 Docker Desktop 单机环境数据，磁盘、CPU 与容器网络均共享宿主机，绝对值不代表生产性能，仅用于横向比较各接口的相对开销。", indent=21)

# ═════════ 6 缺陷 ═════════
head(1, "6  缺陷汇总与分析")
head(2, "6.1  缺陷清单")
defects = [
    ["D-01", "FT-34", "内容块中类型为 image 的块，image_url 传空字符串时仍被接受并成功保存，会生成前台无法渲染的空图块。", "中"],
    ["D-02", "FT-35", "类型为 heading / paragraph 的块，text 为空时仍被接受，会生成空标题或空段落。", "中"],
    ["D-03", "FT-36", "非法 block_type（如 video）未被参数校验拦截，写入数据库时触发 MySQL 1265 Data truncated，最终返回 50000「保存失败」，而非 40000 参数错误。", "中"],
    ["D-04", "FT-51", "分类 / 标签软删除后，以同名再次创建会命中唯一索引（索引不含 deleted_at），返回 50000，用户无法重建同名资源；该问题在数据库设计 §8 中已被标注为待确认项，本轮测试确认其已实际发生。", "高"],
]
table(["编号", "关联用例", "缺陷描述", "严重程度"], defects, [800, 1100, 6172, 1000], size=8.5)
caption("表 6-1  缺陷清单")

head(2, "6.2  修复建议")
table(["编号", "根因", "建议修复方式", "涉及位置"], [
    ["D-01 / D-02", "服务层只做字段映射，未对内容块做字段级必填校验", "在块转换时按类型校验：heading/paragraph 要求 text 去空格后非空；image 要求 url 非空，不满足返回 40000", "service/admin_service.go 内容块组装处"],
    ["D-03", "dto 中切片元素的 binding 未生效（缺少 dive），非法枚举值直达数据库", "为 Blocks 字段补 `binding:\"dive\"`，或在服务层用 switch 显式校验 type 并返回 40000", "dto/admin.go、service/admin_service.go"],
    ["D-03 / D-04", "数据库约束错误未映射，统一落到 50000", "在 repository 层识别 MySQL 1062 / 1265，分别映射为 ErrDuplicate（40900）与参数错误（40000）", "repository 各 *_repo.go、service/errors.go"],
    ["D-04", "唯一索引仅含 name，软删除行仍占用索引", "二选一：唯一索引改为 (name, deleted_at)；或删除时物理清理 / 重建前检测软删除记录并明确提示", "repository/db.go 迁移、meta_service.go"],
], [1400, 2500, 3500, 1672], size=8.5)
caption("表 6-2  修复建议")

para("说明：按项目协作约定，本轮只输出测试结论与修复建议，未改动任何业务代码。上述 4 条缺陷均为接口层可复现的问题，修复后可用同一套脚本回归验证。", indent=21)

# ═════════ 7 结论 ═════════
head(1, "7  测试结论")
para("本轮黑盒测试共执行功能用例 %d 条，通过 %d 条，未通过 %d 条，通过率 %.1f%%；性能测试覆盖 7 个场景，全部零错误。" %
     (total, passed, failed, passed * 100.0 / total), indent=21)
para("核心业务链路（登录鉴权、项目整单保存、内容块全量替换、草稿不上前台、上架可见、下架与删除、分类引用保护、图片上传限制）均按设计预期工作，未发现阻断性问题。", indent=21)
para("未通过的 4 条集中在内容块字段校验与错误码映射两类问题：前者依赖前端约束、服务端缺少兜底，后者把可预期的客户端错误暴露为 50000。二者均不影响主流程演示，但会降低接口健壮性与可维护性，建议在代码终版提交前修复，并用同一套脚本回归。", indent=21)
para("性能方面，系统在本机单实例下可稳定承载 200~500 RPS，读接口 P95 在 35ms 以内，满足课程项目的验收与答辩演示需求。", indent=21)
para("综合判定：系统功能基本完备、行为符合设计，达到可提交状态；建议按第 6 章修复 4 条缺陷后封版。", indent=21)
para("测试执行：脚本化自动执行　｜　测试结论：有条件通过", align=WD_ALIGN_PARAGRAPH.CENTER, space=18, bold=True)

doc.save(OUT)
print("saved", OUT)
