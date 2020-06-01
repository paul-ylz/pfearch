db.create-migration:
	migrate create -ext sql -dir migrations -seq -digits 4 $(MIGRATION_NAME)