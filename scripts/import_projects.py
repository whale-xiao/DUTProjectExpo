# -*- coding: utf-8 -*-
"""把 docs/项目素材/项目汇总.json 导入系统（通过后台接口，自动建标签）。

用法：
    docker compose up -d
    uv run --with requests python scripts/import_projects.py

可用环境变量覆盖 BASE_URL / ADMIN_USER / ADMIN_PASSWORD。
封面统一使用 /uploads/cover/placeholder.png（需先在后台或容器内放入该文件）。

注意：脚本按项目名判断是否已存在。若库中的项目被改过名，脚本会中止以避免重复创建；
确认无误可加 FORCE=1 强制继续。替换演示数据前请先清空 projects 表。
"""
import json, os, sys
import requests

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SRC = os.path.join(ROOT, "docs", "项目素材", "项目汇总.json")
BASE = os.environ.get("BASE_URL", "http://127.0.0.1").rstrip("/")
USER = os.environ.get("ADMIN_USER", "admin")
PWD = os.environ.get("ADMIN_PASSWORD", "admin123456")
COVER = "/uploads/cover/placeholder.png"

# 项目编号 -> 分类名称（分类名必须已存在，否则脚本报错退出，不做静默兜底）
CATEGORY = {
    1: "基于AI智能体的 某上市公司 财务分析系统",
    2: "基于AI智能体的 某上市公司 财务分析系统",
    3: "基于AI智能体的 某上市公司 财务分析系统",
    5: "基于AI 智能体的 房屋租赁系统",
    8: "基于AI 智能体的 房屋租赁系统",
    13: "基于大模型的 个人知识库系统",
    14: "基于大模型的 个人知识库系统",
    17: "基于AI 智能体的 康养系统",
    18: "基于AI 智能体的 康养系统",
    20: "基于AI 智能体的 康养系统",
    21: "基于AI 智能体的 量化分析",
    23: "基于AI 智能体的 量化分析",
    25: "基于yolo模型的 实体识别系统",
    26: "基于yolo模型的 实体识别系统",
    28: "基于yolo模型的 实体识别系统",
    30: "基于华为Iot平台 某物联网系统",
}
DEFAULT_CATEGORY = "基于AI智能体的 某上市公司 财务分析系统"
# 首页推荐位（项目编号 -> 排序）
RECOMMENDED = {2: 1, 17: 2, 25: 3}


def main():
    projects = json.load(open(SRC, encoding="utf-8"))

    # 只导入指定编号：python scripts/import_projects.py 21 30
    only = {int(a) for a in sys.argv[1:] if a.isdigit()}
    if only:
        projects = [p for p in projects if p["number"] in only]
        print("本次仅导入编号:", sorted(only))

    s = requests.Session()

    r = s.post(BASE + "/api/auth/login", json={"username": USER, "password": PWD}, timeout=10).json()
    if r.get("code") != 0:
        print("登录失败:", r)
        return 1
    h = {"Authorization": "Bearer " + r["data"]["token"]}

    cats = s.get(BASE + "/api/categories", timeout=10).json().get("data") or []
    cat_id = {c["name"]: c["id"] for c in cats}
    print("现有分类:", cat_id)

    # 防重复：库里存在数据文件中没有的项目（通常是被改过名），直接中止
    existing = s.get(BASE + "/api/admin/projects?page=1&page_size=200", headers=h, timeout=10) \
        .json().get("data") or {}
    existing_names = {x["name"] for x in existing.get("list", [])}
    json_names = {p["name"] for p in projects}
    unknown = existing_names - json_names
    if unknown and not only and os.environ.get("FORCE") != "1":
        print("\n⚠️ 库中存在与数据文件不匹配的项目（可能已被改名）：")
        for n in sorted(unknown):
            print("   -", n)
        print("\n继续会按名称判定“不存在”而重复创建。确认无误请加 FORCE=1 重跑。")
        return 2

    ok = skip = fail = 0
    for p in projects:
        cname = CATEGORY.get(p["number"], DEFAULT_CATEGORY)
        cid = cat_id.get(cname)
        if cid is None:
            print(f"  ! [{p['number']}] {p['name']}  分类「{cname}」不存在，已跳过")
            fail += 1
            continue

        blocks = [
            {"type": "heading", "text": "项目介绍"},
            {"type": "paragraph", "text": p["intro"] or p["name"]},
        ]
        if p["members"]:
            blocks += [{"type": "heading", "text": "小组成员"},
                       {"type": "paragraph", "text": "、".join(p["members"])}]
        if p["tags"]:
            blocks += [{"type": "heading", "text": "主要技术"},
                       {"type": "paragraph", "text": "、".join(p["tags"])}]

        body = {
            "name": p["name"],
            "summary": p["summary"],
            "categoryId": cid,
            "coverUrl": COVER,
            "status": "published",
            "isRecommended": p["number"] in RECOMMENDED,
            "recommendOrder": RECOMMENDED.get(p["number"], 0),
            "tags": p["tags"],
            "members": [{"name": m, "role": "项目成员"} for m in p["members"]],
            "blocks": blocks,
        }
        r = s.post(BASE + "/api/admin/projects", headers=h, json=body, timeout=20).json()
        code = r.get("code")
        if code == 0:
            ok += 1
            print(f"  + [{p['number']}] {p['name']}  (分类 {cname})")
        elif code == 40900:
            skip += 1
            print(f"  = [{p['number']}] {p['name']}  已存在，跳过")
        else:
            fail += 1
            print(f"  ! [{p['number']}] {p['name']}  失败: {r}")

    print(f"\n导入完成：新增 {ok}，跳过 {skip}，失败 {fail}")
    return 0 if fail == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
