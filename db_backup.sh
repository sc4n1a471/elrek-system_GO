#!/bin/bash

DATETIME=$(date +%F_%H-%M-%S)

BACKUP_FILE="<dir>/<filename>_backup_$DATETIME.sql"

mysqldump --defaults-file=<.my.cnf file location> -u <username> <db_name> > $BACKUP_FILE

if [ $? -eq 0 ]; then
  echo "Backup successfully created at $BACKUP_FILE"
else
  echo "Error creating backup" >&2
  exit 1
fi

# Delete backups older than 14 days
find <backup_dir> -name "*_backup_*.sql" -type f -mtime +14 -exec rm -f {} \;
echo "Old backups deleted"