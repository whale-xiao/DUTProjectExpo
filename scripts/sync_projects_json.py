# -*- coding: utf-8 -*-
"""从当前数据库回导 docs/项目素材/项目汇总.json，使其与库内数据保持一致。

用途：项目在后台被改名/改内容后，让数据文件重新成为“事实快照”。
项目编号（number）取自后台编辑时使用的原始文档编号，按 RENAME 表映射改名项。
"""
import json, os, sys
import requests

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
OUT = os.path.join(ROOT, "docs", "项目素材", "项目汇总.json")
BASE = os.environ.get("BASE_URL", "http://127.0.0.1").rstrip("/")
USER = os.environ.get("ADMIN_USER", "admin")
PWD = os.environ.get("ADMIN_PASSWORD", "admin123456")

# 后台改过名的项目：原文档名 -> 当前库内名称（用于找回编号）
RENAME = {
    "基于AI智能体的某上市公司财务分析系统": "AIFR公司财务数据系统",
    "基于 AI 智能体的上市制造企业财务分析系统": "上市制造企业财务分析系统",
    "基于大模型的个人知识库系统": "Folio个人知识库",
    "WenKB基于大模型的个人知识库系统": "WenKB个人知识库系统",
    "基于 AI 智能体的康养系统（盘锦滨海家庭康养服务系统）": "盘锦滨海家庭康养服务系统",
    "智能财务分析系统（茅台厂近五年财报分析）": "茅台厂近五年财报分析",
}


def main():
    old = json.load(open(OUT, encoding="utf-8")) if os.path.exists(OUT) else []
    num_by_name = {}
    for p in old:
        num_by_name[p["name"]] = p.get("number")
        if p["name"] in RENAME:
            num_by_name[RENAME[p["name"]]] = p.get("number")

    s = requests.Session()
    r = s.post(BASE + "/api/auth/login", json={"username": USER, "password": PWD}, timeout=10).json()
    if r.get("code") != 0:
        print("登录失败:", r)
        return 1
    h = {"Authorization": "Bearer " + r["data"]["token"]}

    listing = s.get(BASE + "/api/admin/projects?page=1&page_size=200", headers=h, timeout=10) \
        .json().get("data", {}).get("list", [])

    records = []
    for item in listing:
        d = s.get(f"{BASE}/api/admin/projects/{item['id']}", headers=h, timeout=10).json().get("data") or {}
        blocks = d.get("blocks") or []
        intro = next((b.get("text") for b in blocks if b.get("type") == "paragraph" and b.get("text")), "")
        records.append({
            "number": num_by_name.get(d.get("name")),
            "id": d.get("id"),
            "name": d.get("name"),
            "summary": d.get("summary"),
            "intro": intro,
            "members": [m["name"] for m in (d.get("members") or [])],
            "tags": d.get("tags") or [],
        })

    records.sort(key=lambda x: (x["number"] is None, x["number"] or 0, x["id"]))
    json.dump(records, open(OUT, "w", encoding="utf-8"), ensure_ascii=False, indent=1)
    print(f"已回导 {len(records)} 个项目 -> {OUT}")
    for x in records:
        print(f"  {str(x['number']):>3} | id={x['id']:<3} | {x['name'][:34]:<34} | 成员 {len(x['members'])} 标签 {len(x['tags'])}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
