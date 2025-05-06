build:
  cd ./apps/backend && go build -o kyob

generate_types:
  protoc -I=./proto --go_out=./apps/backend/types --go_opt=paths=source_relative ./proto/types.proto

generate_schemas:
  cd ./apps/backend && sqlc generate

migrate:
  cd ./apps/backend && dbmate -d ./db/migrations up

rollback:
  cd ./apps/backend && dbmate -d ./db/migrations down

create_migration name:
  cd ./apps/backend && dbmate -d ./db/migrations new {{name}}
