# baxup

Performs backup on target docker compose container and its contents, storing them as a `*.bak.tar.gz` compression file to save on disk space. Optionally, the newly compressed `.tar.gz` file can be copied securely to a remote machine when passed flags to do so. Baxup can be called by crontabs to faciliate both local and remote backups on a nightly schedule.

Although Baxup is primarily designed to perform backups on docker containers defined via `docker compose` with volume mounts in the working directory of the `docker-compose.yml` file, Baxup can also perform the same style of local and remote backups on regular file directories as well. Unfortunately Baxup cannot support backups for docker volume mounts other than `.`. Passing the `-docker` flag at command runtime will allow Baxup to perform `docker compose down` and `docker compose up -d` commands from your docker container's working directories before and after performing backups to ensure data consistency and safety.

Crucially, prior to shutting down a docker compose service and performing a compressive backup, Baxup will also collect the current image digests from the docker services and store them alongside the `docker-compose.yml` file such that the current image digests are always tracked with each backup to help facilitate restoration in a docker container failure/emergency. 

The newly backed-up file, including the aforementioned image digests, can optionally be transferred to a remote machine securely via Rsync using the `-remote-send` flag. Baxup forces SSH transport to ensure security in transit and forces checksum validations to ensure no data is corrupted during the remote transfer process. 


### Example use-cases

Compress a copy of a target directories data, storing it in the defined backup directory
```
## Compresses /$ROOTPATH/$TARGETNAME/ to defined backup directory
./baxup -target=foo
```

Stop target docker container, perform a backup of its current image digests and stored data, and restart it
```
## Stops $TARGETNAME docker container, collects image digests, compresses data to store in backup dir, and restarts container in background
./baxup -target=foobar -docker
```

Backup docker container, storing copy of the backup both Locally and on a Remote backup server
Note that in order for Baxup to be used with Crontabs, an SSH key must be utilized! Otherwise, remote password can be passed at runtime.
```
## Backs up docker container's data and copies it to remote machine
./baxup -target-name=foobar -docker -remote-send=true -remote-user=admin -remote-host=192.168.0.1
```

### Crontab usecase
```
> crontab -e

. . . 
# m h  dom mon dow   command

## Perform local backup on docker container every night at 1:00
0 1 * * * /home/agriffin/go-tools/baxup/baxup -target=bar -docker
```
