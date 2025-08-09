migrate-db:
	@echo "📜 Running migrations from ./sql/schema.sql..."
	
	@cat ./sql/schema.pg.sql | docker exec -i apartment-db \
		env PGPASSWORD=postgres \
		psql -v ON_ERROR_STOP=1 -U postgres -d postgres

	@echo "✅ Migration complete!"
	
swagger:
	@swag init -g cmd/main.go -o ./docs

env:
	@chmod +x env.sh
	@./env.sh

env-quick:
	@chmod +x env.sh
	@./env.sh quick
