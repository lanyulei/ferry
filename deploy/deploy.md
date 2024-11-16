# ferry的Kubernetes部署

## 做出的更改

 通过项目目录下的`Dockerfile`将ferry打包为镜像，推送至dockerhub

镜像名：`beatrueman/ferry:1.0.0`

在项目目录下新增了`deploy`目录，用于ferry的K8s平台部署，其中包含了`helm`和`kubernetes`两个目录

- helm目录：包含了一个ferry的chart
- kubernetes目录：包含部署ferry的资源文件

## kubernetes目录

包含以下文件：

- `config.yaml`：ferry的`ConfigMap`，包含了ferry的`rbac_model.conf`和settings.yml
- `secret.yaml`：用于保存数据库凭据
- `deploy.yaml`：包含ferry主平台的`deployment`，`service`和`pvc`，pvc用于持久化`/opt/workflow/ferry/config`
- `mysql.yaml`：包含用于ferry的mysql数据库的`statefulset`，`service`和`pvc`
- `redis.yaml`：包含用于ferry的redis的`deployment`，`service`和`pvc`
- `sql目录`：保存了`ferry.sql`和`db.sql`，需要用户手动的导入到数据库中

使用如下命令部署：

```
kubectl apply -f <sources>.yaml
```

## helm目录

目录结构如下

```
.
|-- Chart.yaml
|-- charts
|-- templates
|   |-- NOTES.txt
|   |-- _helpers.tpl
|   |-- configmap.yaml # ferry主平台
|   |-- deployment.yaml # ferry主平台
|   |-- mysql # ferry依赖的mysql的资源模板文件
|   |   |-- persistentvolumeclaim.yaml
|   |   |-- service.yaml
|   |   `-- statefulset.yaml
|   |-- persistentvolumeclaim.yaml
|   |-- redis # ferry依赖的redis的资源模板文件
|   |   |-- deployment.yaml
|   |   |-- persistentvolumeclaim.yaml
|   |   `-- service.yaml
|   |-- secret.yaml # ferry主平台
|   `-- service.yaml # ferry主平台
`-- values.yaml # helm配置文件

```

***values.yaml介绍***

```
replicaCount: 1 # ferry、mysql、redis副本数 
namespace: ferry # ferry的命名空间

global:
  storageClassName: longhorn # 用户可以指定存储类

# 数据库凭据，主要用于secret
env:
  ENV: "production"
  MYSQL_ROOT_PASSWORD: "123456"
  MYSQL_USER: "ferry"
  MYSQL_DATABASE: "ferry"
  MYSQL_PASSWORD: "123456"

# ferry的配置项
ferry:
  image:
    repository: beatrueman/ferry
    tag: "1.0.0"
    pullPolicy: IfNotPresent

  service:
    type: NodePort
    port: 8002

  # ferry的持久卷
  persistentVolume:
    accessModes:
      - ReadWriteOnce
    size: 2Gi

# 如果要自用mysql，请将enable设置为false
# 并且需要修改下方configMap.settings.yml中的database.host
mysql:
  enable: false
  image:
    repository: mysql
    tag: 8.4.0-oraclelinux8
  port: 3306
  persistentVolume:
    accessModes:
      - ReadWriteOnce
    size: 2Gi

# 如果要自用redis，请将enable设置为false
# 并且需要修改下方configMap.settings.yml中的redis.url
redis:
  enable: false
  image:
    repository: redis
    tag: 7.0.5-alpine
  port: 6379
  persistentVolume:
    accessModes:
      - ReadWriteOnce
    size: 2Gi

# ferry的配置文件
# 主要关注database和redis
configMap:
  rbac_model_conf: |
    [request_definition]
    r = sub, obj, act

    [policy_definition]
    p = sub, obj, act

    [policy_effect]
    e = some(where (p.eft == allow))

    [matchers]
    m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
  settings_yml: |
    script:
      path: ./static/scripts
    settings:
      application:
        domain: localhost:8002
        host: 0.0.0.0
        ishttps: false
        mode: dev
        name: ferry
        port: "8002"
        readtimeout: 1
        writertimeout: 2
      database:
        dbtype: mysql
        host: ferry-mysql.ferry.svc.cluster.local # 这里使用K8s部署mysql service的DNS，如果使用自用的数据库，请更改
        name: ferry
        password: 123456
        port: 3306
        username: root
      domain:
        gethost: 1
        url: localhost:9527
      email:
        alias: ferry
        host: smtp.163.com
        pass: your password
        port: 465
        user: fdevops@163.com
      gorm:
        logmode: 0
        maxidleconn: 0
        maxopenconn: 20000
      jwt:
        secret: ferry
        timeout: 86400
      ldap:
        anonymousquery: 0
        basedn: dc=fdevops,dc=com
        bindpwd: 123456
        binduserdn: cn=admin,dc=fdevops,dc=com
        host: localhost
        port: 389
        tls: 0
        userfield: uid
      log:
        compress: 1
        consolestdout: 1
        filestdout: 0
        level: debug
        localtime: 1
        maxage: 30
        maxbackups: 300
        maxsize: 10240
        path: ./logs/ferry.log
      public:
        islocation: 0
      redis:
        url: redis://ferry-redis.ferry.svc.cluster.local:6379  # 这里使用K8s部署redis service的DNS，如果使用自用的redis，请更改
      ssl:
        key: keystring
        pem: temp/pem.pem 
        runAsUser: 1000

```

使用如下命令部署：

```
helm install -n <namespace> <release> .
# 建议在ferry命名空间下部署
```

ferry依赖于mysql

如果使用了附带的mysql，当`helm install`后，需要等待mysql容器准备好后，ferry容器才可以正常运行，期间如果ferry没有正常运行，只需要在mysql正常启动后，重启ferry容器即可（delete它）

一切就绪后，注意要把`templates/mysql/sql`下的两个sql文件（`ferry.sql`和`db.sql`）导入名为ferry的数据库，先导入`ferry.sql`，后导入`db.sql`

## 环境介绍

### 集群环境

![image-20241117011316070](https://gitee.com/beatrueman/images/raw/master/img/202411170113207.png)

### helm版本

![image-20241117011418385](https://gitee.com/beatrueman/images/raw/master/img/202411170114447.png)

## 部署成功证明

helm部署

![image-20241117002900813](https://gitee.com/beatrueman/images/raw/master/img/202411170029915.png)

资源文件部署：

使用了自用的数据库

![image-20241117011709830](https://gitee.com/beatrueman/images/raw/master/img/202411170117927.png)