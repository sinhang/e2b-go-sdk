FROM cube-sandbox-cn.tencentcloudcr.com/cube-sandbox/sandbox-code:latest

USER root

# -----------------------------------------------------------------------------
# pip 国内源
# -----------------------------------------------------------------------------
RUN mkdir -p /root/.pip && \
    printf "[global]\nindex-url=https://repo.huaweicloud.com/repository/pypi/simple\ntrusted-host=repo.huaweicloud.com\n" \
    > /root/.pip/pip.conf

# -----------------------------------------------------------------------------
# 安装系统依赖
# -----------------------------------------------------------------------------
RUN apt-get update && \
    apt-get install -y \
        curl \
        wget \
        git \
        ca-certificates \
        gnupg \
        software-properties-common \
        libnss3 \
        libatk-bridge2.0-0 \
        libdrm2 \
        libxkbcommon0 \
        libgbm1 \
        libasound2 \
        libxshmfence1 \
        libxcomposite1 \
        libxdamage1 \
        libxrandr2 \
        libxcursor1 \
        libxinerama1 \
        libxi6 \
        libxtst6 \
        libcups2 \
        libpango-1.0-0 \
        libpangocairo-1.0-0 \
        libcairo2 \
        libgdk-pixbuf2.0-0 \
        libatspi2.0-0 \
        libgtk-3-0 \
        fonts-liberation \
        fonts-noto-cjk && \
    rm -rf /var/lib/apt/lists/*

# -----------------------------------------------------------------------------
# Python 包
# -----------------------------------------------------------------------------
RUN pip install --no-cache-dir \
    numpy \
    pandas \
    matplotlib \
    akshare \
    "markitdown[pptx]" \
    playwright

# -----------------------------------------------------------------------------
# Node.js 20
# -----------------------------------------------------------------------------
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get update && \
    apt-get install -y nodejs && \
    rm -rf /var/lib/apt/lists/*

# -----------------------------------------------------------------------------
# npm / pnpm 国内源
# -----------------------------------------------------------------------------
RUN npm config set registry https://registry.npmmirror.com && \
    npm install -g pnpm playwright

RUN pnpm config set registry https://registry.npmmirror.com

# -----------------------------------------------------------------------------
# Playwright 下载国内镜像
# -----------------------------------------------------------------------------
RUN npx playwright install chromium

# Python Playwright 浏览器
RUN python3 -m playwright install chromium

# -----------------------------------------------------------------------------
# 验证安装
# -----------------------------------------------------------------------------
RUN node --version && \
    npm --version && \
    pnpm --version && \
    python3 --version && \
    pip --version

# -----------------------------------------------------------------------------
# 清理缓存
# -----------------------------------------------------------------------------
RUN npm cache clean --force && \
    pnpm store prune && \
    apt-get clean && \
    rm -rf \
        /root/.cache \
        /var/cache/apt \
        /tmp/*

# -----------------------------------------------------------------------------
# 默认工作目录
# -----------------------------------------------------------------------------
WORKDIR /workspace

# docker build -t registry.i-mall.top/sandbox/cube-code-sandbox-skill:v1 .
# docker push registry.i-mall.top/sandbox/cube-code-sandbox-skill:v1