#!/bin/env bash

# Hybrid display power management for CachyOS (X11 + KDE Wayland)
# Always runs display commands as the logged-in user, not as root.
# Usage: cachyos-wake.sh <on|off|reset> [display_number]

ACTION=$1
DISPLAY_NUM=$2

# Find active user sessions (UID, username, display)
get_sessions() {
    loginctl list-sessions --no-legend 2>/dev/null | while read -r sid uid user seat rest; do
        # Skip non-graphical sessions
        class=$(loginctl show-session "$sid" -p Class --value 2>/dev/null) || continue
        [ "$class" != "user" ] && continue
        display=$(loginctl show-session "$sid" -p Display --value 2>/dev/null || true)
        seat=$(loginctl show-session "$sid" -p Seat --value 2>/dev/null || true)
        echo "$sid $uid $user $display $seat"
    done
}

run_as_user() {
    local uid=$1 user=$2
    shift 2
    sudo -u "$user" DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/$uid/bus" "$@" 2>/dev/null || true
}

do_wake_x11() {
    local user=$1 uid=$2 display=$3 state=$4
    run_as_user "$uid" "$user" xset -display "$display" dpms force "$state"
}

do_wake_wayland() {
    local user=$1 uid=$2 state=$3
    if command -v kscreen-doctor &>/dev/null; then
        if [ "$state" = "on" ] || [ "$state" = "reset" ]; then
            run_as_user "$uid" "$user" kscreen-doctor output.*.dpms.on
        elif [ "$state" = "off" ]; then
            run_as_user "$uid" "$user" kscreen-doctor output.*.dpms.off
        fi
    fi
}

do_reset_x11() {
    local user=$1 uid=$2 display=$3 timeout=$4
    run_as_user "$uid" "$user" xset -display "$display" dpms "$timeout" "$timeout" "$timeout"
    run_as_user "$uid" "$user" xset -display "$display" s "$timeout" "$timeout"
}

get_sessions | while read -r sid uid user display seat; do
    [ -z "$uid" ] && continue
    [ -n "$DISPLAY_NUM" ] && [[ "$display" != *":$DISPLAY_NUM"* ]] && continue

    case "$ACTION" in
        on|off)
            if command -v xset &>/dev/null; then
                do_wake_x11 "$user" "$uid" "$display" "$ACTION"
            else
                do_wake_wayland "$user" "$uid" "$ACTION"
            fi
            ;;
        reset)
            local timeout="${DISPLAY_NUM:-600}"
            if command -v xset &>/dev/null; then
                do_reset_x11 "$user" "$uid" "$display" "$timeout"
            else
                do_wake_wayland "$user" "$uid" on
            fi
            ;;
    esac
done

exit 0
