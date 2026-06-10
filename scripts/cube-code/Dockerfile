FROM cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest

RUN mkdir -p /root/.pip && \
    echo "[global]" > /root/.pip/pip.conf && \
    echo "index-url = https://repo.huaweicloud.com/repository/pypi/simple" >> /root/.pip/pip.conf && \
    echo "trusted-host = repo.huaweicloud.com" >> /root/.pip/pip.conf

CMD ["python3", "/srv/app.py"]
# docker build -t 192.168.1.100:5000/sandbox/cube-code-sandbox:v1 .
# docker push 192.168.1.100:5000/sandbox/cube-code-sandbox:v1


# docker build -t registry.i-mall.top/sandbox/cube-code-sandbox:v1 .
# docker push registry.i-mall.top/sandbox/cube-code-sandbox:v1