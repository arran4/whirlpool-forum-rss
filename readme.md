# README

## Overview
This application grabs the forum topics from Whirlpool forum's newthreads / popular_views summaries and converts them
into an RSS feed for consumption.

## Installation

### Prerequisites
- **Go** (if building from source)

#### Optional
- **Apache or Nginx** (for serving RSS files)

### Build and Install

#### Install only (System level)
Grab the latest binary here: https://github.com/arran4/whirlpool-forum-rss/releases/

#### Install and build as user (User)
Install go 1.23+

Run `go install`:
```bash
go install github.com/arran4/whirlpool-forum-rss/cmd/whirlpoolforumrss@latest
```
This installs to `$HOME/go/bin` (typically; check with `go env`).

### Usage
#### CLI Mode
Generate RSS Feed:
```bash
whirlpoolforumrss -output /var/www/localhost/htdocs/rss/whirlpoolforumnewthreadsrss.xml
```

#### CGI Mode
1. Place `whirlpoolforumrss-cgi` in your server's CGI directory (e.g., `/var/www/htdocs/cgi-bin/whirlpoolforumrss-cgi`).
2. Ensure it is executable:
   ```bash
   chmod +x /var/www/htdocs/cgi-bin/whirlpoolforumrss-cgi
   ```
3. Access it via URL (e.g., `http://example.com/cgi-bin/whirlpoolforumrss-cgi`).

### Deployment

#### rc.d (Cron Job system level)
Add a cron job to run the script periodically:
1. Edit the root crontab:
   ```bash
   sudo crontab -e
   ```
2. Add the following line:
   ```bash
   */15 * * * * /usr/local/bin/whirlpoolforumrss -output /var/www/localhost/htdocs/rss/whirlpoolforumnewthreadsrss.xml
   ```

#### rc.d (Cron Job user level)
Add a cron job to run the script periodically:
1. Edit the user's crontab:
   ```bash
   crontab -e
   ```
2. Add the following line:
   ```bash
   */15 * * * * ~/go/bin/whirlpoolforumrss -output ~/public_html/rss/whirlpoolforumnewthreadsrss.xml
   ```

#### systemd (as root)
1. Create a systemd service file at `/etc/systemd/system/whirlpoolforumrss.service`:
```ini
[Unit]
Description=Whirlpool Forum RSS Feed Creator

[Service]
Type=oneshot
ExecStart=/usr/bin/whirlpoolforumrss -output /var/www/localhost/htdocs/rss/whirlpoolforumnewthreadsrss.xml
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
   sudo systemctl enable --now everyhour@whirlpoolforumrss.timer
   ```

#### systemd (as user)
1. Create a systemd service file at `$HOME/.config/systemd/user/whirlpoolforumrss.service`:
```ini
[Unit]
Description=Whirlpool Forum RSS Feed Creator

[Service]
Type=oneshot
ExecStart=%h/go/bin/whirlpoolforumrss -output %h/public_html/rss/whirlpoolforumnewthreadsrss.xml
```

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
   systemctl --user daemon-reload && systemctl --user enable --now everyhour@whirlpoolforumrss.timer
   ```

#### Apache VirtualHost Configuration
##### User

Refer to documentation for setting up public_html directories

##### Enjoy

http://localhost/~$USERNAME/rss/whirlpoolforumnewthreadsrss.xml

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

#### Note CGI action selection:

Use: `http://localhost/cgi-bin/whirlpoolforumrss-cgi?action=newthreads` or `http://localhost/cgi-bin/whirlpoolforumrss-cgi?action=popular_views` 
for the desired view.