#!/bin/bash

script_dir=$(
  cd $(dirname $0)
  pwd
)                                  # 脚本路径
project_dir=$(dirname $script_dir) # 项目路径

protos=$(find ${project_dir}/cmd/example -type f -name '*.proto')

for proto in $protos; do
  protoc \
    -I ${project_dir} \
    -I ~/go/bin/include \
    --go_out ${project_dir} \
    --go_opt paths=source_relative \
    --ginx_out ${project_dir} \
    --ginx_opt paths=source_relative \
    --ginx_opt rpc_mode=official \
    --ginx_opt use_encoding=true \
    $proto
done

