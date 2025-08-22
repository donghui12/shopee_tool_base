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

ALTER TABLE shop 
ADD COLUMN `account_id` bigint COMMENT 'shopee account id' AFTER `id`;