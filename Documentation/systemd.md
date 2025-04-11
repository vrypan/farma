# Running farma as a linux service

Assuming you have downloaded [the latest binary for your platform](https://github.com/vrypan/farma/releases),
you have extracted the binary, and moved it in `/usr/local/bin/farma`.

## 1. Create user and required directories

```bash
# Create a new user ("farma" in this case), that will run the service
sudo useradd -r -s /bin/false farma

# /etc/farma will hold the config file
sudo mkdir /etc/farma
sudo chown farma:farma /etc/farma

# /usr/lib/farma will host the database
sudo mkdir /var/lib/farma
sudo chown farma:farma /var/lib/farma
```

## 2. Setup
Create the necessary files in `/etc/farma/`:

```
sudo -u farma XDG_CONFIG_HOME=/etc/ /usr/local/bin/farma setup
```

Set the db path in config.yaml:
```
sudo -u farma XDG_CONFIG_HOME=/etc/ /usr/local/bin/farma config set db.path /var/lib/farma/
```


## 3. Configure systemd

First create the service file

```
sudo nano /etc/systemd/system/farma.service
```

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
# Environment=FARMA_DB_PATH=/var/lib/farma/

[Install]
WantedBy=multi-user.target
```

Then enable and start the service
```
sudo systemctl daemon-reload
sudo systemctl start farma.service
```

Check if the service is running: 
```
cat /var/log/farma.log
```

To automatically start the service on boot, enable it:
```
sudo systemctl enable farma.service
```

## Notes

If you have the above configuration, and you are using the comman-line tools like `farma frames-list`,
make sure you export `XDG_CONFIG_HOME` accordingly,

```
export XDG_CONFIG_HOME=/etc/
```

or export `FARMA_PRIVATE_KEY`

```
export FARMA_PRIVATE_KEY=<your private key>
```

You can get the private key from `/etc/farma/config.yaml`. If you have securely copied it somewhere, and you don't
want it to be in the config file, you can remove if (set it to empty string). In this case, you will have to provide it
to commands or tools that use `farma` manually.
