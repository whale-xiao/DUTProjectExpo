# -*- coding: utf-8 -*-
"""简单性能测试：并发压测关键接口，采集 RPS / 平均 / P50 / P95 / P99 / 错误数。"""
import json, os, time, statistics, threading
from concurrent.futures import ThreadPoolExecutor
import requests

HERE = os.path.dirname(os.path.abspath(__file__))
RESULTS_DIR = os.path.join(HERE, "results")
BASE = os.environ.get("BASE_URL", "http://127.0.0.1").rstrip("/")
OUT = os.path.join(RESULTS_DIR, "perf_results.json")
_local = threading.local()


def sess():
    if not hasattr(_local, "s"):
        _local.s = requests.Session()
    return _local.s


def one(method, url, body=None):
    t0 = time.perf_counter()
    try:
        if method == "GET":
            r = sess().get(BASE + url, timeout=15)
        else:
            r = sess().post(BASE + url, json=body, timeout=15)
        ok = r.status_code == 200
        n = len(r.content)
    except Exception:
        ok, n = False, 0
    return (time.perf_counter() - t0) * 1000.0, ok, n


def run(name, desc, method, url, n, conc, body=None):
    # 预热，避免首请求建连影响结果
    with ThreadPoolExecutor(max_workers=conc) as ex:
        list(ex.map(lambda _: one(method, url, body), range(min(20, n))))

    lat, errs, sizes = [], 0, []
    lock = threading.Lock()

    def job(_):
        nonlocal errs
        ms, ok, size = one(method, url, body)
        with lock:
            lat.append(ms)
            sizes.append(size)
            if not ok:
                errs += 1

    t0 = time.perf_counter()
    with ThreadPoolExecutor(max_workers=conc) as ex:
        list(ex.map(job, range(n)))
    wall = time.perf_counter() - t0

    lat.sort()
    p = lambda q: lat[min(len(lat) - 1, int(len(lat) * q))]
    res = {"name": name, "desc": desc, "method": method, "url": url,
           "requests": n, "concurrency": conc, "wall_s": round(wall, 3),
           "rps": round(n / wall, 1), "avg_ms": round(statistics.mean(lat), 2),
           "p50_ms": round(p(0.50), 2), "p95_ms": round(p(0.95), 2),
           "p99_ms": round(p(0.99), 2), "max_ms": round(lat[-1], 2),
           "errors": errs, "avg_size_kb": round(statistics.mean(sizes) / 1024, 1)}
    print("  %-26s %3d 请求 并发%2d  %6.1f RPS  平均%6.1fms  P95%6.1fms  错误%d"
          % (name, n, conc, res["rps"], res["avg_ms"], res["p95_ms"], errs))
    return res


def main():
    print("=== 性能测试开始 ===")
    scen = [
        ("首页聚合 API", "GET /api/home，首页首屏主接口", "GET", "/api/home", 300, 10, None),
        ("项目列表 API", "GET /api/projects，分页 12 条", "GET", "/api/projects?page=1&page_size=12", 400, 20, None),
        ("项目详情 API", "GET /api/projects/1，含内容块/成员/标签", "GET", "/api/projects/1", 400, 20, None),
        ("搜索联想 API", "GET /api/search/suggest，键入 Go", "GET", "/api/search/suggest?keyword=Go", 300, 10, None),
        ("静态首页 HTML", "GET /，nginx 直出 SPA 入口", "GET", "/", 400, 20, None),
        ("静态图片", "GET /uploads/cover/my-govoice.jpg", "GET", "/uploads/cover/my-govoice.jpg", 400, 20, None),
        ("登录接口", "POST /api/auth/login（含 bcrypt 校验）", "POST", "/api/auth/login", 60, 3,
         {"username": "admin", "password": "admin123456"}),
    ]
    results = []
    for s in scen:
        results.append(run(*s))

    summary = {"time": time.strftime("%Y-%m-%d %H:%M:%S"),
               "env": "Docker Compose（nginx + backend + mysql 8.0），本机 Windows 11",
               "data": "15 个项目 / 4 分类 / 30 标签", "results": results}
    os.makedirs(RESULTS_DIR, exist_ok=True)
    with open(OUT, "w", encoding="utf-8") as fh:
        json.dump(summary, fh, ensure_ascii=False, indent=1)
    print("=== 完成，结果写入 %s ===" % OUT)


if __name__ == "__main__":
    main()
