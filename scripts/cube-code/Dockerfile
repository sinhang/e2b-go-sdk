FROM cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest

USER root

RUN mkdir -p /root/.pip && \
    printf "[global]\nindex-url=https://repo.huaweicloud.com/repository/pypi/simple\ntrusted-host=repo.huaweicloud.com\n" \
    > /root/.pip/pip.conf

RUN pip install --no-cache-dir \
    numpy \
    pandas \
    matplotlib \
    akshare

#CMD ["python3", "/srv/app.py"]
# docker build -t 192.168.1.100:5000/sandbox/cube-code-sandbox:v1 .
# docker push 192.168.1.100:5000/sandbox/cube-code-sandbox:v1


# docker build -t registry.i-mall.top/sandbox/cube-code-sandbox:v1 .
# docker push registry.i-mall.top/sandbox/cube-code-sandbox:v1