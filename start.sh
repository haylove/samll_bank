set -e
echo "start migration"
/app/migrate --path /app/migration/ -database "$DB_SOURCE" -verbose up

echo "start app"
exec "$@"
