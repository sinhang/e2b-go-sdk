FROM deltadc/cubesandbox-base:v2

# 安装你自己需要的工具链
RUN apt-get update \
    && apt-get install -y --no-install-recommends python3 python3-pip \
    && rm -rf /var/lib/apt/lists/*


RUN mkdir -p /root/.pip && \
    echo "[global]" > /root/.pip/pip.conf && \
    echo "index-url = https://repo.huaweicloud.com/repository/pypi/simple" >> /root/.pip/pip.conf && \
    echo "trusted-host = repo.huaweicloud.com" >> /root/.pip/pip.conf

RUN pip install --no-cache-dir pandas matplotlib numpy
# docker build -t registry.i-mall.top/sandbox/cube-base-code-sandbox:v1 .
# docker push registry.i-mall.top/sandbox/cube-base-code-sandbox:v1

# 如果你的应用需要作为前台进程运行，在这里设置 CMD 即可。
# envd 仍会作为后台进程持续运行。
# CMD ["python3", "/srv/app.py"]