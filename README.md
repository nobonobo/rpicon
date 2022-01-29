# rpicon

## setup RasPi OS

- use RapberryPi Imager
- config hostname and Wi-Fi, ssh
- sudo apt update
- sudo apt install -y python3-pip
- sudo pip3 install nxbt

## build and install

```
make generate build
scp {rpicon,procon.py} pi@hostname.local:/home/pi/
```

/etc/systemd/system/rpicon.service

```
[Unit]
Description=Rpi ProCon Emulator
After=bluetooth.target

[Service]
ExecStart=/home/pi/rpicon -script=/home/pi/procon.py
Restart=always

[Install]
WantedBy=multi-user.target
```

```
sudo systemctl daemon-reload
sudo systemctl enable rpicon
sudo systemctl start rpicon
```

## using

- open http://hostname.local in browser.
- connect usb or bluetooth for game controller.
- press button for game controller.
- choise game controller from top left select box.
- choice profile for game controller.
- Nintendo Switch standby for connect new controller.
- top-right toggle switch on.
- connect and few seconds later, you can see controller.
- you can use game controller as procon!
