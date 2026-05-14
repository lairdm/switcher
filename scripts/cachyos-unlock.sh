#!/bin/env bash

# Unlocks desktop session on CachyOS/Arch (KDE Plasma, GNOME, etc.)
# Always runs commands as the logged-in user.

# Method 1: unlock each session via logind
loginctl list-sessions --no-legend 2>/dev/null | while read -r sid uid user seat rest; do
    loginctl unlock-session "$sid" 2>/dev/null || true
done

# Method 2: KDE kscreenlocker via qdbus (needs user D-Bus session)
loginctl list-sessions --no-legend 2>/dev/null | while read -r sid uid user seat rest; do
    bus="unix:path=/run/user/$uid/bus"
    sudo -u "$user" DBUS_SESSION_BUS_ADDRESS="$bus" \
        qdbus org.kde.screensaver /ScreenSaver Unlock 2>/dev/null || true
    sudo -u "$user" DBUS_SESSION_BUS_ADDRESS="$bus" \
        qdbus org.freedesktop.ScreenSaver /ScreenSaver Unlock 2>/dev/null || true
done

exit 0
