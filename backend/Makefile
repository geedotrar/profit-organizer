export POSTGRES_URL="postgres://postgres:password@localhost:5440/laba_db?sslmode=disable"

migrate-create:
	@ migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	@ migrate -database ${POSTGRES_URL} -path migrations up

migrate-down:
	@ migrate -database ${POSTGRES_URL} -path migrations down
	