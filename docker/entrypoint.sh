#!/bin/sh
set -e
if [[ ! -f /opt/workflow/ferry/config/settings.yml ]]
then
    cp /opt/workflow/ferry/default_config/* /opt/workflow/ferry/config/
fi
if [[ -f /opt/workflow/ferry/config/needinit ]]
then
    /opt/workflow/ferry/ferry init -c=/opt/workflow/ferry/config/settings.yml
    rm -f /opt/workflow/ferry/config/needinit
fi
/opt/workflow/ferry/ferry server -c=/opt/workflow/ferry/config/settings.yml
