/*
 Navicat Premium Data Transfer

 Source Server         : 本地
 Source Server Type    : MariaDB
 Source Server Version : 100407
 Source Host           : localhost:3308
 Source Schema         : mcpay

 Target Server Type    : MariaDB
 Target Server Version : 100407
 File Encoding         : 65001

 Date: 14/07/2025 23:40:10
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for address
-- ----------------------------
DROP TABLE IF EXISTS `address`;
CREATE TABLE `address`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `network` int(11) NOT NULL DEFAULT 0 COMMENT '网络：1:波场；2:币安',
  `blockNum` int(11) NOT NULL DEFAULT 0 COMMENT '最新区块高度',
  `path_index` int(11) NOT NULL DEFAULT 1,
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '地址',
  `balance` decimal(36, 18) NOT NULL DEFAULT 0.000000000000000000 COMMENT '剩余额度',
  `private_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '私钥',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1:可用；2:使用中；3:禁用',
  `used_times` int(11) NOT NULL DEFAULT 0 COMMENT '使用次数',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `last_used_at` datetime NULL DEFAULT NULL COMMENT '最后使用时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_address`(`address`) USING BTREE,
  UNIQUE INDEX `uk_network_path_index`(`network`, `path_index`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 41 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '钱包地址' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of address
-- ----------------------------
INSERT INTO `address` VALUES (1, 1, 0, 1, 'TKeBRQDnfBiXDg2Ts89HtXt4EF15EJyGo9', 0.000000000000000000, '', 1, 10, '2025-06-26 00:53:00', '2025-07-11 06:52:07', '2025-07-11 00:34:27');
INSERT INTO `address` VALUES (2, 1, 0, 2, 'TVV9MXpDgyRJG4qArA8mGy3gaGEVHHfFiC', 0.000000000000000000, '', 1, 7, '2025-06-26 00:53:00', '2025-07-06 18:09:53', '2025-07-06 18:09:53');
INSERT INTO `address` VALUES (3, 1, 0, 3, 'TSVfF9g7f5hCkfajMx1Q2f8XWjET22TUmB', 0.000000000000000000, '', 1, 4, '2025-06-26 00:53:00', '2025-07-06 14:05:09', '2025-07-06 13:53:46');
INSERT INTO `address` VALUES (4, 1, 0, 4, 'TCg4BG5ueCZmkCaB7kx2LkH7NfVyYy4Y9d', 0.000000000000000000, '', 1, 3, '2025-06-26 00:53:00', '2025-07-06 14:09:15', '2025-07-06 13:58:50');
INSERT INTO `address` VALUES (5, 1, 0, 5, 'TH6YHZ9MHEMGBnPGA23qrRkmskfJtDtBcA', 0.000000000000000000, '', 1, 3, '2025-06-26 00:53:00', '2025-07-06 14:14:15', '2025-07-06 14:04:15');
INSERT INTO `address` VALUES (6, 1, 0, 6, 'TMqyzaQrTMtd4ZnW6iu7xrpfc9EooGvg4D', 0.000000000000000000, '', 1, 1, '2025-06-26 00:53:00', '2025-07-06 14:15:15', '2025-07-06 14:04:30');
INSERT INTO `address` VALUES (7, 1, 0, 7, 'TMSwPbmT3hF8eRTMsJU8a4kkZ2F1dRQvGg', 0.000000000000000000, '', 1, 0, '2025-06-26 00:53:00', '2025-06-26 00:53:00', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (8, 1, 0, 8, 'TNMFdgjdWyTvc5eL7QLEbRDdeZqn9NThEW', 0.000000000000000000, '', 1, 0, '2025-06-26 00:53:00', '2025-06-26 00:53:00', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (9, 1, 0, 9, 'TP1XBFxxk12FT4tYkyyjLoryLnS9C4CNYT', 0.000000000000000000, '', 1, 0, '2025-06-26 00:53:00', '2025-06-26 00:53:00', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (10, 1, 0, 10, 'TXrzepR8um3gfAQ5eg7J3axE2rCUHvfcyb', 0.000000000000000000, '', 1, 0, '2025-06-26 00:53:00', '2025-06-26 00:53:00', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (11, 1, 0, 11, 'TVUGQ3ZiawmbKpbWWRVswo9SNebr96oQ2M', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:15', '2025-07-01 00:21:15', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (12, 1, 0, 12, 'TTnkHoVwey3rUmVZmPzkeBqkmudkuLG4Kc', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:15', '2025-07-01 00:21:15', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (13, 1, 0, 13, 'TQSysGudmYkBzgPU3whJuJcDPC5yEXzCgh', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:15', '2025-07-01 00:21:15', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (14, 1, 0, 14, 'TPuFu7mcLgy1i9scMEBtRbabgk394cDLpA', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:15', '2025-07-01 00:21:15', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (15, 1, 0, 15, 'TDUvQDAcjcXsRadibHD8NX1UHBRm8rH3Wv', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:15', '2025-07-01 00:21:15', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (16, 1, 0, 16, 'TLdZkLcvaWJCFzr27vBJHpxEboo5Wd6jSk', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:15', '2025-07-01 00:21:15', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (17, 1, 0, 17, 'TZFrkunuh5MR1VafYfgvEEZxtHVq5yUL4V', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:15', '2025-07-01 00:21:15', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (18, 1, 0, 18, 'TQvYoqKaxWqkq5XANnzPmYKuN17aHumT9D', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:15', '2025-07-01 00:21:15', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (19, 1, 0, 19, 'TDVABmkBhwEfaX5QaRss93xLomTMTaGVSB', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:15', '2025-07-01 00:21:15', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (20, 1, 0, 20, 'TS7D67GBERVYN7WPpXVguS3xyd3cQB2zeE', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:15', '2025-07-01 00:21:15', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (21, 1, 0, 21, 'TTzBEzimM5Fun97RyM2882zCZZb2K4V8sD', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:27', '2025-07-01 00:21:27', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (22, 1, 0, 22, 'TEoA8849puKfjk93JLwJLj2qy9m1rZr8tk', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:27', '2025-07-01 00:21:27', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (23, 1, 0, 23, 'TTVcj1iDeCza1LUZYcjy2e6KMdy3sVDgvo', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:27', '2025-07-01 00:21:27', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (24, 1, 0, 24, 'TW8565rtCcHxk3xb7zLt5sCrGdVXDVpqSw', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:27', '2025-07-01 00:21:27', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (25, 1, 0, 25, 'TLTrBAzngWKEBwCvRfM2V646RWBqHsrRVr', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:27', '2025-07-01 00:21:27', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (26, 1, 0, 26, 'TZF35jkQUywD2ArP7gzkvWCkoovfE4kNQP', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:27', '2025-07-01 00:21:27', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (27, 1, 0, 27, 'TTuzYXFwKCX3aqyqFkpNHzDxYNUGohx1fN', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:27', '2025-07-01 00:21:27', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (28, 1, 0, 28, 'TKokRGUMFhW8f9rLG8fHcyQPwEjg5hLNTp', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:27', '2025-07-01 00:21:27', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (29, 1, 0, 29, 'TF3mnL1jWQ5ZkkR8YdyUN8BV15A6x3jP6M', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:27', '2025-07-01 00:21:27', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (30, 1, 0, 30, 'TLQNe8aqMQHJ5DKTayVfpzHbVZJ4LXKLEa', 0.000000000000000000, '', 1, 0, '2025-07-01 00:21:27', '2025-07-01 00:21:27', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (31, 2, 0, 1, '0x3D1B17B906F97c87A345229f02abAb96Fb27c52b', 0.000000000000000000, '', 1, 3, '2025-07-11 00:41:12', '2025-07-14 23:17:30', '2025-07-14 23:02:31');
INSERT INTO `address` VALUES (32, 2, 0, 2, '0xEbAa7906fb8011e8dAcd6944d18D2A357BF3b1dC', 0.000000000000000000, '', 1, 2, '2025-07-11 00:41:12', '2025-07-14 23:17:30', '2025-07-14 23:03:25');
INSERT INTO `address` VALUES (33, 2, 0, 3, '0xe8cD9b1292b8E1181184754d094ed57C0d3D59C5', 0.000000000000000000, '', 1, 0, '2025-07-11 00:41:12', '2025-07-11 00:41:12', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (34, 2, 0, 4, '0x2620FFD01f70c2B6A3Bc655F92356FdbEEeE05DD', 0.000000000000000000, '', 1, 0, '2025-07-11 00:41:12', '2025-07-11 00:41:12', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (35, 2, 0, 5, '0xD23B63172125C9Eb65F33a7608125F524Fdc28e8', 0.000000000000000000, '', 1, 0, '2025-07-11 00:41:12', '2025-07-11 00:41:12', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (36, 2, 0, 6, '0x1A1B13a31B11dCa9c4183c45805c0f56d7d442A7', 0.000000000000000000, '', 1, 0, '2025-07-11 00:41:12', '2025-07-11 00:41:12', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (37, 2, 0, 7, '0x2E16faD3BC199038434D78205E08066EdE989DA4', 0.000000000000000000, '', 1, 0, '2025-07-11 00:41:12', '2025-07-11 00:41:12', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (38, 2, 0, 8, '0xc896d8d351678c714A33a3083c4382bCE8b01052', 0.000000000000000000, '', 1, 0, '2025-07-11 00:41:12', '2025-07-11 00:41:12', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (39, 2, 0, 9, '0xBC2ef736522A303d78A0e7e5095625BD09e46dba', 0.000000000000000000, '', 1, 0, '2025-07-11 00:41:12', '2025-07-11 00:41:12', '0000-00-00 00:00:00');
INSERT INTO `address` VALUES (40, 2, 0, 10, '0xDd32178da965653A794cCB5862aF155014d31eF4', 0.000000000000000000, '', 1, 0, '2025-07-11 00:41:12', '2025-07-11 00:41:12', '0000-00-00 00:00:00');

-- ----------------------------
-- Table structure for apps
-- ----------------------------
DROP TABLE IF EXISTS `apps`;
CREATE TABLE `apps`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '0' COMMENT '名称',
  `app_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `app_secret` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1:正常；2:禁用',
  `pay_channel` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '支付渠道',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_app_id`(`app_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of apps
-- ----------------------------
INSERT INTO `apps` VALUES (1, '测试', '654232', '85451FE8F6AFBF3705498D10E36B47F0', 1, '2', '2025-06-08 10:14:20', NULL);

-- ----------------------------
-- Table structure for configs
-- ----------------------------
DROP TABLE IF EXISTS `configs`;
CREATE TABLE `configs`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `k` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `v` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  `beizhu` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of configs
-- ----------------------------
INSERT INTO `configs` VALUES (1, 'mnemonic', 'umbrella moral birth roof ceiling basket lesson rail burger region coconut pass involve bacon alpha chuckle pattern brush autumn oppose damp wire praise image', '助记词');
INSERT INTO `configs` VALUES (2, 'order_code_salt', 'ahd28fsf', NULL);
INSERT INTO `configs` VALUES (3, 'etherscan_apikey', 'TQZBWCQJQP5H8CC2JWKFFA7MHFAXK2M1UC', '文档：https://docs.etherscan.io/etherscan-v2/api-endpoints/accounts#get-a-list-of-erc20-token-transfer-events-by-address');

-- ----------------------------
-- Table structure for order_address
-- ----------------------------
DROP TABLE IF EXISTS `order_address`;
CREATE TABLE `order_address`  (
  `order_id` int(11) NOT NULL,
  `address_id` int(11) NOT NULL COMMENT '地址ID',
  `block_num` int(11) NOT NULL DEFAULT 0 COMMENT '区块高度',
  PRIMARY KEY (`order_id`, `address_id`) USING BTREE,
  INDEX `pay_channel_id`(`address_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '订单表和地址绑定关系表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of order_address
-- ----------------------------
INSERT INTO `order_address` VALUES (542, 5, 0);
INSERT INTO `order_address` VALUES (545, 1, 0);
INSERT INTO `order_address` VALUES (547, 2, 0);
INSERT INTO `order_address` VALUES (550, 1, 0);
INSERT INTO `order_address` VALUES (552, 2, 0);
INSERT INTO `order_address` VALUES (553, 4, 0);
INSERT INTO `order_address` VALUES (554, 5, 0);
INSERT INTO `order_address` VALUES (556, 1, 0);
INSERT INTO `order_address` VALUES (557, 2, 0);
INSERT INTO `order_address` VALUES (558, 1, 0);
INSERT INTO `order_address` VALUES (559, 2, 0);
INSERT INTO `order_address` VALUES (561, 3, 0);
INSERT INTO `order_address` VALUES (562, 1, 0);
INSERT INTO `order_address` VALUES (563, 2, 0);
INSERT INTO `order_address` VALUES (564, 3, 0);
INSERT INTO `order_address` VALUES (565, 4, 0);
INSERT INTO `order_address` VALUES (566, 1, 0);
INSERT INTO `order_address` VALUES (567, 5, 0);
INSERT INTO `order_address` VALUES (568, 6, 0);
INSERT INTO `order_address` VALUES (569, 1, 0);
INSERT INTO `order_address` VALUES (570, 1, 0);
INSERT INTO `order_address` VALUES (571, 2, 0);
INSERT INTO `order_address` VALUES (583, 1, 0);
INSERT INTO `order_address` VALUES (588, 31, 57780338);
INSERT INTO `order_address` VALUES (589, 31, 57810120);
INSERT INTO `order_address` VALUES (590, 32, 57810877);
INSERT INTO `order_address` VALUES (591, 31, 58232935);
INSERT INTO `order_address` VALUES (592, 32, 58233006);

-- ----------------------------
-- Table structure for orders
-- ----------------------------
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NULL DEFAULT NULL,
  `app_id` int(11) NOT NULL,
  `order_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '订单号',
  `amount` decimal(10, 2) UNSIGNED NOT NULL DEFAULT 0.00 COMMENT '订单金额',
  `amount_true` decimal(10, 2) UNSIGNED NOT NULL DEFAULT 0.00 COMMENT '到账金额',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `callback_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '通知地址',
  `callback_state` int(11) NOT NULL DEFAULT 0 COMMENT '通知状态：0:失败；1:成功；',
  `callback_err` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '通知报错信息',
  `callback_date` datetime NULL DEFAULT NULL COMMENT '最后通知时间',
  `callback_times` int(11) NOT NULL DEFAULT 0 COMMENT '通知次数',
  `merchant_order_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '商户订单号',
  `third_order_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '第三方订单号',
  `status` int(11) NULL DEFAULT 1 COMMENT '1:待支付；2:成功；3:失败',
  `redirect_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '支付成功跳转地址',
  `updated_at` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `order_id`(`order_id`) USING BTREE,
  UNIQUE INDEX `app_id_merchant_order_id_unique`(`app_id`, `merchant_order_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 593 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '订单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of orders
-- ----------------------------
INSERT INTO `orders` VALUES (542, NULL, 654232, '2025070106590032635956', 100.00, 0.00, '2025-07-01 06:59:00', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432543', '', 3, '', '2025-07-02 23:43:00');
INSERT INTO `orders` VALUES (545, NULL, 654232, '2025070423520593569681', 100.00, 0.00, '2025-07-04 23:52:05', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432541', '', 3, '', '2025-07-05 00:02:28');
INSERT INTO `orders` VALUES (547, NULL, 654232, '2025070423540298636433', 100.00, 0.00, '2025-07-04 23:54:02', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432542', '', 3, '', '2025-07-05 00:04:28');
INSERT INTO `orders` VALUES (550, NULL, 654232, '2025070500125990900519', 100.00, 0.00, '2025-07-05 00:12:59', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432544', '', 3, '', '2025-07-05 00:23:28');
INSERT INTO `orders` VALUES (552, NULL, 654232, '2025070500134112171528', 100.00, 0.00, '2025-07-05 00:13:41', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432545', '', 3, '', '2025-07-05 00:24:28');
INSERT INTO `orders` VALUES (553, NULL, 654232, '20250705001418585589', 100.00, 0.00, '2025-07-05 00:14:18', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432546', '', 3, '', '2025-07-05 00:24:28');
INSERT INTO `orders` VALUES (554, NULL, 654232, '2025070500235232750714', 100.00, 0.00, '2025-07-05 00:23:52', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432547', '', 3, '', '2025-07-05 00:34:28');
INSERT INTO `orders` VALUES (556, 0, 654232, '2025070612552212245093', 100.00, 100.00, '2025-07-06 12:55:22', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432548', '', 3, '', '2025-07-06 13:06:04');
INSERT INTO `orders` VALUES (557, 0, 654232, '2025070613031021761490', 10.00, 10.00, '2025-07-06 13:03:10', 'http://callbackurl.com', 0, '', '2025-07-06 13:04:20', 1, 'Ad3232432549', '', 2, '', '2025-07-06 13:04:20');
INSERT INTO `orders` VALUES (558, 10000, 654232, '2025070613284736154279', 10.00, 10.00, '2025-07-06 13:28:47', 'http://callbackurl.com', 0, '', '2025-07-06 13:31:21', 1, 'Ad3232432550', '', 3, '', '2025-07-06 13:39:34');
INSERT INTO `orders` VALUES (559, 10000, 654232, '202507061334014998578', 10.00, 10.00, '2025-07-06 13:34:01', 'http://callbackurl.com', 0, '', '2025-07-06 13:35:54', 1, 'Ad3232432551', '', 3, '', '2025-07-06 13:44:34');
INSERT INTO `orders` VALUES (561, 10000, 654232, '202507061338112015552', 10.00, 10.00, '2025-07-06 13:38:11', 'http://callbackurl.com', 0, '', '2025-07-06 13:41:34', 1, 'Ad3232432552', '', 3, '', '2025-07-06 13:49:24');
INSERT INTO `orders` VALUES (562, 10000, 654232, '2025070613480317798945', 10.00, 10.00, '2025-07-06 13:48:03', 'http://callbackurl.com', 0, '', '2025-07-06 13:48:39', 1, 'Ad3232432553', '', 3, '', '2025-07-06 13:58:16');
INSERT INTO `orders` VALUES (563, 10000, 654232, '2025070613503691181511', 10.00, 10.00, '2025-07-06 13:50:36', 'http://callbackurl.com', 0, '', '2025-07-06 13:51:20', 1, 'Ad3232432554', '', 3, '', '2025-07-06 14:01:16');
INSERT INTO `orders` VALUES (564, 10000, 654232, '20250706135346534990', 10.00, 10.00, '2025-07-06 13:53:46', 'http://callbackurl.com', 0, '', '2025-07-06 13:54:29', 1, 'Ad3232432555', '', 3, '', '2025-07-06 14:05:09');
INSERT INTO `orders` VALUES (565, 10000, 654232, '2025070613585011380817', 10.00, 0.00, '2025-07-06 13:58:50', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432556', '', 3, '', '2025-07-06 14:09:15');
INSERT INTO `orders` VALUES (566, 10000, 654232, '2025070613591750108702', 10.00, 10.00, '2025-07-06 13:59:17', 'http://callbackurl.com', 0, '', '2025-07-06 13:59:46', 1, 'Ad3232432557', '', 2, '', '2025-07-06 13:59:46');
INSERT INTO `orders` VALUES (567, 10000, 654232, '2025070614041530256676', 10.00, 0.00, '2025-07-06 14:04:15', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432558', '', 3, '', '2025-07-06 14:14:15');
INSERT INTO `orders` VALUES (568, 10000, 654232, '2025070614043096008459', 10.00, 0.00, '2025-07-06 14:04:30', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432559', '', 3, '', '2025-07-06 14:15:15');
INSERT INTO `orders` VALUES (569, 10000, 654232, '2025070614050460245871', 10.00, 10.00, '2025-07-06 14:05:04', 'http://callbackurl.com', 0, '', '2025-07-06 14:05:45', 1, 'Ad3232432560', '', 2, '', '2025-07-06 14:05:45');
INSERT INTO `orders` VALUES (570, 10000, 654232, '20250706180324203259', 10.00, 10.00, '2025-07-06 18:03:24', 'http://callbackurl.com', 0, '', '2025-07-06 18:04:05', 1, 'Ad3232432561', '', 2, '', '2025-07-06 18:04:05');
INSERT INTO `orders` VALUES (571, 10000, 654232, '2025070618095375601505', 10.00, 10.00, '2025-07-06 18:09:53', 'http://callbackurl.com', 0, '', '2025-07-06 18:13:49', 1, 'Ad3232432562', '', 2, '', '2025-07-06 18:13:49');
INSERT INTO `orders` VALUES (572, 10000, 654232, '2025071100221583128517', 10.00, 0.00, '2025-07-11 00:22:15', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432563', '', 3, '', '2025-07-11 00:35:24');
INSERT INTO `orders` VALUES (573, 10000, 654232, '2025071100263954436500', 10.00, 0.00, '2025-07-11 00:26:39', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432564', '', 3, '', '2025-07-11 00:38:11');
INSERT INTO `orders` VALUES (575, 10000, 654232, '2025071100272796502378', 10.00, 0.00, '2025-07-11 00:27:27', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432565', '', 3, '', '2025-07-11 00:38:11');
INSERT INTO `orders` VALUES (576, 10000, 654232, '2025071100275148346412', 10.00, 0.00, '2025-07-11 00:27:51', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432566', '', 3, '', '2025-07-11 00:38:11');
INSERT INTO `orders` VALUES (577, 10000, 654232, '2025071100304227518998', 10.00, 0.00, '2025-07-11 00:30:42', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432567', '', 3, '', '2025-07-11 00:41:18');
INSERT INTO `orders` VALUES (578, 10000, 654232, '2025071100311285736005', 10.00, 0.00, '2025-07-11 00:31:12', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432568', '', 3, '', '2025-07-11 00:41:18');
INSERT INTO `orders` VALUES (579, 10000, 654232, '2025071100313442309481', 10.00, 0.00, '2025-07-11 00:31:34', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432569', '', 3, '', '2025-07-11 00:43:46');
INSERT INTO `orders` VALUES (580, 10000, 654232, '2025071100320286727736', 10.00, 0.00, '2025-07-11 00:32:02', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432570', '', 3, '', '2025-07-11 00:43:46');
INSERT INTO `orders` VALUES (583, 10000, 654232, '2025071100342761701306', 10.00, 0.00, '2025-07-11 00:34:27', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432573', '', 3, '', '2025-07-11 06:52:07');
INSERT INTO `orders` VALUES (588, 10000, 654232, '2025071100413394096009', 10.00, 0.00, '2025-07-11 00:41:33', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432578', '', 3, '', '2025-07-11 06:52:07');
INSERT INTO `orders` VALUES (589, 10000, 654232, '2025071106540185286513', 10.00, 0.00, '2025-07-11 06:54:01', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432579', '', 3, '', '2025-07-14 23:00:11');
INSERT INTO `orders` VALUES (590, 10000, 654232, '2025071107032921177329', 10.00, 0.00, '2025-07-11 07:03:29', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432580', '', 3, '', '2025-07-14 23:00:11');
INSERT INTO `orders` VALUES (591, 10000, 654232, '202507142302316168707', 10.00, 0.00, '2025-07-14 23:02:31', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432581', '', 3, '', '2025-07-14 23:17:30');
INSERT INTO `orders` VALUES (592, 10000, 654232, '2025071423032566070282', 2.00, 0.00, '2025-07-14 23:03:25', 'http://callbackurl.com', 0, '', '0000-00-00 00:00:00', 0, 'Ad3232432582', '', 3, '', '2025-07-14 23:17:30');

-- ----------------------------
-- Table structure for pay_channels
-- ----------------------------
DROP TABLE IF EXISTS `pay_channels`;
CREATE TABLE `pay_channels`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of pay_channels
-- ----------------------------
INSERT INTO `pay_channels` VALUES (1, 'Tron');
INSERT INTO `pay_channels` VALUES (2, 'Bsc');

-- ----------------------------
-- Table structure for transactions
-- ----------------------------
DROP TABLE IF EXISTS `transactions`;
CREATE TABLE `transactions`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_id` int(11) NULL DEFAULT NULL COMMENT '订单ID',
  `address_id` int(11) NULL DEFAULT NULL COMMENT '地址ID',
  `network` int(11) NOT NULL DEFAULT 0 COMMENT '网络：1:波场；2:币安',
  `tx_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '交易hash',
  `from_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '转账地址',
  `amount` decimal(10, 2) NULL DEFAULT NULL COMMENT '交易金额',
  `block_time` datetime NULL DEFAULT NULL COMMENT '交易时间',
  `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `callback_state` int(11) NOT NULL DEFAULT 0 COMMENT '通知状态：0:失败；1:成功；',
  `callback_err` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '通知报错信息',
  `callback_date` datetime NULL DEFAULT NULL COMMENT '最后通知时间',
  `callback_times` int(11) NOT NULL DEFAULT 0 COMMENT '通知次数',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `network_haxh_unique`(`network`, `tx_hash`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of transactions
-- ----------------------------
INSERT INTO `transactions` VALUES (1, 556, 1, 1, 'e9e47c0e74a0cdc25319e2ac10f91d1d58f841cecdec4cdf91c893ad0ff43018', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 100.00, '2025-07-06 12:56:39', '2025-07-06 13:01:38', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (2, 557, 2, 1, 'c80c90561583f1423a1ab1f10458f0c06247d7fd9a61a04b749562d33b73754a', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 10.00, '2025-07-06 13:03:51', '2025-07-06 13:04:19', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (3, 558, 1, 1, 'd137d07e302146920871243b18b5f20cf88e6b90f4935993777824e5765a573e', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 5.00, '2025-07-06 13:29:57', '2025-07-06 13:30:06', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (4, 558, 1, 1, 'b93ccb29740b14bb73ab527ecca369afe39bd8b1f4b660c4c80eb6f206f4c860', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 5.00, '2025-07-06 13:31:12', '2025-07-06 13:31:21', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (5, 559, 2, 1, '6c5ea7baf9271e351e76f01af6dc11d68fdbb90b4e5ed756fd525f4cfb1f1b20', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 10.00, '2025-07-06 13:34:30', '2025-07-06 13:34:48', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (6, 561, 3, 1, '0b51e529a03d4da10781d5f03d59098d4e41605634f5bb60e62f3dcae7a3c3d6', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 10.00, '2025-07-06 13:41:27', '2025-07-06 13:41:34', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (7, 562, 1, 1, 'a3a0eb6a510ee4c97f622ce7bd56bc2ddb384b46dde8eacf7ba76e1669d2fa80', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 10.00, '2025-07-06 13:48:30', '2025-07-06 13:48:39', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (8, 563, 2, 1, 'eda6061ad2c52280f0e2ba70ffd02005bc4a14ed7400f85607ae07f7a0dbe8fa', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 10.00, '2025-07-06 13:51:06', '2025-07-06 13:51:20', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (9, 564, 3, 1, 'f8b3d215c16efd59870121bf46975b165f33482bd1c870fbf79572f54cd150e4', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 10.00, '2025-07-06 13:54:15', '2025-07-06 13:54:29', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (10, 566, 1, 1, 'c641d825d190afc7ec2f7009f4c778ce00cc48971e82797d4da8601ddee288b9', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 10.00, '2025-07-06 13:59:39', '2025-07-06 13:59:46', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (11, 569, 1, 1, '32436630f7df7dd2a81a79b36cabcb22b95ea86d468ca56b940845a78bc9826d', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 10.00, '2025-07-06 14:05:30', '2025-07-06 14:05:45', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (12, 570, 1, 1, '7071082df2c9a9384ba0237b4654a21f6ada6e735aad5f35a1f06c9e6a16d024', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 10.00, '2025-07-06 18:03:48', '2025-07-06 18:04:05', 0, '', '0000-00-00 00:00:00', 0);
INSERT INTO `transactions` VALUES (13, 571, 2, 1, '5d4f666eb86cce9ff1156ead4f1d3454554cf2560d80535fece818a67aec5b75', 'TEckQwtjYS1tgVQrRzbTgNWRpB2y8cDW8s', 10.00, '2025-07-06 18:13:33', '2025-07-06 18:13:49', 0, '', '0000-00-00 00:00:00', 0);

SET FOREIGN_KEY_CHECKS = 1;
