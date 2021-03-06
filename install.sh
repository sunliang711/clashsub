#!/bin/bash
if [ -z "${BASH_SOURCE}" ];then
    this=${PWD}
else
    rpath="$(readlink ${BASH_SOURCE})"
    if [ -z "$rpath" ];then
        rpath=${BASH_SOURCE}
    fi
    this="$(cd $(dirname $rpath) && pwd)"
fi

export PATH=$PATH:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

user="${SUDO_USER:-$(whoami)}"
home="$(eval echo ~$user)"

# export TERM=xterm-256color

# Use colors, but only if connected to a terminal, and that terminal
# supports them.
if which tput >/dev/null 2>&1; then
  ncolors=$(tput colors 2>/dev/null)
fi
if [ -t 1 ] && [ -n "$ncolors" ] && [ "$ncolors" -ge 8 ]; then
    RED="$(tput setaf 1)"
    GREEN="$(tput setaf 2)"
    YELLOW="$(tput setaf 3)"
    BLUE="$(tput setaf 4)"
    CYAN="$(tput setaf 5)"
    BOLD="$(tput bold)"
    NORMAL="$(tput sgr0)"
else
    RED=""
    GREEN=""
    YELLOW=""
    CYAN=""
    BLUE=""
    BOLD=""
    NORMAL=""
fi

_err(){
    echo "$*" >&2
}

_runAsRoot(){
    cmd="${*}"
    local rootID=0
    if [ "${EUID}" -ne "${rootID}" ];then
        # echo -n "Not root, try to run '${cmd}' as root.."
        # or sudo sh -c ${cmd} ?
        # if eval "sudo ${cmd}";then
        if sudo sh -c "${cmd}";then
            # echo "ok"
            return 0
        else
            # echo "failed"
            return 1
        fi
    else
        # or sh -c ${cmd} ?
        eval "${cmd}"
    fi
}

rootID=0
function _root(){
    if [ ${EUID} -ne ${rootID} ];then
        echo "Need run as root!"
        echo "Requires root privileges."
        exit 1
    fi
}

ed=vi
if command -v vim >/dev/null 2>&1;then
    ed=vim
fi
if command -v nvim >/dev/null 2>&1;then
    ed=nvim
fi
if [ -n "${editor}" ];then
    ed=${editor}
fi
###############################################################################
# write your code below (just define function[s])
# function is hidden when begin with '_'
###############################################################################
# TODO
link="https://source711.oss-cn-shanghai.aliyuncs.com/clashsub/clashsub.tar.bz2"

install(){
    local installDir=${1:?'missing instal location'}
    if [ ! -d ${installDir} ];then
        echo "${installDir} not exist!"
	mkdir -p ${installDir}
    fi

    installDir=$(cd ${installDir} && pwd)
    echo "Install dir: ${installDir}"

    local downloadDir=/tmp/clashSub
    [ ! -d ${downloadDir} ] && mkdir ${downloadDir}
    cd ${downloadDir}
    curl -LO ${link}
    tarName=${link##*/}
    echo "tarName: ${tarName}"
    echo "Extract clashsub to ${installDir} ..."
    ls -l ${tarName}
    tar -C ${installDir} -xjvf ${tarName}

    local start="${installDir}/clashsub/clashsub -c config.yaml"
    local pwd=${installDir}/clashsub
    cd ${this}
    sed -e "s|<START>|${start}|g" \
        -e "s|<PWD>|${pwd}|g" \
        clashsub.service >/tmp/clashsub.service
    _runAsRoot "mv /tmp/clashsub.service /etc/systemd/system"
    _runAsRoot "systemctl daemon-reload"
    _runAsRoot "systemctl enable --now clashsub"
}

em(){
    $ed $0
}

###############################################################################
# write your code above
###############################################################################
function _help(){
    cd "${this}"
    cat<<EOF2
Usage: $(basename $0) ${bold}CMD${reset}

${bold}CMD${reset}:
EOF2
    # perl -lne 'print "\t$1" if /^\s*(\w+)\(\)\{$/' $(basename ${BASH_SOURCE})
    # perl -lne 'print "\t$2" if /^\s*(function)?\s*(\w+)\(\)\{$/' $(basename ${BASH_SOURCE}) | grep -v '^\t_'
    perl -lne 'print "\t$2" if /^\s*(function)?\s*(\S+)\s*\(\)\s*\{$/' $(basename ${BASH_SOURCE}) | perl -lne "print if /^\t[^_]/"
}

case "$1" in
     ""|-h|--help|help)
        _help
        ;;
    *)
        "$@"
esac
