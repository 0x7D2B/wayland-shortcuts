hello let's install wayland-shortcuts

it's probably these commands but i'm not sure

1. let's make a group
```
sudo groupadd --force uinput
```

2. let's make a new user
```
sudo useradd --no-create-home --groups uinput wayland-shortcuts
```

3. let's get uinput perms
```
echo 'KERNEL=="uinput", GROUP="uinput", MODE:="0660"' | sudo tee /etc/udev/rules.d/99-uinput.rules
```

4. let's copy our cool systemd service
```
sudo cp wayland-shortcuts.service /etc/systemd/system/wayland-shortcuts.service
```

5. it's a cool serivce pls admire it
```
cat wayland-shortcuts.service
```

6. enable it with systemd maybe
```
sudo systemctl daemon-reload
sudo systemctl enable wayland-shortcuts
```

7. probably ready to go once you reboot but don't take my word for it i haven't tried it 

if it breaks check perms and groups for /dev/uinput

¯\\\_(ツ)\_/¯
