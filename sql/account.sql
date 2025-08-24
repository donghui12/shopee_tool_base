-- 创建 account 表
CREATE TABLE IF NOT EXISTS `account` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `account_id` bigint NOT NULL,
    `username` varchar(255) NOT NULL COMMENT '虾皮账户名',
    `password` varchar(255) DEFAULT NULL COMMENT '虾皮密码',
    `phone` varchar(255) DEFAULT NULL COMMENT '手机号',
    `machine_code` varchar(255) DEFAULT NULL COMMENT '机器码',
    `active_code` varchar(255) DEFAULT NULL COMMENT '激活码',
    `expired_at` varchar(255) DEFAULT NULL COMMENT '过期时间',
    `cookies` text COMMENT 'cookies信息',
    `session` text COMMENT 'session信息',
    `status` int NOT NULL DEFAULT '1' COMMENT '状态：1=有效 0=无效',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_account_id` (`account_id`),
    KEY `idx_username` (`username`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='虾皮账号信息表';

ALTER TABLE `accounts`
ADD COLUMN `merchant_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'shopee 账户' AFTER `account_id`;