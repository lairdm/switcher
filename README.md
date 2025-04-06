# Monitor Switcher

Listens to a MQTT topic for commands to switch monitor inputs. It's designed to allow Home Assistant to send commands so I can use a small RPi powered touch screen beside my desk as a KVM control.

Do not attempt to learn anything from this code, it was my first attempt with Golang and completed in about 5 hours. I haven't even started using it in production yet. :)

## Sample Config

```json
{
  "mqtt": {
    "broker": "192.168.1.10",
    "port": 1883,
    "topic": "/home/office/monitors"
  },
  "monitors": {
   "left": {
     "Serial": "123456",
     "Inputs": {
       "hdmi": "0x10",
       "dp": "0x11"
     }
   },
    "right": {
      "Serial": "1234567",
      "Display": "0",
      "Inputs": {
        "hdmi": "0x10",
        "usbc": "0x1c"
      }
    }
  },
  "ddcutil": {
    "bin": "/usr/bin/ddcutil"
  },
  "monitorctl": {
    "bin": "scripts/wakeUserDisplays.sh"
  },
  "blanking": {
    "bin": "scripts/resetBlanking.sh"
  },
  "amixer": {
    "bin": "/usr/bin/amixer"
  },
  "xlock": {
    "bin": "cinnamon-unlock-desktop"
  },
  "lirc": {
    "bin": "/usr/bin/irsend",
    "remote": "NEC"
  }
}
```

## Sample message

```json
[
  {
    "command": "monitor",
    "data":
      {
      "monitor": "left",
      "input": "hdmi",
      "power": "wake"
      }
  },
  {
    "command": "monitor",
    "data":
      {
        "monitor": "right",
        "input": "usbc"
      }
  },
  {
    "command": "speaker",
    "data":
      {
        "volume": "up"
      }
  },
  {
    "command": "lock",
    "data":
      {
        "unlock": true
      }
  }
]
```

## Notes

ddcutil seems a little finicky depending on the video card. ddcutil segfaults immediately on my laptop with external monitors plugged in. Trying Arch as described in the [ddcutil FAQ](https://www.ddcutil.com/faq/) gives a partial communication failure - it sees the monitors but isn't able to communicate sufficiently to send commands.

On my old desktop the Arch method is able to issue commands.

I had contemplated it could be the input type (USB-C DP alt-mode on the laptop, HDMI on the desktop), however ddcutil works perfectly fine from the normal shell. Anyhow, we'll see how this works on the desktop when I'm ready to go live, otherwise I might just run the application using supervisord.

## Personal note

> CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on /usr/lib/go-1.22/bin/go build -a -o switcher main.go
