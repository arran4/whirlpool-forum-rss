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
   */15 * * * * ./go/bin/abcjustinrss -output ~/public_html/rss/abcjustinrss.xml
   ```

#### systemd (as root)
1. Create a systemd service file at `/etc/systemd/system/abcjustinrss.timer`:
```ini
[Unit]
Description=ABC News Just-in RSS Feed Creator

[Service]
Type=oneshot
ExecStart=/usr/local/bin/abcjustinrss -output /var/www/localhost/htdocs/rss/abcjustinrss.xml
User=apache
Group=apache

[Timer]
OnCalendar=*-*-* *:~15:00

[Install]
WantedBy=timers.target
```
2. Reload systemd and start the service:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now abcjustinrss.timer
   ```

#### systemd (as user)
1. Create a systemd service file at `$HOME/.config/systemd/user/abcjustinrss.timer`:
```ini
[Unit]
Description=ABC News Just-in RSS Feed Creator

[Service]
Type=oneshot
ExecStart=/usr/local/bin/abcjustinrss -output ~/public_html/rss/abcjustinrss.xml

[Timer]
OnCalendar=*-*-* *:~15:00

[Install]
WantedBy=timers.target
```
2. Reload systemd and start the service:
   ```bash
   systemctl --user enable --now abcjustinrss.timer
   ```

#### Apache VirtualHost Configuration
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

