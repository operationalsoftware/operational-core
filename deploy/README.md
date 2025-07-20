## Deployment Procedure

1. Run a new cloud server using the `server-setup.sh` script in user data.
2. (Optional) Restore database backup once the server is up.
    - `curl -L -o batten_allen.dump 'https://presigned-url...'`
    - `pg_restore -U postgres -d batten_allen batten_allen.dump`
4. Run the deploy/deploy.sh

## Some commands for Database backup service:

```
sudo systemctl daemon-reload
sudo systemctl enable --now db-backup.timer
sudo systemctl list-timers
systemctl status db-backup.timer
journalctl -u db-backup.service
```
