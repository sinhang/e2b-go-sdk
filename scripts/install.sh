#!/bin/bash

docker compose up -d

docker exec -it cube-sandbox bash

apt update && apt install python3 curl wget qemu-system-x86 qemu-utils ripgrep openssh-client make iproute2 python3-pip python3.12-venv -y
./prepare_image.sh
# /home/CubeSandbox/dev-env/.workdir/
# 新开终端
./run_vm.sh
ss -lntp
# 新开终端
./login.sh
curl -sL https://cnb.cool/CubeSandbox/CubeSandbox/-/git/raw/master/deploy/one-click/online-install.sh | CUBE_PVM_ENABLE=1 MIRROR=cn bash

# [run_vm][INFO]   SSH        : ssh -p 10022 opencloudos@127.0.0.1
# [run_vm][INFO]   Cube API   : http://127.0.0.1:13000 -> guest:3000
# [run_vm][INFO]   CubeProxy  : http://127.0.0.1:11080 -> guest:80
# [run_vm][INFO]   CubeProxy  : https://127.0.0.1:11443 -> guest:443
# [run_vm][INFO]   WebUI      : http://127.0.0.1:12088 -> guest:12088
