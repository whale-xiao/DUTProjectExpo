# 演示图片来源

> 这些图仅作**原型/演示占位**，正式内容请替换为实训基地真实项目截图。
> 统一规格：**1200 × 800（3:2）JPEG**，与 `docs/09-前台设计方案.md` §5.3 的封面规格一致。

## 一、Pexels（免费商用许可）

> 许可：Pexels License —— 免费使用、可商用、无需署名。
> 来源页：https://www.pexels.com/photo/&lt;photo-slug&gt;-&lt;id&gt;/

| 文件 | 对应 | Pexels Photo | 原图页路径片段 |
| --- | --- | --- | --- |
| my-govoice.jpg | MyGoVoice 封面 | 177598 | black-laptop-computer-turned-on-showing-computer-codes-177598 |
| my-govoice-shot.jpg | MyGoVoice 正文截图占位 | 34803994 | modern-laptop-on-wooden-desk-with-code-displayed-34803994 |
| ai-qa.jpg | 智能问答平台 封面 | 30530411 | laptop-displaying-ai-conversation-interface-30530411 |
| ai-qa-shot.jpg | 智能问答平台 正文截图占位 | 30530403 | laptop-screen-displaying-ai-interface-at-night-30530403 |
| iot-monitor.jpg | IoT 环境监测站 封面 | 35686441 | macro-shot-of-electronic-sensor-module-on-yellow-background-35686441 |

下载 URL 模板：`https://images.pexels.com/photos/<id>/pexels-photo-<id>.jpeg?auto=compress&cs=tinysrgb&w=1200&h=800&fit=crop`

## 二、Lorem Picsum（其余演示项目封面）

`demo-01.jpg` ~ `demo-12.jpg` 用于扩充演示项目（见 `repository/seed.go`）。
这 12 张是**通用随机照片**，与项目主题无关，仅用于让首页在演示数据下有足够内容量
（撑起「热力榜」「分类合集行」等板块），**必须替换为真实截图**。

> 许可：Lorem Picsum 图片来自 Unsplash，遵循 Unsplash License —— 免费使用、可商用、无需署名。
> 来源：https://picsum.photos/

| 文件 | 用于（seed 演示项目） | 真实内容应替换为 |
| --- | --- | --- |
| demo-01.jpg | 校园二手交易平台 | 平台首页/商品列表截图 |
| demo-02.jpg | 分布式文件同步工具 | 集群状态或同步日志界面 |
| demo-03.jpg | 代码查重引擎 | 查重报告对比界面 |
| demo-04.jpg | 手写公式识别 | 公式识别前后对比图 |
| demo-05.jpg | 校园舆情分析 | 情感趋势与话题聚类图 |
| demo-06.jpg | 零件缺陷检测 | 缺陷标注框选效果图 |
| demo-07.jpg | 智能宿舍门禁 | 门禁终端实拍 |
| demo-08.jpg | 实验室环境监控 | 监控看板或硬件接线图 |
| demo-09.jpg | 智能小车导航 | 小车实拍 + 建图界面 |
| demo-10.jpg | 课程表助手 | 课表主界面 |
| demo-11.jpg | 数据可视化大屏 | 大屏实拍 |
| demo-12.jpg | 校园地图导览 | 地图导览界面 |

重新下载命令（如需换图）：

```bash
curl -L "https://picsum.photos/1200/800?random=<种子>" -o demo-NN.jpg
```
