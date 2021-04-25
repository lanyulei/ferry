#!/bin/bash
set -e
if [[ ! -f /opt/ferry/config/settings.yml ]]
then
    cp /opt/ferry/default_config/* /opt/ferry/config/
fi
if [[ -f /opt/ferry/config/needinit ]]
then
    /opt/ferry/ferry init -c=/opt/ferry/config/settings.yml
    rm -f /opt/ferry/config/needinit
fi
/opt/ferry/ferry server -c=/opt/ferry/config/settings.yml
