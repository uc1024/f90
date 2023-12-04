#!/bin/bash

script_dir=$(
  cd $(dirname $0)
  pwd
)                                  # 脚本路径
project_dir=$(dirname $script_dir) # 项目路径

protos=$(find ${project_dir}/lib/examplepb -type f -name '*.proto')

for proto in $protos; do
  protoc \
    -I ${project_dir} \
    -I ~/go/bin/include \
    --go_out ${project_dir} \
    --go_opt paths=source_relative \
    $proto
done

