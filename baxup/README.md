# baxup

Performs backup on target docker compose container and its contents, storing them as a `*.bak.tar.gz` compression file to save on disk space. Optionally, the newly compressed `.tar.gz` file can be copied securely to a remote machine when passed flags to do so. So long as SSH keys are set up between the local machine running baxup and the target remote machine, (or only local backups used) Baxup can be called by crontabs to faciliate both local and remote backups on a nightly schedule.

Baxup is primarily designed to perform backups on docker containers defined via `docker compose` *with volume mounts in the working directory of the `docker-compose.yml` file*, Baxup can also perform the same style of local and remote backups on regular everyday file directories as well. Unfortunately Baxup cannot support backups for docker volume mounts other than `.` without some extra effort. Passing the `-docker` flag at command runtime will allow Baxup to perform `docker compose down` and `docker compose up -d` commands from your docker container's working directories before and after performing backups to ensure data consistency and safety.

Crucially, prior to shutting down a docker compose service and performing a compressive backup, Baxup will also collect the current docker image digests from the docker services and store them alongside the `docker-compose.yml` file such that the current image digests are *always* tracked with each backup at shutdown to help facilitate restoration in a docker container failure/emergency, especially those regarding updates on images that result in DB errors.

The newly compressed file, including the aforementioned image digests, can optionally be transferred to a remote machine securely and reliably via Rsync using the `-remote-send` flag. Baxup forces SSH transport to ensure security in transit and forces checksum validations during receipt to ensure no data is corrupted or altered during the remote transfer process. 

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

Backup docker container, storing copy of the backup and img digests both Locally and on a Remote backup server
Note that in order for Baxup to be used with Crontabs, an SSH key must be utilized! Otherwise, remote password can be passed at runtime.
```
## Backs up docker container's data and copies it to remote machine
./baxup -target-name=foobar -docker -remote-send=true -remote-user=admin -remote-host=192.168.0.1
```

### Crontab usecase
Again, please note that crontab backups will require SSH keys between the local and target/remote machine (unless you wanna be around to enter the remote password at runtime in the middle of the night lol)
These work great for nightly docker backups and copy critical docker container data to a remote machine for easy restoration in the event of an emergency, even being self contained enough to be rebuilt on the remote machine in the event the local one no longer exists.
```
> crontab -e

. . . 
# m h  dom mon dow   command

## Perform local backup on docker container every night at 1:00
0 1 * * * /home/agriffin/go-tools/baxup/baxup -target=bar -docker
```
