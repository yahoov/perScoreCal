set -e
source .env
# sudo -
# ssh -i private_key -p 22 user@host "uptime"
# psql postgres -c "CREATE DATABASE mytemplate1 WITH ENCODING 'UTF8' TEMPLATE template0"
# sudo printf "CREATE USER $psqluser WITH PASSWORD '$psqlpass';\nCREATE DATABASE $psqldb WITH OWNER $psqluser;" > cartaro.sql

# psqluser="koko28"   # Database username
# psqlpass="pass123"  # Database password
# psqldb="kokodb28"   # Database name

sudo printf "CREATE USER $DEV_USERNAME WITH PASSWORD '$DEV_PASSWORD';\nCREATE DATABASE $DEV_DBNAME WITH OWNER $DEV_USERNAME;" > setupPostgres.sql
sudo -u postgres psql -f setupPostgres.sql
