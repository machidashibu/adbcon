# Supported commands

* Android Debug Bridge version 1.0.40
* Version 28.0.2-5303910

Memo for development.

`body` column is `w/o` means no body request. i.e. not need `-s` option.

| Command | Implemented | body | Note |
| ------- | ----------- | ---- | ---- |
| devices | `/api/devices` | w/o | auto polling |
| help | x | w/o | |
| version | x | w/o | |
| connect | x | | |
| disconnect | x | | |
| forward | x | | |
| ppp | x | | |
| reverse | x | | |
| push | `/api/adb/push` | | |
| pull | `/api/adb/pull` | | |
| sync | x | | |
| shell | `/api/adb/shell` | | |
| emu | x | | |
| install | x | | |
| install-multiple | x | | |
| install-multi-package | x | | |
| uninstall | x | | |
| bugreport | x | | |
| jdwp | x | | |
| logcat | `/api/adb/logcat` | | disable GUI because screen freeze. |
| disable-verity | x | | |
| enable-verity | x | | |
| keygen | x | | |
| wait-for[-TRANSPORT]-STATE | x | | `TRANSPORT`: usb, local, or any (default)<br> `STATE`: device, recovery, sideload, or bootloader |
| get-state | x | | |
| get-serialno | x | | |
| get-devpath | x | | |
| remount | x | | |
| reboot | `/api/adb/reboot` | | |
| sideload | x | | |
| root | `/api/adb/root` | | |
| unroot | `/api/adb/unroot` | | |
| usb | x | | |
| tcpip | x | | |
| start-server | `/api/adb/start-server` | w/o | |
| kill-server | `/api/adb/start-server` | w/o ||
| reconnect | x | w/o | |
