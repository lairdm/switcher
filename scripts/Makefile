CFLAGS = -std=gnu99 -O2
LDLIBS =

all: wake-inject

wake-inject: wake-inject.c
	$(CC) $(CFLAGS) $(LDFLAGS) -o $@ $< $(LDLIBS)

clean:
	rm -f wake-inject

.PHONY: all clean
