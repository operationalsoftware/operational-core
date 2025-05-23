## Deployment Procedure

1. Run a new cloud server using the `server-setup.sh` script in user data.
2. (Optional) Restore database backup once the server is up.
    - Use `scp` command to copy the .dump file to server at /home/debian.
    - Update permissions of dump file using chown to change ownership to `app` user.
    - Switch to the app user using the su command: su app (password `app`)
    - Run the following pg_restore command for restoring database backup:
        * pg_restore -U postgres -d batten_allen backup_file.dump
3. Update the *.env file for the corresponding branch.
4. Run the deploy/deploy.sh

## Some commands for Database backup service:

```
sudo systemctl daemon-reload
sudo systemctl enable --now db-backup.timer
sudo systemctl list-timers
systemctl status db-backup.timer
journalctl -u db-backup.service
```