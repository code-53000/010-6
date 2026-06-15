-- 榨油坊柜台管理系统 - 数据库初始化脚本
CREATE DATABASE IF NOT EXISTS oilpress DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE oilpress;

-- 价格配置（初始值，可在前端页面修改）
INSERT INTO price_config (id, price_per_kg, oil_rate, processing_fee, cake_take_fee, updated_at) VALUES
('peanut',    30.00, 45.0, 50.00, 10.00, NOW()),
('rapeseed',  25.00, 38.0, 50.00,  8.00, NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();
