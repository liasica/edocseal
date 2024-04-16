#!/usr/bin/env bash

ENTPATH=./internal/ent

# 清理并生成ent
for path in "$ENTPATH"/*; do
    [ "$path" = "$ENTPATH/schema" ] && continue
    [ "$path" = "$ENTPATH/template" ] && continue
    [ "$path" = "$ENTPATH/connect.go" ] && continue
    rm -rf "path"
done

# 重新生成ent
go run -mod=mod entgo.io/ent/cmd/ent generate "$ENTPATH/schema" --feature schema/snapshot,namedges,sql/modifier,sql/execquery,sql/upsert --template "$ENTPATH/template/upsert.tmpl"
