#!/bin/bash
#=============================================================================
#
# Author: zhouzhibo
#
# id : AL1009
#
# Last modified:	2021-04-12 21:18
#
# Filename:		build.sh
#
# Description: 
#
#=============================================================================


CUR_USE=$(whoami)    
BASE_DIR=$(cd $(dirname $0) >/dev/null 2>&1 && pwd)    
FILE_BASE_NAME=$(basename $0)    

function echo_red() {
    echo -e "\033[1;31m$1\033[0m"
}

function echo_green() {
    echo -e "\033[1;32m$1\033[0m"
}

function echo_yellow() {
    echo -e "\033[1;33m$1\033[0m"
}

function echo_done() {
    echo "$(gettext 'complete')"
}

function read_from_input() {
    var=$1
    msg=$2
    choices=$3
    default=$4
    if [[ ! -z "${choices}" ]]; then
        msg="${msg} (${choices}) "
    fi
    if [[ -z "${default}" ]]; then
        msg="${msg} ($(gettext 'no default'))"
    else
        msg="${msg} ($(gettext 'default') ${default})"
    fi
    echo -n "${msg}: "
    read input
    if [[ -z "${input}" && ! -z "${default}" ]]; then
        export ${var}="${default}"
    else
        export ${var}="${input}"
    fi
}

function usage {
    cat << EOF
  echo -e "\nUsage: $0 (install|start|stop)"    
    echo "Examle:     
    bash $0 install 
EOF
}

function check_soft {
    local _soft_name=$1
    command -v ${_soft_name} > /dev/null || {
    echo_red "请安装 ${_soft_name} 后再执行本脚本安装ferry。"
        exit 1
    }
}


function prepare_check {
    cat << EOF

      执行此脚本之前，请确认一下软件是否安装或者是否有现成的连接地址。

      若未没有请根据不同的系统，自行百度一下安装教程。

          1. git 最新版本即可
          1. MySQL >= 5.7
          2. Go >= 1.14
          3. Redis 最新版本即可
          4. node >= v12 （稳定版本）
          5. npm >= v6.14.8
          6. 安装docker

EOF
    check_soft docker
    check_soft git
    check_soft go
    check_soft npm
}

function isDirExist {
    # 判断目录是否存在，不存在则新建目录
    local _dir_name=$1
    [ ! -d "$1" ] && mkdir -p $1
}

function mk_ferry_dir {
    echo "检查创建确认 build 以及子目录是否存在"
    isDirExist "${BASE_DIR}/build/log"
    isDirExist "${BASE_DIR}/build/template"
    isDirExist "${BASE_DIR}/build/config"
}

function init(){
    mk_ferry_dir
    echo_green "\n>>> $(gettext '开始迁移配置信息...')"
    [ -f "${BASE_DIR}/config/db.sql" ] && cp -pf ${BASE_DIR}/config/db.sql ${BASE_DIR}/build/config
    [ -f "${BASE_DIR}/config/settings.yml" ] && cp -pf ${BASE_DIR}/config/settings.yml ${BASE_DIR}/build/config
    [ -f "${BASE_DIR}/config/rbac_model.conf" ] && cp -pf ${BASE_DIR}/config/rbac_model.conf ${BASE_DIR}/build/config

    echo_green "\n>>> $(gettext '开始迁移静态文件...')"
    isDirExist "${BASE_DIR}/build/static/scripts"
    isDirExist "${BASE_DIR}/build/static/template"
    isDirExist "${BASE_DIR}/build/static/uploadfile"
    [ -f "${BASE_DIR}/static/template/email.html" ] && cp -pf ${BASE_DIR}/static/template/email.html ${BASE_DIR}/build/static/template/email.html
    if [ -f "${BASE_DIR}/build/config/settings.yml" ];then
        CONFIG_FILE=${BASE_DIR}/build/config/settings.yml  
    else
        echo_red "配置文件: ${BASE_DIR}/build/config/settings.yml 不存在，请检查。"
        exit 1
    fi
}

function config_mysql {
    echo_green "\n>>> $(gettext '需注意: 邮件服务器信息若是暂时没有，可暂时不修改，但是MySQL和Redis是必须配置正确的')"
    read_from_input confirm "$(gettext '请确认是否安装MySQL')?" "y/n" "y"

    if [[ "${confirm}" == "y" ]]; then
        echo ""
        echo "请在此处暂停一下，将数据库配置信息，写入到配置文件中，${BASE_DIR}/build/config/settings.yml，<settings.database> 下面数据库相关配置。"
    else
        echo_red "未安装Mysql结束此次编译"
        exit 1
    fi
}

function config_redis {
    echo_green "\n>>> $(gettext '回车前请确保你已经安装了Redis,且启动服务')"
    read_from_input confirm "$(gettext '请确认是否安装Redis')?" "y/n" "y"

    if [[ "${confirm}" == "y" ]]; then
        echo ""
        echo "请在此处暂停一下，将 Redis 配置信息，写入到配置文件中，${BASE_DIR}/build/config/settings.yml，<settings.redis> 下面是Redis相关配置，若是不知道如何配置URL，可自行百度一下。"
    else
        echo_red "未安装Redis结束此次编译"
        exit 1
    fi
}

function get_variables {
    read_from_input front_url "$(gettext '请输入您的程序访问地址: ')" "" ""
    read_from_input front_clone_from "$(gettext '请选择从哪里拉取前端代码，默认是gitee: 1:gitee, 2: github, 3:自定义地址')" "" "1"

    if [ $front_clone_from == 1 ]; then
        ui_address="https://gitee.com/yllan/ferry_web.git"
    elif [ $front_clone_from == 2 ]; then
        ui_address="https://github.com/lanyulei/ferry_web.git"
    else
        ui_address=${front_clone_from}
    fi

    config_mysql
    config_redis
    echo_done

}

function config_front {
    echo_green "\n>>> $(gettext '替换程序访问地址...')"
    cat > ${BASE_DIR}/ferry_web/.env.production << EOF
# just a flag
ENV = 'production'

# base api
VUE_APP_BASE_API = '$front_url'
EOF

}

function install_front {
    echo_green "\n>>> $(gettext '开始拉取前端程序...')"
    read_from_input confirm "$(gettext '此处会执行 rm -rf ./ferry_web 的命令，若此命令不会造成当前环境的损伤则请继续')?" "y/n[y]" "y"
    if [[ "${confirm}" != "y" ]]; then
        echo_red "结束此次编译"
        exit 1
    fi


    if [ -d "${BASE_DIR}/ferry_web" ]; then
        echo_green "\n>>> $(gettext '请稍等，正在删除 ferry_web ...')"
        rm -rf ${BASE_DIR}/ferry_web
    fi
    git clone $ui_address 

    if [ "$?" -ne 0 ];then
        echo_red "克隆代码失败，请检查git地址: ${ui_address}或者网络质量"
        exit 1
    fi
    config_front
    echo_green "\n>>> $(gettext '开始安装前端依赖...')"
    cnpm_base_dir=$(dirname $(dirname $(which npm)))
    npm install -g cnpm --registry=https://registry.npm.taobao.org --prefix ${cnpm_base_dir}
    cd ferry_web && cnpm install && npm run build:prod && cp -r web ../build/template && cp -r web/static/* ../build/static/

}

function install_backend {
    echo_green "\n>>> $(gettext '开始编译后端程序...')"

    cd ${BASE_DIR} 
    if [ "$(uname)" == "Darwin" ];then
        CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ferry main.go
        result=$?
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ];then
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ferry main.go
        result=$?
    elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ];then
        CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ferry main.go
        result=$?
    fi

    cd - &>/dev/null

    if [ "$result" -ne 0 -o ! -f "${BASE_DIR}/ferry" ];then
        echo_red "编译后端代码失败，退出安装"
        exit 1
    fi
    cp -r ${BASE_DIR}/ferry ${BASE_DIR}/build/
    cd ${BASE_DIR}/build 
    ${BASE_DIR}/build/ferry init -c=config/settings.yml
    cd - &>/dev/null
}

function install_app() {
    prepare_check
    init
    get_variables
    install_front
    install_backend
}

function start_backend {
    cd ${BASE_DIR}/build
    ./ferry server -c=config/settings.yml 
}

function main {
    action=${1-}
    target=${2-}
    args=("$@")

    case "${action}" in
        install)
            install_app
            ;;
        uninstall)
            echo "功能暂未开发, 敬请期待。"
            ;;
        start)
            start_backend
            ;;
        stop)
            echo "功能暂未开发, 敬请期待。"
            ;;
        help)
            usage
            ;;
        --help)
            usage
            ;;
        -h)
            usage
            ;;
        *)
            echo "No such command: ${action}"
            usage
            ;;
    esac
}

main "$@"
