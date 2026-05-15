#!/bin/env bash

# Hybrid display power management for CachyOS (X11 + KDE Wayland)
# Always runs display commands as the logged-in user, not as root.
# Usage: cachyos-wake.sh <on|off|reset> [display_number]

ACTION=$1
DISPLAY_NUM=$2
LOG="/tmp/cachyos-wake.log"

log() { echo "[$(date '+%H:%M:%S')] $*" >> "$LOG"; }

# Find active user sessions (UID, username, session type, display)
get_sessions() {
    loginctl list-sessions --no-legend 2>/dev/null | while read -r sid uid user seat rest; do
        class=$(loginctl show-session "$sid" -p Class --value 2>/dev/null) || continue
        [ "$class" != "user" ] && continue
        type=$(loginctl show-session "$sid" -p Type --value 2>/dev/null || true)
        display=$(loginctl show-session "$sid" -p Display --value 2>/dev/null || true)
        seat=$(loginctl show-session "$sid" -p Seat --value 2>/dev/null || true)
        echo "$sid|$uid|$user|$type|$display|$seat"
    done
}

run_as_user() {
    local uid=$1 user=$2
    shift 2
    if [ "$(id -u)" = "$uid" ]; then
        DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/$uid/bus" \
        XDG_RUNTIME_DIR="/run/user/$uid" "$@"
    else
        sudo -u "$user" \
            DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/$uid/bus" \
            XDG_RUNTIME_DIR="/run/user/$uid" "$@" 2>/dev/null || true
    fi
}

do_wake_x11() {
    local user=$1 uid=$2 display=$3 state=$4
    log "X11: xset -display $display dpms force $state"
    run_as_user "$uid" "$user" xset -display "$display" dpms force "$state"
}

do_wake_wayland() {
    local user=$1 uid=$2 state=$3
    if ! command -v kscreen-doctor &>/dev/null; then
        return
    fi

    if [ "$state" = "on" ]; then
        WAKER="$(dirname "$0")/wake-inject"
        if [ -x "$WAKER" ]; then
            log "wake via $WAKER"
            "$WAKER"
        fi
        log "restore brightness"
        timeout 5 run_as_user "$uid" "$user" \
            kscreen-doctor "output.*.brightness.75" 2>/dev/null || true
    elif [ "$state" = "off" ]; then
        log "brightness 0 on all outputs"
        timeout 5 run_as_user "$uid" "$user" \
            kscreen-doctor "output.*.brightness.0" 2>/dev/null || true
    fi
}

do_reset_x11() {
    local user=$1 uid=$2 display=$3 timeout=$4
    log "X11 reset: xset dpms $timeout"
    run_as_user "$uid" "$user" xset -display "$display" dpms "$timeout" "$timeout" "$timeout" 2>/dev/null || true
    run_as_user "$uid" "$user" xset -display "$display" s "$timeout" "$timeout" 2>/dev/null || true
}

log "=== cachyos-wake.sh $* ==="

get_sessions | while IFS='|' read -r sid uid user type display seat; do
    [ -z "$uid" ] && continue
    [ -n "$DISPLAY_NUM" ] && [ -n "$display" ] && [[ "$display" != *":$DISPLAY_NUM"* ]] && continue

    case "$ACTION" in
        on|off)
            if [ "$type" = "x11" ]; then
                do_wake_x11 "$user" "$uid" "$display" "$ACTION"
            elif [ "$type" = "wayland" ]; then
                do_wake_wayland "$user" "$uid" "$ACTION"
            fi
            ;;
        reset|[0-9]*)
            if [[ "$ACTION" =~ ^[0-9]+$ ]]; then
                local timeout="$ACTION"
            else
                local timeout="${DISPLAY_NUM:-600}"
            fi
            if [ "$type" = "x11" ]; then
                do_reset_x11 "$user" "$uid" "$display" "$timeout"
            fi
            # Wayland: KWin manages DPMS timeouts internally, no-op
            ;;
    esac
done

exit 0
