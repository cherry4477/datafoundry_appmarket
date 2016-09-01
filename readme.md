
## 数据库设计

```
CREATE TABLE IF NOT EXISTS DF_SAAS_APP
(
   APP_ID             BIGINT NOT NULL AUTO_INCREMENT,
   PROVIDER           VARCHAR(120) NOT NULL,
   NAME               VARCHAR(255) NOT NULL,
   URL                VARCHAR(255) NOT NULL,
   VERSION            VARCHAR(64) NOT NULL,
   CATEGORY           VARCHAR(64),
   DESCRIPTION        TEXT,
   ICON_URL           VARCHAR(255) NOT NULL,
   CREATE_TIME        DATETIME,
   HOTNESS            INT NOT NULL,
   PRIMARY KEY (APP_ID)
)  DEFAULT CHARSET=UTF8;
```

## API设计

### POST /saas/v1/apps

创建一个SaaS App。

Body Parameters:
```
provider: 提供方
url: 应用地址
name: 应用名称
version: 当前应用版本
category: 应用类别
description: 应用描述
```

Return Result (json):
```
code: 返回码
msg: 返回信息
data.id: 应用id
```

### DELETE /saas/v1/apps/{id}

删除一个SaaS App。

Path Parameters:
```
id: 应用id
```

Return Result (json):
```
code: 返回码
msg: 返回信息
```

### PUT /saas/v1/apps/{id}

修改一个SaaS App。

Path Parameters:
```
id: 应用id
```

Body Parameters:
```
provider: 提供方
url: 应用地址
name: 应用名称
version: 当前应用版本
category: 应用类别
description: 应用描述
```

Return Result (json):
```
code: 返回码
msg: 返回信息
```

### GET /saas/v1/apps/{id}

查询一个SaaS App。

Path Parameters:
```
id: 应用id
```

Return Result (json):
```
code: 返回码
msg: 返回信息
data.id
data.provider
data.url
data.name
data.version
data.category
data.description
data.iconUrl
data.createTime
```

### GET /saas/v1/apps?category={category}&orderby={orderby}&provider={provider}

查询SaaS App列表。

Query Parameters:
```
category: app的类别。可选。如果忽略，表示所有类别。
orderby: 排序依据。可选。合法值包括hotness|createtime，默认为hotness。
provider: 提供方。可选。
```

Return Result (json):
```
code: 返回码
msg: 返回信息
data.total
data.results
data.results[0].id
data.results[0].provider
data.results[0].url
data.results[0].name
data.results[0].version
data.results[0].category
data.results[0].description
data.results[0].iconUrl
data.results[0].createTime
...
```

## 部署

```
oc new-instance MysqlForAppMarket --service=Mysql --plan=NoCase

oc new-app --name datafoundryappmarket https://github.com/asiainfoLDP/datafoundry_appmarket.git#develop \
    -e  CLOUD_PLATFORM="dataos" \
    \
    -e  DATAFOUNDRY_HOST_ADDR="xxx" \
    \
    -e  ENV_NAME_MYSQL_ADDR="BSI_MYSQL_MYSQLFORAPPMARKET_HOST" \
    -e  ENV_NAME_MYSQL_PORT="BSI_MYSQL_MYSQLFORAPPMARKET_PORT" \
    -e  ENV_NAME_MYSQL_DATABASE="BSI_MYSQL_MYSQLFORAPPMARKET_NAME" \
    -e  ENV_NAME_MYSQL_USER="BSI_MYSQL_MYSQLFORAPPMARKET_USERNAME" \
    -e  ENV_NAME_MYSQL_PASSWORD="BSI_MYSQL_MYSQLFORAPPMARKET_PASSWORD" \
    \
    -e  MYSQL_CONFIG_DONT_UPGRADE_TABLES="false" \
    -e  LOG_LEVEL="debug"

oc bind MysqlForAppMarket datafoundryappmarket

oc expose service datafoundryappmarket --hostname=datafoundry-appmarket.app.dataos.io

oc start-build datafoundryappmarket

```
