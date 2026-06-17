FROM cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest

USER root

# 复制 envd
COPY --from=ghcr.io/tencentcloud/cubesandbox-base:2026.16 \
     /usr/bin/envd /usr/bin/envd
COPY --from=ghcr.io/tencentcloud/cubesandbox-base:2026.16 \
     /usr/local/bin/cube-entrypoint.sh /usr/local/bin/cube-entrypoint.sh

RUN chmod +x /usr/bin/envd /usr/local/bin/cube-entrypoint.sh

# 配置国内源（同上面 1-3 步）
RUN mkdir -p /root/.pip && \
    echo "[global]" > /root/.pip/pip.conf && \
    echo "index-url = https://repo.huaweicloud.com/repository/pypi/simple" >> /root/.pip/pip.conf && \
    echo "trusted-host = repo.huaweicloud.com" >> /root/.pip/pip.conf

# 如果 apt 源存在则替换
RUN if [ -f /etc/apt/sources.list ]; then \
        sed -i 's/archive.ubuntu.com/mirrors.tencent.com/g' /etc/apt/sources.list 2>/dev/null || true; \
        sed -i 's/security.ubuntu.com/mirrors.tencent.com/g' /etc/apt/sources.list 2>/dev/null || true; \
    fi

RUN pip install --no-cache-dir pandas matplotlib numpy akshare

# 保持原镜像的 ENTRYPOINT/CMD
# 但是添加一个前置启动脚本
COPY --chmod=755 <<'EOF' /docker-entrypoint-init.d/start-envd.sh
#!/bin/bash
/usr/bin/envd -port 49983 > /var/log/envd.log 2>&1 &
EOF


#CMD ["python3", "/srv/app.py"]
# docker build -t 192.168.1.100:5000/sandbox/cube-code-sandbox:v1 .
# docker push 192.168.1.100:5000/sandbox/cube-code-sandbox:v1


# docker build -t registry.i-mall.top/sandbox/cube-code-sandbox:v3 .
# docker push registry.i-mall.top/sandbox/cube-code-sandbox:v3