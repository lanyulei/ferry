# 安装

> 需注意因使用到了json类型的字段，因此MySQL需是5.7以上的版本。
> 
> MySQL > 5.7
> 
> Go >= 1.14
> 
> Redis

若是安装出错，请先确认redis及MySQL是否安装配置成功，若是还有问题，可在群内提问。

## 配置文件介绍

    script:
      path: ./static/scripts # 任务脚本路径
    settings:
      application:
        domain: localhost:8002 # 用于将HTTP请求重定向到HTTPS的主机名
        host: 0.0.0.0 # 启动的地址，主机ip 或者域名，默认0.0.0.0
        ishttps: false # 是否为HTTPS
        mode: dev # 开发模式
        name: ferry-test # 服务名称
        port: "8002" # 启动端口
        readtimeout: 1 # 请求读取超时时间，从连接被接受(accept)到request body完全被读取(如果你不读取body，那么时间截止到读完header为止)
        writertimeout: 2 # 从request header的读取结束开始，到response write结束为止(也就是ServeHTTP 方法的声明周期)
      database:
        dbtype: mysql # 数据库类型
        host: 127.0.0.1 # 数据库地址
        name: ferry # 数据库名称
        password: 123456 # 数据库密码
        port: 3306 # 数据库端口
        username: ferry # 数据库用户名
      email:
        alias: ferry # 邮箱别名
        host: smtp.163.com # 邮件服务器
        pass: your password # 邮箱密码
        port: 465 # 邮件服务器端口
        user: fdevops@163.com # 邮箱账号
      gorm:
        logmode: 0 # gorm详细日志输出，0表示不输出，1表示输出
        maxidleconn: 0 # 最大空闲连接
        maxopenconn: 20000 # 最大连接数据
      jwt:
        secret: ferry # JWT加密字符串
        timeout: 3600 # 过期时间单位：秒
      log:
        dir: logs # 日志路径
        operdb: false
      ssl:
        key: keystring
        pem: temp/pem.pem

## 本地开发

后端程序启动：

    # 1\. 拉取代码，以下命令二选一即可：
    git clone https://github.com/lanyulei/ferry.git
    git clone https://gitee.com/yllan/ferry.git

    # 2\. 进入工作路径
    cd ferry

    # 3\. 修改配置
    vim config/settings.dev.yml
      1). 修改为自己的数据库信息
      2). 修改为自己的邮件服务器地址
    其他的根据情况来修改调整

    # 4\. 安装依赖
    go get

    # 5\. 连接数据库，并创建数据库
    create database ferry charset 'utf8mb4';

    # 6\. 初始化数据结构
    go run main.go init -c=config/settings.dev.yml

    # 7\. 测试启动程序，没有报错及没有问题
    go run main.go server -c=config/settings.dev.yml

    # 8\. 热加载方式启动
    air

前端程序启动：

    # 1\. 拉取代码，以下命令二选一即可：
    git clone https://github.com/lanyulei/ferry_web.git
    git clone https://gitee.com/yllan/ferry_web.git

    # 2\. 进入工作路径
    cd ferry_web

    # 3\. 安装依赖
    npm config set registry https://registry.npm.taobao.org
    npm install
    # 若npm install安装失败，可尝试使用一下命令安装
    npm install --unsafe-perm

    # 推荐使用cnpm
    npm install -g cnpm --registry=https://registry.npm.taobao.org
    cnpm install

    # 4\. 启动程序
    npm run dev

    # 5\. 访问http://localhost:9527，是否可正常访问

## 部署线上

后端部署：

    # 1\. 拉取代码，以下命令二选一即可：
    git clone https://github.com/lanyulei/ferry.git
    git clone https://gitee.com/yllan/ferry.git

    # 2\. 进入工作路径
    cd ferry

    # 3\. 交叉编译（centos）
    env GOOS=linux GOARCH=amd64 go build
    更多交叉编译内容，请访问 https://www.fdevops.com/2020/03/08/go-locale-configuration

    # 4\. config目录上传到项目根路径下，并确认配置信息是否正确
    vim config/settings.yml
      1). 修改为自己的数据库信息
      2). 修改为自己的邮件服务器地址
    其他的根据情况来修改调整

    # 4\. 创建日志路径及静态文件经历
    mkdir -p log static/uploadfile static/scripts static/template

    # 5\. 将本地项目下static/template目录下的所有文件上传的到，服务器对应的项目目录下static/template

    # 6\. 连接数据库，并创建数据库
    create database ferry charset 'utf8mb4';

    # 7\. 初始化数据
    ./ferry init -c=config/settings.yml

    # 8\. 启动程序，推荐通过"进程管理工具"进行启动维护
    nohup ./ferry server -c=config/settings.yml > /dev/null 2>&1 &

前端部署：

    # 1\. 拉取代码，以下命令二选一即可：
    git clone https://github.com/lanyulei/ferry_web.git
    git clone https://gitee.com/yllan/ferry_web.git

    # 2\. 进入工作路径
    cd ferry_web

    # 3\. 安装依赖
    npm config set registry https://registry.npm.taobao.org
    npm install
    # 若npm install安装失败，可尝试使用一下命令安装
    npm install --unsafe-perm

    # 推荐使用cnpm
    npm install -g cnpm --registry=https://registry.npm.taobao.org
    cnpm install

    # 4\. 修改 .env.production 文件
    # base api
    VUE_APP_BASE_API = 'http://fdevops.com:8001'  # 修改为您自己的域名

    # 5\. 编译
    npm run build:prod

    # 6\. 将dist目录上传至项目路径下即可。
    mv dist web

    # 7\. nginx配置，根据业务自行调整即可
      server {
        listen 8001; # 监听端口
        server_name localhost; # 域名可以有多个，用空格隔开

        #charset koi8-r;

        #access_log  logs/host.access.log  main;
        location / {
          root /data/ferry/web;
          index index.html index.htm; #目录内的默认打开文件,如果没有匹配到index.html,则搜索index.htm,依次类推
        }

        #ssl配置省略
        location /api {
          # rewrite ^.+api/?(.*)$ /$1 break;
          proxy_pass http://127.0.0.1:8002; #node api server 即需要代理的IP地址
          proxy_redirect off;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
          proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }

        # 登陆
        location /login {
          proxy_pass http://127.0.0.1:8002; #node api server 即需要代理的IP地址
        }

        # 刷新token
        location /refresh_token {
          proxy_pass http://127.0.0.1:8002; #node api server 即需要代理的IP地址
        }

        # 接口地址
        location /swagger {
          proxy_pass http://127.0.0.1:8002; #node api server 即需要代理的IP地址
        }

        # 后端静态文件路径
        location /static/uploadfile {
          proxy_pass http://127.0.0.1:8002; #node api server 即需要代理的IP地址
        }

        #error_page  404              /404.html;    #对错误页面404.html 做了定向配置

        # redirect server error pages to the static page /50x.html
        #将服务器错误页面重定向到静态页面/50x.html
        #
        error_page 500 502 503 504 /50x.html;
        location = /50x.html {
          root html;
        }
      }
