build:
  cd ./apps/backend && go build -o kyob

generate_types:
  protoc -I=./proto --go_out=./apps/backend/src/types --go_opt=paths=source_relative ./proto/types.proto

generate_schemas:
  cd ./apps/backend && sqlc generate

migrate:
  cd ./apps/backend && dbmate -d ./migrations up

rollback:
  cd ./apps/backend && dbmate -d ./migrations down

create_migration name:
  cd ./apps/backend && dbmate -d ./migrations new {{name}}
