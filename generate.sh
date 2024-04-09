#!/usr/bin/env bash

ENTPATH=./internal/ent
# 清理并生成ent
for dir in "$ENTPATH"/*; do
    [ "$dir" = "$ENTPATH/schema" ] && continue
    rm -rf "$dir"
done
go run -mod=mod entgo.io/ent/cmd/ent generate "$ENTPATH/schema" --feature schema/snapshot,namedges,sql/modifier,sql/execquery,sql/upsert
