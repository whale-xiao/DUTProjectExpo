-- 清空演示数据脚本（seed.go 的注释引用本文件）
--
-- 用途：SeedDemo 只在「项目中表为空」时写入演示数据。若想重新灌入演示数据
--       （例如演示项目清单有更新），先执行本脚本清空，再重启后端 `go run .`，
--       启动时会自动重新建表并写入最新的一批演示数据。
--
-- 执行：
--   mysql -uroot -p < backend/scripts/drop.sql
--
-- ⚠️ 会删除全部业务数据（项目、内容块、成员、图片、标签、分类、访问日志），
--    不可恢复。仅用于开发/演示环境；生产环境请勿执行。
--    管理员账号（users 表）保留，不删除。

USE project_showcase;

-- 关闭外键检查，避免删除顺序问题
SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE project_content_blocks;
TRUNCATE TABLE project_members;
TRUNCATE TABLE project_images;
TRUNCATE TABLE project_tags;
TRUNCATE TABLE visit_logs;
TRUNCATE TABLE projects;
TRUNCATE TABLE tags;
TRUNCATE TABLE categories;

SET FOREIGN_KEY_CHECKS = 1;

SELECT
  (SELECT COUNT(*) FROM projects) AS 剩余项目,
  (SELECT COUNT(*) FROM categories) AS 剩余分类,
  (SELECT COUNT(*) FROM tags) AS 剩余标签;

-- 执行完成后重启后端：cd backend && go run .
-- 预期日志：[seed] ✅ 已写入演示数据：分类 4、标签 7、示例项目 15
