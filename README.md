## 日历提醒服务

### 系统有MySQL两张表

```sql
# 用户表
(
    id         INT PRIMARY KEY AUTO_INCREMENT COMMENT '唯一标识用户的ID',
    mobile     VARCHAR(20)  NOT NULL UNIQUE COMMENT '用户手机号，国内格式',
    creator_id VARCHAR(128) NOT NULL UNIQUE COMMENT '用户唯一标识，用于关联提醒信息',
    created_at DATETIME NOT NULL COMMENT '用户创建时间',
    updated_at DATETIME NOT NULL COMMENT '用户信息最后更新时间'
)
```

```sql
# 日历提醒表
(
    id         INT PRIMARY KEY AUTO_INCREMENT COMMENT '唯一标识提醒信息的ID',
    creator_id VARCHAR(128) NOT NULL COMMENT '提醒信息创建者的ID',
    content    TEXT         NOT NULL COMMENT '提醒内容',
    remind_at  DATETIME     NOT NULL COMMENT '提醒的具体时间', 
    created_at DATETIME     NOT NULL  COMMENT '提醒信息创建时间',
    updated_at DATETIME     NOT NULL  COMMENT '提醒信息最后更新时间'
)
```

### 通过登录注册来实现每个用户只能管理本人的提醒信息

redis来储存token和手机短信信息实现登录/注册功能；日历提醒的CRUD功能；rabbitmq实现延迟消息队列推送手机短信通知提醒。
系统通过docker容器部署在阿里云服务器上，可以通过接口文档 进行访问。
**切记：REST API风格访问**

