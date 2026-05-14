- Q: Listen: listen tcp 169.254.254.53:53: bind: cannot assign requested address
- A: 
```shell
# 创建虚拟接口
sudo ip link add cube-dns0 type dummy
# 分配 IP 地址
sudo ip addr add 169.254.254.53/32 dev cube-dns0
# 启用接口
sudo ip link set cube-dns0 up
# 配置 systemd-resolved 使用该接口解析 cube.app 域名
sudo resolvectl dns cube-dns0 169.254.254.53
sudo resolvectl domain cube-dns0 ~cube.app
# 验证配置
resolvectl status cube-dns0
# 输出结果
#Link 5 (cube-dns0)
#    Current Scopes: DNS
#         Protocols: -DefaultRoute -LLMNR -mDNS -DNSOverTLS DNSSEC=no/unsupported
#Current DNS Server: 169.254.254.53
#       DNS Servers: 169.254.254.53
#        DNS Domain: ~cube.app

# 重启 cube-proxy-coredns 容器
```