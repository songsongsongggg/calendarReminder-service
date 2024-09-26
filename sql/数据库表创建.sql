CREATE
DATABASE `calendar_reminder`;
USE
`calendar_reminder`;

-- 删除 users 表，如果存在
DROP TABLE IF EXISTS users;
-- 创建 users 表
CREATE TABLE users
(
    id         INT PRIMARY KEY AUTO_INCREMENT COMMENT '唯一标识用户的ID',
    mobile     VARCHAR(20)  NOT NULL UNIQUE COMMENT '用户手机号，国内格式',
    creator_id VARCHAR(128) NOT NULL UNIQUE COMMENT '用户唯一标识，用于关联提醒信息',
    created_at DATETIME NOT NULL COMMENT '用户创建时间',
    updated_at DATETIME NOT NULL COMMENT '用户信息最后更新时间'
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 删除 reminders 表，如果存在
DROP TABLE IF EXISTS reminders;
-- 创建 reminders 表
CREATE TABLE reminders
(
    id         INT PRIMARY KEY AUTO_INCREMENT COMMENT '唯一标识提醒信息的ID',
    creator_id VARCHAR(128) NOT NULL COMMENT '提醒信息创建者的ID',
    content    TEXT         NOT NULL COMMENT '提醒内容',
    remind_at  DATETIME     NOT NULL COMMENT '提醒的具体时间', -- 将 TIMESTAMP 改为 DATETIME
    created_at DATETIME     NOT NULL  COMMENT '提醒信息创建时间', -- 改为 DATETIME
    updated_at DATETIME     NOT NULL  COMMENT '提醒信息最后更新时间', -- 改为 DATETIME
    INDEX      idx_creator_id (creator_id(20)) -- 只索引前 20 个字符
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;