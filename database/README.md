# Database 启动指南
首先在系统上启动 mysql 与 redis

为 mysql 创建用户，数据库，并赋予用户操作数据库的权限：
```sql
CREATE USER 'test'@'localhost' IDENTIFIED BY '123456';
CREATE DATABASE app;
grant all privileges on *.* to test@'localhost' identified by "123456" with grant option;
```
用户，数据库，密码的设定在`mysql.go`中。

进入数据库，创建表单。
```sql
use app;
CREATE TABLE rules(
    id int UNSIGNED AUTO_INCREMENT,
    aid INT UNSIGNED,
    hit_count INT UNSIGNED DEFAULT 0,
    download_count INT UNSIGNED DEFAULT 0,  
    platform CHAR(16),
    download_url VARCHAR(128),
    update_version_code	VARCHAR(128),
    device_list TEXT,
    md5	VARCHAR(128),
    max_update_version_code	VARCHAR(128),
    min_update_version_code	VARCHAR(128),
    max_os_api	TINYINT UNSIGNED,
    min_os_api	TINYINT UNSIGNED,
    cpu_arch	TINYINT UNSIGNED,
    channel	VARCHAR(128),
    title	VARCHAR(256),
    update_tips	VARCHAR(1024),
    enabled	BOOLEAN DEFAULT true,
	create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ( id )
);
```

启动 redis，默认即可。
```bash
redis-server
```