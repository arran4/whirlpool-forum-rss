# README

## Overview
This application scrapes the "Just In" section of ABC News and generates an RSS feed. 

## Installation

### Prerequisites
- **Go** (if building from source)

#### Optional
- **Apache or Nginx** (for serving RSS files)

### Build and Install

#### Install only (System level)

Grab the latest binary here, deb, rpm, etc: https://github.com/arran4/abc-justin-rss/releases/

#### Install and build as user (User)

Install go 1.23+

Run `go install`:
```bash
go install github.com/arran4/abc-justin-rss/cmd/abcjustinrss@latest
```

This installs to `$HOME/go/bin` (typically check with `go env`)

### Usage
#### Generate RSS Feed
```bash
abcjustinrss -output /var/www/localhost/htdocs/rss/abcjustinrss.xml
```

### Deployment

#### rc.d (Cron Job system level)
Add a cron job to run the script periodically:
1. Edit the root crontab:
   ```bash
   sudo crontab -e
   ```
2. Add the following line:
   ```bash
   */15 * * * * /usr/local/bin/abcjustinrss -output /var/www/localhost/htdocs/rss/abcjustinrss.xml
   ```

#### rc.d (Cron Job user level)
Add a cron job to run the script periodically:
1. Edit the user's crontab:
   ```bash
   crontab -e
   ```
2. Add the following line:
   ```bash
   */15 * * * * ~/go/bin/abcjustinrss -output ~/public_html/rss/abcjustinrss.xml
   ```

Remember to modify `~` with the correct value and `go/bin` too if you're using a custom go env location

#### systemd (as root)
1. Create a systemd service file at `/etc/systemd/system/abcjustinrss.service`:
```ini
[Unit]
Description=ABC News Just-in RSS Feed Creator

[Service]
Type=oneshot
ExecStart=/usr/bin/abcjustinrss -output /var/www/localhost/htdocs/rss/abcjustinrss.xml
User=apache
Group=apache
```

2. Create a systemd timer file at `/etc/systemd/system/everyhour@.timer`:

```ini
[Unit]
Description=Monthly Timer for %i service

[Timer]
OnCalendar=*-*-* *:00:00
AccuracySec=1h
RandomizedDelaySec=1h
Persistent=true
Unit=%i.service

[Install]
WantedBy=default.target
```

3. Reload systemd and start the service:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now everyhour@abcjustinrss.timer
   ```

#### systemd (as user)
1. Create a systemd service file at `$HOME/.config/systemd/user/abcjustinrss.service`:
```ini
[Unit]
Description=ABC News Just-in RSS Feed Creator

[Service]
Type=oneshot
ExecStart=$HOME/go/bin/abcjustinrss -output ${HOME}/public_html/rss/abcjustinrss.xml
```

Remember to modify $HOME with the correct value and `go/bin` too if you're using a custom go env location

2. Create a systemd timer file at `$HOME/.config/systemd/user/everyhour@.timer`:

```ini
[Unit]
Description=Monthly Timer for %i service

[Timer]
OnCalendar=*-*-* *:00:00
AccuracySec=1h
RandomizedDelaySec=1h
Persistent=true
Unit=%i.service

[Install]
WantedBy=default.target
```

3. Reload systemd and start the service:
   ```bash
   systemctl --user daemon-reload && systemctl --user enable --now everyhour@abcjustinrss.timer
   ```

#### Apache VirtualHost Configuration
##### User

Refer to documentation for setting up public_html directories

##### Enjoy

http://localhost/~$USERNAME/rss/abcjustinrss.xml

##### System

Add the following configuration to your Apache setup (e.g., `/etc/httpd/conf.d/rss.conf`):
```apache
<VirtualHost *:80>
    ServerName example.com
    DocumentRoot /var/www/localhost/htdocs/rss
    <Directory "/var/www/localhost/htdocs/rss">
        Options Indexes FollowSymLinks
        AllowOverride None
        Require all granted
    </Directory>
</VirtualHost>
```


#### Nginx Configuration
##### User

Refer to documentation for setting up public_html directories

##### System

Add this to your Nginx server block:
```nginx
server {
    listen 80;
    server_name example.com;

    location /rss/ {
        root /var/www/localhost/htdocs;
        autoindex on;
    }
}
```

