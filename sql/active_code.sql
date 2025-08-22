-- 创建 active_code 表
CREATE TABLE IF NOT EXISTS `active_code` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `code` varchar(255) NOT NULL COMMENT '激活码',
    `expired_at` timestamp NOT NULL COMMENT '过期时间',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_expired_at` (`expired_at`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='激活码表';