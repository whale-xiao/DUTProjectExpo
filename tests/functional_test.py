# -*- coding: utf-8 -*-
"""黑盒功能测试：真实调用正在运行的 Docker 全栈，逐条断言并输出 JSON 结果。"""
import json, os, sys, io, time, random, string
import requests

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.dirname(HERE)
RESULTS_DIR = os.path.join(HERE, "results")

# 目标地址与账号可用环境变量覆盖（默认指向 docker compose 部署的 nginx）
BASE = os.environ.get("BASE_URL", "http://127.0.0.1").rstrip("/")
ADMIN = os.environ.get("ADMIN_USER", "admin")
PWD = os.environ.get("ADMIN_PASSWORD", "admin123456")
OUT = os.path.join(RESULTS_DIR, "functional_results.json")
MARK = "【测试】"
RESULTS = []


def rec(cid, module, title, steps, expected, actual, ok):
    RESULTS.append({"id": cid, "module": module, "title": title, "steps": steps,
                    "expected": expected, "actual": actual, "status": "通过" if ok else "未通过"})
    print(("  [PASS] " if ok else "  [FAIL] ") + cid + " " + title + "  -> " + str(actual))


def code(r):
    try:
        return r.json().get("code")
    except Exception:
        return None


def body(r):
    try:
        return r.json()
    except Exception:
        return {}


def main():
    s = requests.Session()
    token = None

    def hdr():
        return {"Authorization": "Bearer " + token} if token else {}

    def rnd(n=6):
        return "".join(random.choices(string.ascii_lowercase + string.digits, k=n))

    created_projects = []

    # ───────────── A. 前台展示 ─────────────
    M = "A 前台展示"
    r = s.get(BASE + "/api/home", timeout=10)
    d = body(r).get("data") or {}
    ok = r.status_code == 200 and code(r) == 0 and "recommended" in d and len(d.get("recommended", [])) > 0
    rec("FT-01", M, "首页聚合接口返回推荐项目与分类", "GET /api/home", "code=0 且 recommended 非空",
        "code=%s recommended=%d" % (code(r), len(d.get("recommended", []))), ok)

    r = s.get(BASE + "/api/projects?page=1&page_size=8", timeout=10)
    d = body(r).get("data") or {}
    ok = code(r) == 0 and len(d.get("list", [])) == 8 and d.get("total", 0) >= 8
    rec("FT-02", M, "项目列表默认分页返回 8 条", "GET /api/projects?page=1&page_size=8",
        "本页 8 条，total≥8", "本页=%d total=%s" % (len(d.get("list", [])), d.get("total")), ok)

    cats = (body(s.get(BASE + "/api/categories", timeout=10)).get("data") or [])
    cid = cats[0]["id"] if cats else 1
    r = s.get(BASE + "/api/projects?category_id=%d&page_size=100" % cid, timeout=10)
    d = body(r).get("data") or {}
    ok = code(r) == 0 and all(p["categoryId"] == cid for p in d.get("list", [])) and len(d.get("list", [])) > 0
    rec("FT-03", M, "按分类筛选项目", "GET /api/projects?category_id=%d" % cid,
        "返回项目均属于该分类", "命中=%d 条，分类一致=%s" % (len(d.get("list", [])),
        "是" if all(p["categoryId"] == cid for p in d.get("list", [])) else "否"), ok)

    tags = (body(s.get(BASE + "/api/tags", timeout=10)).get("data") or [])
    tid = tags[0]["id"] if tags else 1
    r = s.get(BASE + "/api/projects?tag_id=%d&page_size=100" % tid, timeout=10)
    d = body(r).get("data") or {}
    ok = code(r) == 0 and len(d.get("list", [])) > 0
    rec("FT-04", M, "按标签筛选项目", "GET /api/projects?tag_id=%d" % tid,
        "code=0 且命中项目", "命中=%d 条" % len(d.get("list", [])), ok)

    views = []
    for sort in ("views", "newest"):
        d = body(s.get(BASE + "/api/projects?sort=%s&page_size=10" % sort, timeout=10)).get("data") or {}
        lst = d.get("list", [])
        views.append([p["viewCount"] for p in lst] if sort == "views" else [p["createdAt"] for p in lst])
    ok = views[0] == sorted(views[0], reverse=True) and views[1] == sorted(views[1], reverse=True)
    rec("FT-05", M, "排序参数 views / newest 生效", "GET /api/projects?sort=views|newest",
        "浏览量降序、创建时间降序", "两者均为降序=%s" % ok, ok)

    d = body(s.get(BASE + "/api/projects?page_size=500", timeout=10)).get("data") or {}
    ok = d.get("page_size", 0) <= 100
    rec("FT-06", M, "page_size 超上限被截断", "GET /api/projects?page_size=500",
        "page_size 被限制为 ≤100", "page_size=%s" % d.get("page_size"), ok)

    r = s.get(BASE + "/api/projects/1", timeout=10)
    d = body(r).get("data") or {}
    blocks = d.get("blocks") or []
    ok = code(r) == 0 and len(blocks) > 0 and d.get("name")
    rec("FT-07", M, "项目详情返回基础信息与内容块", "GET /api/projects/1",
        "code=0 且 blocks 非空", "项目=%s blocks=%d" % (d.get("name"), len(blocks)), ok)

    r = s.get(BASE + "/api/projects/999999", timeout=10)
    ok = code(r) == 40400
    rec("FT-08", M, "查询不存在的项目返回 40400", "GET /api/projects/999999",
        "code=40400", "code=%s http=%s" % (code(r), r.status_code), ok)

    d = body(s.get(BASE + "/api/categories", timeout=10)).get("data") or []
    ok = isinstance(d, list) and len(d) >= 1 and "name" in d[0]
    rec("FT-09", M, "分类列表接口", "GET /api/categories", "返回分类数组", "分类数=%d" % len(d), ok)

    d = body(s.get(BASE + "/api/tags", timeout=10)).get("data") or []
    ok = isinstance(d, list) and len(d) >= 1
    rec("FT-10", M, "标签列表接口", "GET /api/tags", "返回标签数组", "标签数=%d" % len(d), ok)

    # ───────────── B. 搜索 ─────────────
    M = "B 搜索"
    d = body(s.get(BASE + "/api/projects?keyword=Go&page_size=50", timeout=10)).get("data") or {}
    hits = d.get("list", [])
    ok = len(hits) > 0
    rec("FT-11", M, "关键词搜索命中项目", "GET /api/projects?keyword=Go",
        "返回命中项目", "命中=%d 条" % len(hits), ok)

    d = body(s.get(BASE + "/api/search/suggest?keyword=Go", timeout=10)).get("data") or []
    ok = isinstance(d, list) and 0 < len(d) <= 6
    rec("FT-12", M, "搜索联想返回不超过 6 条", "GET /api/search/suggest?keyword=Go",
        "0 < 条数 ≤ 6", "条数=%d" % len(d), ok)

    r = s.get(BASE + "/api/search/suggest?keyword=", timeout=10)
    d = body(r).get("data")
    ok = r.status_code == 200 and code(r) == 0 and (d == [] or d is None)
    rec("FT-13", M, "空关键词联想不报错", "GET /api/search/suggest?keyword=",
        "code=0，返回空数组", "code=%s data=%s" % (code(r), d), ok)

    d = body(s.get(BASE + "/api/projects?keyword=zzzz_not_exist_9x&page_size=50", timeout=10)).get("data") or {}
    ok = d.get("total") == 0 and d.get("list") == []
    rec("FT-14", M, "无结果关键词返回空列表", "GET /api/projects?keyword=zzzz_not_exist_9x",
        "total=0 且 list 为空", "total=%s" % d.get("total"), ok)

    d = body(s.get(BASE + "/api/search/suggest?keyword=' OR 1=1 --", timeout=10))
    ok = d.get("code") == 0
    rec("FT-15", M, "搜索输入 SQL 注入串不报错", "keyword=' OR 1=1 --",
        "服务正常返回，不出现 500", "code=%s" % d.get("code"), ok)

    # ───────────── C. 认证与权限 ─────────────
    M = "C 认证与权限"
    r = s.post(BASE + "/api/auth/login", json={"username": ADMIN, "password": PWD}, timeout=10)
    d = body(r).get("data") or {}
    token = d.get("token")
    ok = r.status_code == 200 and code(r) == 0 and token and len(token) > 20
    rec("FT-16", M, "正确账号密码登录成功", "POST /api/auth/login",
        "code=0 且返回 token", "code=%s token长度=%s" % (code(r), len(token) if token else 0), ok)

    r = s.post(BASE + "/api/auth/login", json={"username": ADMIN, "password": "wrong-pwd"}, timeout=10)
    ok = code(r) == 40100
    rec("FT-17", M, "错误密码登录被拒", "POST /api/auth/login (错误密码)",
        "code=40100", "code=%s" % code(r), ok)

    r = s.post(BASE + "/api/auth/login", json={"username": "no_such_user", "password": "x"}, timeout=10)
    ok = code(r) == 40100
    rec("FT-18", M, "不存在账号登录被拒", "POST /api/auth/login (不存在账号)",
        "code=40100", "code=%s" % code(r), ok)

    r = s.get(BASE + "/api/admin/overview", timeout=10)
    ok = code(r) == 40100
    rec("FT-19", M, "无 token 访问后台接口被拒", "GET /api/admin/overview (无 token)",
        "code=40100", "code=%s" % code(r), ok)

    r = s.get(BASE + "/api/admin/overview", headers={"Authorization": "Bearer fake.invalid.token"}, timeout=10)
    ok = code(r) == 40100
    rec("FT-20", M, "非法 token 被拒", "GET /api/admin/overview (伪造 token)",
        "code=40100", "code=%s" % code(r), ok)

    r = s.get(BASE + "/api/auth/me", headers=hdr(), timeout=10)
    d = body(r).get("data") or {}
    ok = code(r) == 0 and d.get("username") == ADMIN and "passwordHash" not in json.dumps(d)
    rec("FT-21", M, "当前用户接口返回信息且不含密码", "GET /api/auth/me",
        "username=admin 且无 passwordHash", "username=%s 含密码字段=%s" % (d.get("username"),
        "passwordHash" in json.dumps(d)), ok)

    # ───────────── D. 后台项目管理 ─────────────
    M = "D 项目管理"
    d = body(s.get(BASE + "/api/admin/overview", headers=hdr(), timeout=10)).get("data") or {}
    ok = d.get("projectTotal", 0) > 0 and "published" in d
    rec("FT-22", M, "仪表盘统计返回项目总数", "GET /api/admin/overview",
        "含 projectTotal/published", "projectTotal=%s" % d.get("projectTotal"), ok)

    d = body(s.get(BASE + "/api/admin/projects?page=1&page_size=5", headers=hdr(), timeout=10)).get("data") or {}
    ok = len(d.get("list", [])) == 5 and "status" in d["list"][0]
    rec("FT-23", M, "后台项目列表返回状态字段", "GET /api/admin/projects",
        "返回 5 条且每行含 status", "条数=%d" % len(d.get("list", [])), ok)

    pname = MARK + "自动化测试项目-" + "".join(random.choices(string.ascii_lowercase, k=5))
    payload = {"name": pname, "summary": "黑盒测试创建，测试结束删除", "categoryId": cid,
               "coverUrl": "", "status": "draft", "isRecommended": False, "recommendOrder": 0,
               "tags": [tags[0]["name"]] if tags else [],
               "members": [{"name": "测试员", "role": "QA"}],
               "blocks": [{"type": "heading", "text": "测试小节"},
                          {"type": "paragraph", "text": "这是自动化测试写入的段落。"},
                          {"type": "image", "url": "/uploads/cover/my-govoice.jpg", "caption": "测试配图"}]}
    r = s.post(BASE + "/api/admin/projects", headers=hdr(), json=payload, timeout=15)
    pid = (body(r).get("data") or {}).get("id")
    ok = code(r) == 0 and pid
    rec("FT-24", M, "新建项目（草稿）", "POST /api/admin/projects",
        "code=0 且返回新项目 id", "code=%s id=%s" % (code(r), pid), ok)

    if pid:
        d = body(s.get(BASE + "/api/admin/projects/%d" % pid, headers=hdr(), timeout=10)).get("data") or {}
        ok = len(d.get("blocks", [])) == 3 and len(d.get("members", [])) == 1 and d.get("status") == "draft"
        rec("FT-25", M, "新建后可回显内容块与成员", "GET /api/admin/projects/%d" % pid,
            "blocks=3 members=1 status=draft",
            "blocks=%d members=%d status=%s" % (len(d.get("blocks", [])), len(d.get("members", [])), d.get("status")), ok)

        r = s.get(BASE + "/api/projects/%d" % pid, timeout=10)
        ok = code(r) == 40400
        rec("FT-26", M, "草稿项目前台不可见", "GET /api/projects/%d (前台)" % pid,
            "code=40400", "code=%s" % code(r), ok)

        payload2 = dict(payload)
        payload2["status"] = "published"
        payload2["blocks"] = [{"type": "heading", "text": "改写后的小节"},
                              {"type": "paragraph", "text": "整单保存应全量替换旧块。"}]
        payload2["members"] = [{"name": "测试员二", "role": "QA"}]
        r = s.put(BASE + "/api/admin/projects/%d" % pid, headers=hdr(), json=payload2, timeout=15)
        d = body(s.get(BASE + "/api/admin/projects/%d" % pid, headers=hdr(), timeout=10)).get("data") or {}
        ok = code(r) == 0 and len(d.get("blocks", [])) == 2 and d.get("blocks")[0].get("text") == "改写后的小节"
        rec("FT-27", M, "整单保存全量替换内容块", "PUT /api/admin/projects/%d" % pid,
            "blocks 由 3 变 2 且内容已替换", "blocks=%d 首块=%s" % (len(d.get("blocks", [])),
            (d.get("blocks") or [{}])[0].get("text")), ok)

        r = s.get(BASE + "/api/projects/%d" % pid, timeout=10)
        ok = code(r) == 0
        rec("FT-28", M, "上架后前台可见", "PUT 置 published 后 GET 前台详情",
            "code=0 前台可见", "code=%s" % code(r), ok)

        r = s.put(BASE + "/api/admin/projects/%d/recommend" % pid, headers=hdr(),
                  json={"isRecommended": True, "recommendOrder": 9}, timeout=10)
        d = body(s.get(BASE + "/api/admin/projects/%d" % pid, headers=hdr(), timeout=10)).get("data") or {}
        ok = code(r) == 0 and d.get("isRecommended") is True and d.get("recommendOrder") == 9
        rec("FT-29", M, "设置推荐位与排序", "PUT /api/admin/projects/%d/recommend" % pid,
            "isRecommended=true recommendOrder=9",
            "isRecommended=%s order=%s" % (d.get("isRecommended"), d.get("recommendOrder")), ok)

        r = s.put(BASE + "/api/admin/projects/%d/status" % pid, headers=hdr(),
                  json={"status": "archived"}, timeout=10)
        d = body(s.get(BASE + "/api/admin/projects/%d" % pid, headers=hdr(), timeout=10)).get("data") or {}
        ok = code(r) == 0 and d.get("status") == "archived"
        rec("FT-30", M, "下架项目", "PUT /api/admin/projects/%d/status {archived}" % pid,
            "status=archived", "status=%s" % d.get("status"), ok)

        r = s.post(BASE + "/api/admin/projects", headers=hdr(), json=payload, timeout=10)
        ok = code(r) == 40900
        rec("FT-31", M, "项目重名被拒", "POST /api/admin/projects (同名)",
            "code=40900", "code=%s" % code(r), ok)

        r = s.delete(BASE + "/api/admin/projects/%d" % pid, headers=hdr(), timeout=10)
        r2 = s.get(BASE + "/api/admin/projects/%d" % pid, headers=hdr(), timeout=10)
        ok = code(r) == 0 and code(r2) == 40400
        rec("FT-32", M, "删除项目后查询返回 40400", "DELETE 后 GET 详情",
            "删除成功且后续 40400", "delete=%s 后续=%s" % (code(r), code(r2)), ok)

    r = s.get(BASE + "/api/admin/projects/999999", headers=hdr(), timeout=10)
    ok = code(r) == 40400
    rec("FT-33", M, "后台查询不存在项目返回 40400", "GET /api/admin/projects/999999",
        "code=40400", "code=%s" % code(r), ok)

    # ───────────── E. 内容块校验 ─────────────
    M = "E 内容块校验"
    def make_bad(blocks, tag):
        p = dict(payload)
        p["name"] = MARK + "块" + tag + "-" + "".join(random.choices(string.ascii_lowercase, k=6))
        p["blocks"] = blocks
        return p

    b1 = make_bad([{"type": "image", "url": ""}], "图片无URL")
    r = s.post(BASE + "/api/admin/projects", headers=hdr(), json=b1, timeout=10)
    if code(r) == 0:
        created_projects.append((body(r).get("data") or {}).get("id"))
    ok = code(r) == 40000
    rec("FT-34", M, "图片块缺 URL 被拒", "POST 项目 blocks=[{type:image,url:''}]",
        "code=40000", "code=%s（服务端未校验，项目被创建）" % code(r), ok)

    b2 = make_bad([{"type": "heading", "text": ""}], "标题无文本")
    r = s.post(BASE + "/api/admin/projects", headers=hdr(), json=b2, timeout=10)
    if code(r) == 0:
        created_projects.append((body(r).get("data") or {}).get("id"))
    ok = code(r) == 40000
    rec("FT-35", M, "标题块缺文本被拒", "POST 项目 blocks=[{type:heading,text:''}]",
        "code=40000", "code=%s（服务端未校验，项目被创建）" % code(r), ok)

    b3 = make_bad([{"type": "video", "text": "x"}], "非法类型")
    r = s.post(BASE + "/api/admin/projects", headers=hdr(), json=b3, timeout=10)
    if code(r) == 0:
        created_projects.append((body(r).get("data") or {}).get("id"))
    ok = code(r) == 40000
    rec("FT-36", M, "非法内容块类型被拒", "POST 项目 blocks=[{type:'video'}]",
        "code=40000", "code=%s" % code(r), ok)

    # ───────────── F. 分类与标签 ─────────────
    M = "F 分类与标签"
    cname = MARK + "测试分类-" + rnd()
    r = s.post(BASE + "/api/admin/categories", headers=hdr(),
               json={"name": cname, "description": "黑盒测试", "sortOrder": 99}, timeout=10)
    # 该接口成功时 data 为 null，需回查列表确认落库并取得 id
    _cats = body(s.get(BASE + "/api/admin/categories", headers=hdr(), timeout=10)).get("data") or []
    cid2 = next((c["id"] for c in _cats if c.get("name") == cname), None)
    ok = code(r) == 0 and cid2 is not None
    rec("FT-37", M, "新建分类（回查确认落库）", "POST /api/admin/categories + GET 列表",
        "code=0 且列表中可查到", "code=%s id=%s" % (code(r), cid2), ok)

    r = s.post(BASE + "/api/admin/categories", headers=hdr(),
               json={"name": cname, "description": "", "sortOrder": 0}, timeout=10)
    ok = code(r) == 40900
    rec("FT-38", M, "分类重名被拒", "POST /api/admin/categories (同名)",
        "code=40900", "code=%s" % code(r), ok)

    tname = MARK + "测试标签-" + rnd()
    r = s.post(BASE + "/api/admin/tags", headers=hdr(), json={"name": tname, "color": "#ff0000"}, timeout=10)
    _tags = body(s.get(BASE + "/api/admin/tags", headers=hdr(), timeout=10)).get("data") or []
    tid2 = next((t["id"] for t in _tags if t.get("name") == tname), None)
    ok = code(r) == 0 and tid2 is not None
    rec("FT-39", M, "新建标签（回查确认落库）", "POST /api/admin/tags + GET 列表",
        "code=0 且列表中可查到", "code=%s id=%s" % (code(r), tid2), ok)

    r = s.post(BASE + "/api/admin/tags", headers=hdr(), json={"name": tname, "color": ""}, timeout=10)
    ok = code(r) == 40900
    rec("FT-40", M, "标签重名被拒", "POST /api/admin/tags (同名)",
        "code=40900", "code=%s" % code(r), ok)

    # 被项目引用的分类不能删：categories[0] 下有演示项目
    r = s.delete(BASE + "/api/admin/categories/%d" % cid, headers=hdr(), timeout=10)
    ok = code(r) == 40900
    rec("FT-41", M, "被项目引用的分类不可删", "DELETE /api/admin/categories/<有项目的分类>",
        "code=40900 引用保护", "code=%s" % code(r), ok)

    r = s.delete(BASE + "/api/admin/categories/%d" % cid2, headers=hdr(), timeout=10) if cid2 else None
    ok = r is not None and code(r) == 0
    rec("FT-42", M, "删除无引用的分类成功", "DELETE 刚建的测试分类",
        "code=0", "code=%s" % (code(r) if r else None), ok)

    r = s.delete(BASE + "/api/admin/tags/%d" % tid2, headers=hdr(), timeout=10) if tid2 else None
    ok = r is not None and code(r) == 0
    rec("FT-43", M, "删除测试标签成功", "DELETE 刚建的测试标签",
        "code=0", "code=%s" % (code(r) if r else None), ok)

    # ───────────── G. 图片上传 ─────────────
    M = "G 图片上传"
    src = os.path.join(ROOT, "backend", "uploads", "cover", "my-govoice.jpg")
    if os.path.exists(src):
        with open(src, "rb") as fh:
            r = s.post(BASE + "/api/admin/upload", headers=hdr(),
                       files={"file": ("test.jpg", fh, "image/jpeg")}, timeout=20)
        url = (body(r).get("data") or {}).get("url") or (body(r).get("data") or {}).get("url")
        ok = code(r) == 0 and url
        if ok:
            r2 = s.get(BASE + url, timeout=10)
            ok = r2.status_code == 200
        rec("FT-44", M, "上传合法图片并可通过 URL 访问", "POST /api/admin/upload (jpg)",
            "code=0 且返回 url 可访问", "code=%s url=%s" % (code(r), url), ok)
    else:
        rec("FT-44", M, "上传合法图片并可通过 URL 访问", "POST /api/admin/upload (jpg)",
            "code=0 且返回 url 可访问", "跳过：样例图片不存在", False)

    bad_txt = os.path.join(HERE, "bad_upload.txt")
    big_jpg = os.path.join(HERE, "big_upload.jpg")
    with open(bad_txt, "wb") as fh:
        fh.write(b"not an image")
    with open(bad_txt, "rb") as fh:
        r = s.post(BASE + "/api/admin/upload", headers=hdr(),
                   files={"file": ("evil.txt", fh, "text/plain")}, timeout=10)
    ok = code(r) == 40000
    rec("FT-45", M, "上传非图片扩展名被拒", "POST /api/admin/upload (.txt)",
        "code=40000", "code=%s" % code(r), ok)

    with open(big_jpg, "wb") as fh:
        fh.write(b"\xff\xd8\xff" + os.urandom(6 * 1024 * 1024))
    with open(big_jpg, "rb") as fh:
        r = s.post(BASE + "/api/admin/upload", headers=hdr(),
                   files={"file": ("big.jpg", fh, "image/jpeg")}, timeout=30)
    ok = code(r) == 40000
    rec("FT-46", M, "上传超过 5MB 的图片被拒", "POST /api/admin/upload (6MB)",
        "code=40000", "code=%s" % code(r), ok)

    # ───────────── H. 异常与边界 ─────────────
    M = "H 异常与边界"
    r = s.get(BASE + "/api/projects/abc", timeout=10)
    ok = code(r) == 40000
    rec("FT-47", M, "非法 id 格式返回参数错误", "GET /api/projects/abc",
        "code=40000", "code=%s" % code(r), ok)

    long_name = MARK + "超长" * 200
    r = s.post(BASE + "/api/admin/projects", headers=hdr(),
               json={"name": long_name[:600], "summary": "x", "categoryId": cid,
                     "status": "draft", "blocks": []}, timeout=15)
    ok = code(r) == 40000
    rec("FT-48", M, "超长项目名被拒", "POST 项目 name 长度 > 100",
        "code=40000", "code=%s" % code(r), ok)

    r = s.post(BASE + "/api/admin/projects", headers=hdr(),
               json={"name": "", "summary": "x", "categoryId": cid, "status": "draft"}, timeout=10)
    ok = code(r) == 40000
    rec("FT-49", M, "项目名为空被拒", "POST 项目 name=''",
        "code=40000", "code=%s" % code(r), ok)

    r = s.post(BASE + "/api/admin/projects", headers=hdr(),
               json={"name": MARK + "不存在分类", "summary": "x", "categoryId": 999999,
                     "status": "draft"}, timeout=10)
    ok = code(r) == 40000
    rec("FT-50", M, "引用不存在的分类被拒", "POST 项目 categoryId=999999",
        "code=40000", "code=%s" % code(r), ok)

    # 软删除后以同名重建（对应 docs/03 §8 标记的"待确认"项）
    c3 = MARK + "软删重建-" + rnd()
    s.post(BASE + "/api/admin/categories", headers=hdr(),
           json={"name": c3, "description": "", "sortOrder": 0}, timeout=10)
    _cs = body(s.get(BASE + "/api/admin/categories", headers=hdr(), timeout=10)).get("data") or []
    c3id = next((c["id"] for c in _cs if c.get("name") == c3), None)
    if c3id:
        s.delete(BASE + "/api/admin/categories/%d" % c3id, headers=hdr(), timeout=10)
    r = s.post(BASE + "/api/admin/categories", headers=hdr(),
               json={"name": c3, "description": "", "sortOrder": 0}, timeout=10)
    ok = code(r) in (0, 40900)
    rec("FT-51", M, "软删除后同名资源可重建或明确报冲突",
        "新建分类 → 删除 → 再建同名", "code=0 或 40900（不应是 50000）",
        "code=%s" % code(r), ok)

    # ── 清理：删除本次创建的所有测试项目 + 临时文件 ──
    for pid_ in [p for p in created_projects if p]:
        s.delete(BASE + "/api/admin/projects/%d" % pid_, headers=hdr(), timeout=10)
    for f in (bad_txt, big_jpg):
        try:
            os.remove(f)
        except OSError:
            pass

    passed = sum(1 for x in RESULTS if x["status"] == "通过")
    summary = {"total": len(RESULTS), "passed": passed, "failed": len(RESULTS) - passed,
               "time": time.strftime("%Y-%m-%d %H:%M:%S"), "results": RESULTS}
    os.makedirs(RESULTS_DIR, exist_ok=True)
    with open(OUT, "w", encoding="utf-8") as fh:
        json.dump(summary, fh, ensure_ascii=False, indent=1)
    print("\n=== 合计 %d 条，通过 %d，未通过 %d ===" % (len(RESULTS), passed, len(RESULTS) - passed))


if __name__ == "__main__":
    main()
