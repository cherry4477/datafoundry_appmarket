


oc new-backingserviceinstance datahubstars-mysql --backingservice_name=Mysql --planid=56660431-6032-43D0-A114-FFA3BF521B66

oc secrets new-basicauth liuxu-dataos --username=user --password=password

oc new-app --name service-star https://code.dataos.io/asiainfoLDP/datahub_stars.git#develop \
    -e  SERVICE_NAME="datahub_stars" \
    \
    -e  CLOUD_PLATFORM="dataos" \
    \
    -e  ENV_NAME_KAFKA_ADDR="BSI_DATAHUBKAFKA_HOST" \
    -e  ENV_NAME_KAFKA_PORT="BSI_DATAHUBKAFKA_PORT" \
    \
    -e  ENV_NAME_MYSQL_ADDR="BSI_DATAHUBSTARSMYSQL_HOST" \
    -e  ENV_NAME_MYSQL_PORT="BSI_DATAHUBSTARSMYSQL_PORT" \
    -e  ENV_NAME_MYSQL_DATABASE="BSI_DATAHUBSTARSMYSQL_NAME" \
    -e  ENV_NAME_MYSQL_USER="BSI_DATAHUBSTARSMYSQL_USERNAME" \
    -e  ENV_NAME_MYSQL_PASSWORD="BSI_DATAHUBSTARSMYSQL_PASSWORD" \
    \
    -e  API_SERVER="service-discovery-datahub-develop.app.dataos.io" \
    -e  API_PORT="80" \
    \
    -e  REPOSIROTY_SERVICE_API_SERVER="service-repository" \
    -e  REPOSIROTY_SERVICE_API_PORT="8080" \
    -e  SUBSCRIPTION_SERVICE_API_SERVER="service-subscription" \
    -e  SUBSCRIPTION_SERVICE_API_PORT="8081" \
    -e  USER_SERVICE_API_SERVER="service-user" \
    -e  USER_SERVICE_API_PORT="80" \
    -e  BILL_SERVICE_API_SERVER="service-bill" \
    -e  BILL_SERVICE_API_PORT="80" \
    -e  DEAMON_SERVICE_API_SERVER="service-deamon" \
    -e  DEAMON_SERVICE_API_PORT="80" \
    \
    -e  MYSQL_CONFIG_DONT_UPGRADE_TABLES="false" \
    -e  LOG_LEVEL="debug" \
    -e  IS_STAGE_OR_DEV="false" \
    -e  TEST_OPSTAT_SENDING="no"

oc bind datahubstars-mysql service-star
oc bind datahub-kafka service-star

oc edit bc service-star

    source:
        git:
            ref: develop
            uri: https://code.dataos.io/asiainfoLDP/datahub_stars.git
        //>>>
    sourceSecret:
      name: liuxu-dataos
        //<<<
        secrets: []
        type: Git

oc start-build service-star
