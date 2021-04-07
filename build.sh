#!/bin/bash

cat << EOF

      执行此脚本之前，请确认一下软件是否安装或者是否有现成的连接地址。

      若未没有请根据不同的系统，自行百度一下安装教程。

          1. git 最新版本即可
          1. MySQL >= 5.7
          2. Go >= 1.14
          3. Redis 最新版本即可
          4. node >= v12 （稳定版本）
          5. npm >= v6.14.8

EOF

# 判断目录是否存在，不存在则新建目录
isDirExist() {
  if [ ! -d "$1" ]; then
    mkdir -p $1
  fi
}

echo "确认 build 目录是否存在"
isDirExist "./build"
isDirExist "./build/log"
isDirExist "./build/template"

echo "开始迁移配置信息..."
isDirExist "./build/config"
cp -r ./config/db.sql ./build/config
cp -r ./config/settings.yml ./build/config/settings.yml
cp -r ./config/rbac_model.conf ./build/config/rbac_model.conf

echo "开始迁移静态文件..."
isDirExist "./build/static/scripts"
isDirExist "./build/static/template"
isDirExist "./build/static/uploadfile"
cp -r ./static/template/email.html ./build/static/template/email.html

# 编译前端程序，再此处需输入程序的访问地址，来进行前端程序的编译
read -p "请输入您的程序访问地址，例如：https://fdevops.com:8001，不可为空: " url
if [ -z "$url" ]; then
  echo "url输入不能为空"
  exit 1
fi

echo "请选择从哪里拉取前端代码，默认是gitee: "
cat << EOF

  1. gitee
  2. github
  3. 自定义拉取地址

EOF

read -p "请选择[1]: " ui
if [ -z "$ui" ]; then
  ui=1
fi

if [ $ui == 1 ]; then
  ui_address="https://gitee.com/yllan/ferry_web.git"
elif [ $ui == 2 ]; then
  ui_address="https://github.com/lanyulei/ferry_web.git"
elif [ $ui == 3 ]; then
  read -p "请输入拉取地址: " ui_address
else
  echo "选项不正确，请重新输入"
  exit 1
fi

echo "开始拉取前端程序..."
read -p "此处会执行 rm -rf ./ferry_web 的命令，若此命令不会造成当前环境的损伤则请继续，y/n[y] :" s
if [ ! -z "$s" ]; then
  if [ $s == "n" ]; then
    echo "结束此次编译"
    exit 1
  elif [ $s != "y" ]; then
    echo "结束此次编译"
    exit 1
  fi
fi

 if [ -d "./ferry_web" ]; then
   echo "请稍等，正在删除 ferry_web ..."
   rm -rf ./ferry_web
 fi
 git clone $ui_address

echo "替换程序访问地址..."
cat > ./ferry_web/.env.production << EOF
# just a flag
ENV = 'production'

# base api
VUE_APP_BASE_API = '$url'
EOF

echo "开始安装前端依赖..."
npm install -g cnpm --registry=https://registry.npm.taobao.org
cd ferry_web && cnpm install && npm run build:prod && cp -r web ../build/template

echo "\n需注意: 邮件服务器信息若是暂时没有，可暂时不修改，但是MySQL和Redis是必须配置正确的\n"
read -p "请确认是否配置MySQL、Redis及邮件服务器信息，配置文件地址: build/config/settings.yml，y/n[y]: " config_status
if [ ! -z "$config_status" ]; then
  if [ $config_status == "n" ]; then
    echo "结束此次编译"
    exit 1
  elif [ $config_status != "y" ]; then
    echo "结束此次编译"
    exit 1
  fi
fi

read -p "请确认是否创建配置文件中的MySQL库，y/n[y]: " mysql_db_status
if [ ! -z "$mysql_db_status" ]; then
  if [ $mysql_db_status == "n" ]; then
    echo "结束此次编译"
    exit 1
  elif [ $mysql_db_status != "y" ]; then
    echo "结束此次编译"
    exit 1
  fi
fi

cat <<EOF

    请选择程序运行的平台:

        1. Mac
        2. Linux
        3. Windows

EOF

read -p "请选择[2]: " run_platform
if [ -z "$run_platform" ]; then
  run_platform=2
fi

echo "开始编译后端程序..."

if [ $run_platform == 1 ]; then
  CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ferry main.go
elif [ $run_platform == 2 ]; then
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ferry main.go
elif [ $run_platform == 3 ]; then
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ferry main.go
else
  echo "没有您选择的平台，请确认"
  exit 1
fi

cp -r ferry ./build/

cd build && ./ferry init -c=config/settings.yml
if [ $? != 0 ]; then

  cat << EOF

    同步数据结构及数据失败，请确认 build/config/settings.yml 中数据库的配置是否正确。

    数据库配置信息正确后，可手动执行以下步骤，完成编译：

        # 1. 进入工作目录
        cd build

        # 2. 重新同步任务
        ./ferry init -c=config/settings.yml

        # 3. 启动服务
        ./ferry server -c=config/settings.yml

EOF

  exit 1

fi

echo "编译完成"

cat << EOF

    执行以下命令，启动程序：

        cd build
        ./ferry server -c=config/settings.yml

EOF
