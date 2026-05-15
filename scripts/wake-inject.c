#define _GNU_SOURCE
#include <linux/uinput.h>
#include <unistd.h>
#include <fcntl.h>
#include <string.h>
#include <stdio.h>

// Build: make -C scripts
//   or: gcc -std=gnu99 -O2 -o wake-inject wake-inject.c
// Requires /dev/uinput (root or uinput group)

int main() {
    int fd = open("/dev/uinput", O_WRONLY | O_NONBLOCK);
    if (fd < 0) { perror("open /dev/uinput"); return 1; }

    ioctl(fd, UI_SET_EVBIT, EV_REL);
    ioctl(fd, UI_SET_RELBIT, REL_X);
    ioctl(fd, UI_SET_RELBIT, REL_Y);
    ioctl(fd, UI_SET_EVBIT, EV_KEY);
    ioctl(fd, UI_SET_KEYBIT, BTN_LEFT);
    ioctl(fd, UI_SET_EVBIT, EV_SYN);

    struct uinput_setup us = {0};
    strcpy(us.name, "wake-injector");
    us.id.bustype = BUS_USB;
    ioctl(fd, UI_DEV_SETUP, &us);
    ioctl(fd, UI_DEV_CREATE);
    usleep(200000);

    struct input_event ev = {0};
    ev.type = EV_REL;
    ev.code = REL_X;
    ev.value = 1;
    write(fd, &ev, sizeof(ev));

    ev.type = EV_SYN;
    ev.code = SYN_REPORT;
    ev.value = 0;
    write(fd, &ev, sizeof(ev));

    usleep(50000);
    ioctl(fd, UI_DEV_DESTROY);
    close(fd);
    return 0;
}
