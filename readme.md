# README

## Overview
This application scrapes the "Just In" section of ABC News and generates an RSS feed. 

## Installation

### Prerequisites
- **Go** (if building from source)
- **Apache or Nginx** (for serving RSS files)
- **Docker** (optional, for containerized deployment)

### Build and Install

### Usage
#### Generate RSS Feed
```bash
abcjustinrss -output /var/www/localhost/htdocs/rss/feed.xml
```

### Deployment

#### rc.d (Cron Job)
Add a cron job to run the script periodically:
1. Edit the user's crontab:
   ```bash
   crontab -e
   ```
2. Add the following line:
   ```bash
   */15 * * * * /usr/local/bin/abcjustinrss -output ~/public_html/rss/feed.xml
   ```

#### systemd
1. Create a systemd service file at `/etc/systemd/system/abcjustinrss.service`:
```ini
[Unit]
Description=RSS Feed Creator
After=network.target

[Service]
ExecStart=/usr/local/bin/abcjustinrss -output /var/www/localhost/htdocs/rss/feed.xml
User=apache
Group=apache
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
2. Reload systemd and start the service:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now abcjustinrss
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

