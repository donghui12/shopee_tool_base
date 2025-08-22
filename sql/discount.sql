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