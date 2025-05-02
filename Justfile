build:
  cd ./apps/backend && go build -o kyob

generate_types:
  protoc -I=./proto --go_out=./apps/backend/src/types --go_opt=paths=source_relative ./proto/types.proto
