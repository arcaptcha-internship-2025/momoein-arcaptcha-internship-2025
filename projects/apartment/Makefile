migrate-dev-db:
	@echo "📜 Running migrations from ./sql/schema.sql..."
	
	@cat ./sql/schema.sql | docker exec -i apartment-dev-db \
		env PGPASSWORD=postgres \
		psql -v ON_ERROR_STOP=1 -U postgres -d postgres

	@echo "✅ Migration complete!"
	