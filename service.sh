#!/bin/bash

# worktime 服务管理脚本

WORKTIME_CMD="./worktime"
WORKTIME_ARGS="start"  # 启动web服务的参数
SERVICE_NAME="worktime"

start() {
    echo "正在启动 $SERVICE_NAME 服务..."
    if pgrep -f $SERVICE_NAME > /dev/null; then
        echo "$SERVICE_NAME 服务已经在运行中!"
        return 1
    fi

    chmod +x $WORKTIME_CMD
    nohup $WORKTIME_CMD $WORKTIME_ARGS > /dev/null 2>&1 &

    # 检查服务是否启动成功
    sleep 2
    if pgrep -f $SERVICE_NAME > /dev/null; then
        echo "$SERVICE_NAME 服务启动成功!"
    else
        echo "$SERVICE_NAME 服务启动失败，请检查!"
        return 1
    fi
}

stop() {
    echo "正在停止 $SERVICE_NAME 服务..."
    if ! pgrep -f $SERVICE_NAME > /dev/null; then
        echo "$SERVICE_NAME 服务未运行!"
        return 1
    fi

    pkill -f $SERVICE_NAME

    # 等待服务完全停止
    sleep 3

    if pgrep -f $SERVICE_NAME > /dev/null; then
        echo "$SERVICE_NAME 服务停止失败，可能需要手动kill!"
        return 1
    else
        echo "$SERVICE_NAME 服务已停止!"
    fi
}

restart() {
    echo "正在重启 $SERVICE_NAME 服务..."
    stop
    start
}

status() {
    if pgrep -f $SERVICE_NAME > /dev/null; then
        echo "$SERVICE_NAME 服务正在运行"
        pgrep -f $SERVICE_NAME
    else
        echo "$SERVICE_NAME 服务未运行"
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    *)
        echo "用法: $0 {start|stop|restart|status}"
        exit 1
        ;;
esac

exit 0