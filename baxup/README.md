# baxup

Performs backup on target directory and its contents, storing them as a `*.bak.tar.gz` compression file to save on space. Optionally, the newly compressed `.tar.gz` file can be copied securely to a remote machine when passed flags to do so. Baxup can be called by crontabs to faciliate both local and remote backups on a nightly schedule.

Though Baxup can be used to back up a copy of any directory, it is primarily designed to perform backups on docker containers that are defined via `docker compose` with volume mounts in the working directory of the `docker-compose.yml` file. Unfortunately Baxup cannot support backups for docker volume mounts other than `.`. Passing the `-docker` flag at command runtime will allow Baxup to perform `docker compose down` and `docker compose up -d` commands from your docker container's working directories before and after performing backups to ensure data consistency and safety. 

The newly backed-up file can optionally be transferred to a remote machine securely via Rsync using the `-remote-send` flag. Baxup enforces SSH transport to ensure security and forces checksum validations to ensure no data is corrupted during transfer. 


### Example use-cases

Compress a copy of a target directories data, storing it in the defined backup directory
```
## Compresses /$ROOTPATH/$TARGETNAME/ to defined backup directory
./baxup -target-name=foo
```

Stop target docker container, perform a backup of it's data, and restart the it
```
## Stops $TARGETNAME docker container, compresses data to store in backup dir, restart container in background
./baxup -target-name=vaultwarden -docker
```

Backup docker container, storing copy of the backup both Locally and on a Remote backup server
```
## Backs up docker container's data and copies it to remote machine
./baxup -target-name=vaultwarden -docker -remote-send=true -remote-user=admin -remote-host=192.168.0.1
```

### Crontab usecase
```
> crontab -e

. . . 
# m h  dom mon dow   command

## Perform local backup on docker container every night at 1:00
0 1 * * * /home/agriffin/go-tools/baxup/baxup -target-name=bar -docker
```
