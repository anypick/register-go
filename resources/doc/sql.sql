CREATE TABLE `site` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `site_id` bigint(20) NOT NULL COMMENT '门店id',
  `site_code` varchar(10) NOT NULL DEFAULT '' COMMENT '门店编码',
  `site_name` varchar(240) NOT NULL DEFAULT '' COMMENT '门店名称',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `site_id` (`site_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10001 DEFAULT CHARSET=utf8;