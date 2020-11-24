#!/bin/sh
# vim:sw=4:ts=4:et

grep -r -o "http://fdevops.com:8001" /opt/web |awk -F ':' '{print $1}' | xargs sed -i s'#http://fdevops.com:8001#http://192.168.1.1:8001#g'
grep -r -o "VUE_APP_BASE_API" /opt/web |awk -F ':' '{print $1}' | xargs sed -i s'#VUE_APP_BASE_API#http://192.168.1.1:8001#g'
grep -r -o "localhost" /opt/web/static/web/js |awk -F ':' '{print $1}' | xargs sed -i s'#localhost#192.168.1.1#g'

set -e

if [ -z "${NGINX_ENTRYPOINT_QUIET_LOGS:-}" ]; then
    exec 3>&1
else
    exec 3>/dev/null
fi

if [ "$1" = "nginx" -o "$1" = "nginx-debug" ]; then
    if /usr/bin/find "/docker-entrypoint.d/" -mindepth 1 -maxdepth 1 -type f -print -quit 2>/dev/null | read v; then
        echo >&3 "$0: /docker-entrypoint.d/ is not empty, will attempt to perform configuration"

        echo >&3 "$0: Looking for shell scripts in /docker-entrypoint.d/"
        find "/docker-entrypoint.d/" -follow -type f -print | sort -n | while read -r f; do
            case "$f" in
                *.sh)
                    if [ -x "$f" ]; then
                        echo >&3 "$0: Launching $f";
                        "$f"
                    else
                        # warn on shell scripts without exec bit
                        echo >&3 "$0: Ignoring $f, not executable";
                    fi
                    ;;
                *) echo >&3 "$0: Ignoring $f";;
            esac
        done

        echo >&3 "$0: Configuration complete; ready for start up"
    else
        echo >&3 "$0: No files found in /docker-entrypoint.d/, skipping configuration"
    fi
fi

exec "$@"