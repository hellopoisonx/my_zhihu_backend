---
title: 我的知乎后端
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# 我的知乎后端

## 部署
采用docker部署
### 依赖
- golang 1.26.0
- docker
- docker-buildx
- docker-compose
### 相关端口
- mysql:3306
- redis:6379
- app:8080
```shell
$ docker-compose build --no-cache # 构建
$ docker-compose up -d # 运行
$ docker-compose logs -f app # 运行日志
```

## 语言相关
基于go 1.26.0版本 使用泛型简化代码 使用option模式编写更新接口

## 日志系统
采用zap作为日志系统来实现结构化日志 借助golang的构建约束来区分生产和开发环境

## 配置管理
采用viper作为配置管理 可以从环境变量、yaml文件和默认值中读取配置 支持热更新

## api设计
借助 apifox 来实现接口的设计、测试、数据mock

## 用户权限设计
采用双token方案(refreshToken + accessToken) accessToken采用短时效jwt以实现无状态凭证存储减轻服务器压力 同时采用长时效有状态refreshToken+redis以实现用户状态的无感刷新、单点登录和服务器主动控制用户上下线

## 用户存储相关
- 采用雪花算法来生成用户全局唯一id 保证后续拆表等操作的便利性
- 用户密码先使用sha512统一长度后采用bcrypt算法来加盐加密存储 
- 客户端登录凭证为邮箱加明文密码 后续认证凭证为accessToken
- 用户的简介、头像和个人设置采用mysql json结构存储 保证后续拓展性 同时减少联表查询操作

## 限流
使用 [time/rate](https://pkg.go.dev/golang.org/x/time/rate) 包提供的令牌桶实现了基于主机地址的限流

## TODO
- [ ] 回答、评论 相关功能
- [ ] 管理员控制
- [ ] 缓存
- [ ] 移除部分硬编码数据

## 接口响应结构
统一返回 200 http 状态码 采用业务码判断请求结果状态
### 通用错误码

| 错误码 | 常量名              | 描述   |
|-----|------------------|------|
| 0   | `ErrCodeOk`      | 成功   |
| 1   | `ErrCodeUnknown` | 未知错误 |
### 用户相关错误码 (10001-10014)

| 错误码   | 常量名                                 | 描述     | 对应错误变量                    |
|-------|-------------------------------------|--------|---------------------------|
| 10001 | `ErrCodeUserNotExists`              | 用户不存在  | `ErrUserNotExists`        |
| 10002 | `ErrCodeUserAlreadyExists`          | 用户已存在  | `ErrUserAlreadyExists`    |
| 10003 | `ErrCodeInvalidParameters`          | 参数无效   | -                         |
| 10004 | `ErrCodeInvalidAuthorizationHeader` | 认证头无效  | -                         |
| 10005 | `ErrCodeTimeout`                    | 超时     | `ErrTimeout`              |
| 10006 | `ErrCodeUserPermissionDenied`       | 用户权限不足 | `ErrUserPermissionDenied` |
| 10007 | `ErrCodeUserNotAuthorized`          | 用户未授权  | `ErrUserNotAuthorized`    |
| 10008 | `ErrCodeUserWrongPassword`          | 密码错误   | `ErrUserWrongPassword`    |
| 10009 | `ErrCodeUserInvalidToken`           | 令牌无效   | `ErrUserInvalidToken`     |
| 10010 | `ErrCodeUserWrongTokenType`         | 令牌类型错误 | `ErrUserWrongTokenType`   |
| 10011 | `ErrCodeQuestionNotFound`           | 问题未找到  | `ErrQuestionNotFound`     |
| 10012 | `ErrCodeAnswerNotFound`             | 回答未找到  | `ErrAnswerNotFound`       |
| 10013 | `ErrCodeCommentNotFound`            | 评论未找到  | `ErrCommentNotFound`      |
| 10014 | `ErrCodeTooManyRequest`             | 请求频繁    | `ErrTooManyRequest`       |
### 系统相关错误码 (20001-20004)

| 错误码   | 常量名                 | 描述          |
|-------|---------------------|-------------|
| 20001 | `ErrCodeMysql`      | MySQL 数据库错误 |
| 20002 | `ErrCodeRedis`      | Redis 错误    |
| 20003 | `ErrCodeUserToken`  | 用户令牌错误      |
| 20004 | `ErrCodeEncryption` | 加密错误        |

Base URLs:

* <a href="http://localhost:8080">开发环境: http://localhost:8080</a>

# Authentication

- HTTP Authentication, scheme: bearer

# users

## POST 创建用户

POST /users

> Body 请求参数

```json
{
  "username": "杞婷婷",
  "email": "p1pvgw1@yahoo.com.cn",
  "password": "3VmAKjMINUM55N9",
  "gender": 1,
  "region": "华中",
  "settings": {
    "hide_privacy": true
  },
  "other": {
    "introduction": "飞行员发烧友",
    "icon": "https://avatars.githubusercontent.com/u/39475011"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 是 |none|
|» username|body|string| 是 |none|
|» email|body|string| 是 |none|
|» password|body|string| 是 |none|
|» gender|body|integer(int32)| 否 |none|
|» region|body|string| 否 |none|
|» settings|body|[UserSettings](#schemausersettings)| 否 |none|
|»» hide_privacy|body|boolean| 否 |none|
|» other|body|[UserOtherInfo](#schemauserotherinfo)| 否 |none|
|»» introduction|body|string| 否 |none|
|»» icon|body|string| 否 |none|

#### 枚举值

|属性|值|
|---|---|
|» gender|0|
|» gender|1|
|» gender|2|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": {
    "id": 0,
    "username": "string",
    "email": "string",
    "gender": 0,
    "region": "string",
    "other": {
      "introduction": "string",
      "icon": "string"
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» ok|boolean|true|none||none|
|» internal_error|boolean|true|none||none|
|» message|string|true|none||none|
|» body|object|true|none||none|
|»» id|number|true|none||none|
|»» username|string|true|none||none|
|»» email|string|true|none||none|
|»» gender|integer(int32)|true|none||none|
|»» region|string|true|none||none|
|»» other|object|true|none||none|
|»»» introduction|string|true|none||none|
|»»» icon|string|true|none||none|

#### 枚举值

|属性|值|
|---|---|
|gender|0|
|gender|1|
|gender|2|

## GET 搜索用户

GET /users

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|username|query|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": [
    0
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» ok|boolean|true|none||none|
|» internal_error|boolean|true|none||none|
|» message|string|true|none||none|
|» body|[integer]|true|none|array of UserId|none|

## DELETE 删除用户

DELETE /users/me

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|[Response](#schemaresponse)|

## PATCH 更新用户信息

PATCH /users/me

> Body 请求参数

```json
{
  "username": "string",
  "password": "string",
  "gender": 0,
  "region": "string",
  "settings": {
    "hide_privacy": true
  },
  "other": {
    "introduction": "string",
    "icon": "string"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 是 |none|
|» username|body|string| 否 |none|
|» password|body|string| 否 |none|
|» gender|body|integer(int32)| 否 |none|
|» region|body|string| 否 |none|
|» settings|body|[UserSettings](#schemausersettings)| 否 |none|
|»» hide_privacy|body|boolean| 否 |none|
|» other|body|[UserOtherInfo](#schemauserotherinfo)| 否 |none|
|»» introduction|body|string| 否 |none|
|»» icon|body|string| 否 |none|

#### 枚举值

|属性|值|
|---|---|
|» gender|0|
|» gender|1|
|» gender|2|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": {
    "id": 0,
    "username": "string",
    "email": "string",
    "gender": 0,
    "region": "string",
    "other": {
      "introduction": "string",
      "icon": "string"
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» ok|boolean|true|none||none|
|» internal_error|boolean|true|none||none|
|» message|string|true|none||none|
|» body|object|true|none||none|
|»» id|number|true|none||none|
|»» username|string|true|none||none|
|»» email|string|true|none||none|
|»» gender|integer(int32)|true|none||none|
|»» region|string|true|none||none|
|»» other|object|true|none||none|
|»»» introduction|string|true|none||none|
|»»» icon|string|true|none||none|

#### 枚举值

|属性|值|
|---|---|
|gender|0|
|gender|1|
|gender|2|

## GET 获取用户信息

GET /users/{id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|path|integer| 是 |UserId|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": {
    "id": 0,
    "username": "string",
    "email": "string",
    "gender": 0,
    "region": "string",
    "other": {
      "introduction": "string",
      "icon": "string"
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» ok|boolean|true|none||none|
|» internal_error|boolean|true|none||none|
|» message|string|true|none||none|
|» body|object|true|none||none|
|»» id|number|true|none||none|
|»» username|string|true|none||none|
|»» email|string|true|none||none|
|»» gender|integer(int32)|true|none||none|
|»» region|string|true|none||none|
|»» other|object|true|none||none|
|»»» introduction|string|true|none||none|
|»»» icon|string|true|none||none|

#### 枚举值

|属性|值|
|---|---|
|gender|0|
|gender|1|
|gender|2|

## POST 关注用户

POST /users/follow/{id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|path|integer| 是 |UserId of the following user|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|[Response](#schemaresponse)|

## DELETE 取关用户

DELETE /users/follow/{id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|path|integer| 是 |UserId of the following user|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|[Response](#schemaresponse)|

## GET 获取粉丝列表

GET /users/followers/{id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|path|integer| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": [
    0
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» ok|boolean|true|none||none|
|» internal_error|boolean|true|none||none|
|» message|string|true|none||none|
|» body|[integer]|true|none||none|

## GET 获取关注列表

GET /users/followings/{id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|path|integer| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": [
    0
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» ok|boolean|true|none||none|
|» internal_error|boolean|true|none||none|
|» message|string|true|none||none|
|» body|[integer]|true|none||none|

# auth

## POST 登录

POST /auth

> Body 请求参数

```json
{
  "id": 0,
  "email": "string",
  "password": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 是 |none|
|» id|body|number| 否 |none|
|» email|body|string| 是 |none|
|» password|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": {
    "access_token": {
      "token": "string",
      "expire_at": "string"
    },
    "refresh_token": {
      "token": "string",
      "expire_at": "string"
    }
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» ok|boolean|true|none||none|
|» internal_error|boolean|true|none||none|
|» message|string|true|none||none|
|» body|[AuthResponse](#schemaauthresponse)|true|none||none|
|»» access_token|[Token](#schematoken)|true|none||none|
|»»» token|string|true|none||none|
|»»» expire_at|string|true|none||none|
|»» refresh_token|[Token](#schematoken)|true|none||none|
|»»» token|string|true|none||none|
|»»» expire_at|string|true|none||none|

## PATCH 续期

PATCH /auth

> Body 请求参数

```json
{
  "user_id": 0,
  "refresh_token": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 是 |none|
|» user_id|body|integer| 是 |none|
|» refresh_token|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": {
    "token": "string",
    "expire_at": "string"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» ok|boolean|true|none||none|
|» internal_error|boolean|true|none||none|
|» message|string|true|none||none|
|» body|[Token](#schematoken)|true|none||none|
|»» token|string|true|none||none|
|»» expire_at|string|true|none||none|

## DELETE 登出

DELETE /auth

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|[Response](#schemaresponse)|

# article

## POST 创建问题

POST /questions

> Body 请求参数

```json
{
  "title": "string",
  "content": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 是 |none|
|» title|body|string| 是 |none|
|» content|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": {
    "id": 0,
    "title": "string",
    "content": "string",
    "author_id": 0,
    "is_available": true,
    "updated_at": "string"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» ok|boolean|true|none||none|
|» internal_error|boolean|true|none||none|
|» message|string|true|none||none|
|» body|[Question](#schemaquestion)|false|none||none|
|»» id|number|true|none||none|
|»» title|string|true|none||none|
|»» content|string|true|none||none|
|»» author_id|number|true|none||none|
|»» is_available|boolean|true|none||none|
|»» updated_at|string|true|none||none|

## GET 搜索

GET /questions

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|page|query|integer| 否 |none|
|size|query|integer| 否 |none|
|keywords|query|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": [
    0
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» ok|boolean|true|none||none|
|» internal_error|boolean|true|none||none|
|» message|string|true|none||none|
|» body|[number]|true|none||none|

## PATCH 更新问题

PATCH /questions/{id}

> Body 请求参数

```json
{
  "title": "string",
  "content": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|path|string| 是 |none|
|body|body|object| 是 |none|
|» title|body|string| 否 |none|
|» content|body|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": {
    "id": 0,
    "title": "string",
    "content": "string",
    "author_id": 0,
    "is_available": true,
    "updated_at": "string"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» ok|boolean|true|none||none|
|» internal_error|boolean|true|none||none|
|» message|string|true|none||none|
|» body|[Question](#schemaquestion)|false|none||none|
|»» id|number|true|none||none|
|»» title|string|true|none||none|
|»» content|string|true|none||none|
|»» author_id|number|true|none||none|
|»» is_available|boolean|true|none||none|
|»» updated_at|string|true|none||none|

## DELETE 删除问题

DELETE /questions/{id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|[Response](#schemaresponse)|

# 数据模型

<h2 id="tocS_User">User</h2>

<a id="schemauser"></a>
<a id="schema_User"></a>
<a id="tocSuser"></a>
<a id="tocsuser"></a>

```json
{
  "id": 0,
  "username": "string",
  "email": "string",
  "password": "string",
  "gender": 0,
  "region": "string",
  "settings": {
    "hide_privacy": true
  },
  "other": {
    "introduction": "string",
    "icon": "string"
  }
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|id|number|false|none||none|
|username|string|true|none||none|
|email|string|true|none||none|
|password|string|true|none||none|
|gender|integer(int32)|false|none||none|
|region|string|false|none||none|
|settings|[UserSettings](#schemausersettings)|false|none||none|
|other|[UserOtherInfo](#schemauserotherinfo)|false|none||none|

#### 枚举值

|属性|值|
|---|---|
|gender|0|
|gender|1|
|gender|2|

<h2 id="tocS_UserSettings">UserSettings</h2>

<a id="schemausersettings"></a>
<a id="schema_UserSettings"></a>
<a id="tocSusersettings"></a>
<a id="tocsusersettings"></a>

```json
{
  "hide_privacy": true
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|hide_privacy|boolean|false|none||none|

<h2 id="tocS_UserOtherInfo">UserOtherInfo</h2>

<a id="schemauserotherinfo"></a>
<a id="schema_UserOtherInfo"></a>
<a id="tocSuserotherinfo"></a>
<a id="tocsuserotherinfo"></a>

```json
{
  "introduction": "string",
  "icon": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|introduction|string|false|none||none|
|icon|string|false|none||none|

<h2 id="tocS_Response">Response</h2>

<a id="schemaresponse"></a>
<a id="schema_Response"></a>
<a id="tocSresponse"></a>
<a id="tocsresponse"></a>

```json
{
  "code": 0,
  "ok": true,
  "internal_error": true,
  "message": "string",
  "body": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|code|integer|true|none||none|
|ok|boolean|true|none||none|
|internal_error|boolean|true|none||none|
|message|string|true|none||none|
|body|any|false|none||none|

oneOf

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|string|false|none||none|

xor

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|integer|false|none||none|

xor

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|boolean|false|none||none|

xor

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|array|false|none||none|

xor

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|object|false|none||none|

xor

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|number|false|none||none|

<h2 id="tocS_UserId">UserId</h2>

<a id="schemauserid"></a>
<a id="schema_UserId"></a>
<a id="tocSuserid"></a>
<a id="tocsuserid"></a>

```json
0

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|*anonymous*|integer|false|none||none|

<h2 id="tocS_AuthResponse">AuthResponse</h2>

<a id="schemaauthresponse"></a>
<a id="schema_AuthResponse"></a>
<a id="tocSauthresponse"></a>
<a id="tocsauthresponse"></a>

```json
{
  "access_token": {
    "token": "string",
    "expire_at": "string"
  },
  "refresh_token": {
    "token": "string",
    "expire_at": "string"
  }
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|access_token|[Token](#schematoken)|true|none||none|
|refresh_token|[Token](#schematoken)|true|none||none|

<h2 id="tocS_Token">Token</h2>

<a id="schematoken"></a>
<a id="schema_Token"></a>
<a id="tocStoken"></a>
<a id="tocstoken"></a>

```json
{
  "token": "string",
  "expire_at": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|token|string|true|none||none|
|expire_at|string|true|none||none|

<h2 id="tocS_Question">Question</h2>

<a id="schemaquestion"></a>
<a id="schema_Question"></a>
<a id="tocSquestion"></a>
<a id="tocsquestion"></a>

```json
{
  "id": 0,
  "title": "string",
  "content": "string",
  "author_id": 0,
  "is_available": true,
  "updated_at": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|id|number|true|none||none|
|title|string|true|none||none|
|content|string|true|none||none|
|author_id|number|true|none||none|
|is_available|boolean|true|none||none|
|updated_at|string|true|none||none|

