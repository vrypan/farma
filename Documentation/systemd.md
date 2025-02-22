# Running farma as a linux service

Assuming you have downloaded [the latest binary for your platform](https://github.com/vrypan/farma/releases),
you have extracted the binary, and moved it in `/usr/local/bin/farma`.

## 1. Create user and required directories

Create a new user (`farma` in this case), that will run the service

`sudo useradd -r -s /bin/false farma`

Then create the `/etc/farma` (config files) and `/var/lib/farma` (data files) and set the ownership to `farma`.

`sudo mkdir /etc/farma`

`sudo chown farma:farma /etc/farma`

`sudo mkdir /var/lib/farma`

`sudo chown farma:farma /var/lib/farma`

## 2. Run the setup script
`sudo -u farma XDG_CONFIG_HOME=/etc/ /usr/local/bin/farma setup`

This will add the necessary files to

## 3. Configure systemd

First create the service file

`sudo nano /etc/systemd/system/farma.service`

Add the following content to the file
```
[Unit]
Description=Farma Server
After=network.target

[Service]
ExecStart=/usr/local/bin/farma server
Restart=always
User=farma
WorkingDirectory=/tmp
StandardOutput=append:/var/log/farma.log
StandardError=append:/var/log/farma.log
Environment=XDG_CONFIG_HOME=/etc/
Environment=FARMA_DB_PATH=/var/lib/farma/

[Install]
WantedBy=multi-user.target
```

Then enable and start the service
`sudo systemctl daemon-reload`

`sudo systemctl start farma.service`

Check if the service is running: `cat /var/log/farma.log`

To automatically start the service on boot, enable it:
`sudo systemctl enable farma.service`

## Notes

If you have the above configuration, and you are using the comman-line tools like `farma frames-list`,
make sure you export `XDG_CONFIG_HOME` accordingly,

`export XDG_CONFIG_HOME=/etc/`

or export `FARMA_PRIVATE_KEY`

`export FARMA_PRIVATE_KEY=<your private key>`

You can get the private key from `/etc/farma/config.yaml`
