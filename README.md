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
      "Inputs": {
        "hdmi": "0x10",
        "usbc": "0x1c"
      }
    }
  },
  "ddcutil": {
    "bin": "/usr/bin/ddcutil"
  }
}
```

## Sample message

```json
[
  {
  "monitor": "left",
  "input": "hdmi"
  },
  {
    "monitor": "right",
    "input": "usbc"
  }
]
```