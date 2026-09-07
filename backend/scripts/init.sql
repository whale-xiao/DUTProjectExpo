-- 建库脚本（对齐 .env.example 的 DB_NAME=project_showcase）
-- 执行：mysql -uroot -p < backend/scripts/init.sql
-- 说明：utf8mb4_general_ci 在本地 5.7 与云端 8.0 均兼容；
--       数据表无需手写 DDL，后端启动时由 GORM AutoMigrate 自动创建（对应 model/ 9 张表）。
CREATE DATABASE IF NOT EXISTS `project_showcase`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;
