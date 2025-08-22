-- 初始化数据库表结构
-- github.com/donghui12/shopee_tool_base 项目所需的所有表

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

-- 创建 shop 表
CREATE TABLE IF NOT EXISTS `shop` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `shop_id` varchar(64) NOT NULL COMMENT '店铺ID',
    `region` varchar(16) NOT NULL COMMENT '地区',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_shop_id` (`shop_id`),
    KEY `idx_region` (`region`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='店铺信息表';

-- 创建 discount 表
CREATE TABLE IF NOT EXISTS `discount` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `name` varchar(255) DEFAULT NULL COMMENT '折扣名称',
    `shop_id` varchar(64) NOT NULL COMMENT '店铺ID',
    `discount_id` bigint NOT NULL COMMENT '折扣ID',
    `status` int NOT NULL DEFAULT '1' COMMENT '状态',
    `start_time` varchar(255) DEFAULT NULL COMMENT '开始时间',
    `end_time` varchar(255) DEFAULT NULL COMMENT '结束时间',
    `item_count` int DEFAULT NULL COMMENT '商品数量',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_discount_id` (`discount_id`),
    KEY `idx_shop_id` (`shop_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='折扣信息表';

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