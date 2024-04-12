FROM --platform=linux/amd64 python:latest

RUN mkdir -p /app \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apt-get clean && apt-get autoclean \
    && rm -rf /var/cache/apt/archives/* \
    && rm -rf /var/lib/apt/lists/*

#RUN <<EOF cat >> ~/.pip/pip.conf
#[global]
#index-url = http://mirrors.aliyun.com/pypi/simple/
#
#[install]
#trusted-host=mirrors.aliyun.com
#EOF

COPY ./build/release/edocseal /app/
COPY ./signer.py /app/

RUN pip install --upgrade pip \
    && pip install --no-cache-dir 'urllib3<2.0' \
    && pip install --no-cache-dir pyHanko \
    && pip install --no-cache-dir 'pyHanko[pkcs11,image-support,opentype,xmp,opentype,image-support]' \
    && chmod +x /app/signer.py

WORKDIR /app

ENTRYPOINT ["/app/edocseal", "server", "-c", "/app/config/config.yaml"]
