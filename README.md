# UMG
user management system

#### Backup database
```bash
# creat backup
docker exec -i user_management_db pg_dump -U user-umg umg > backup.sql

# restore backup
docker exec -i user_management_db psql -U user-umg umg < backup.sql
```

### API Document

API document is available [here](https://github.com/boof/ptrack/backend/umg-docs/-/blob/master/swagger.yaml)
