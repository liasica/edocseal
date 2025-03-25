FROM python:latest

RUN mkdir -p /app \
    && apt update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends tzdata unzip \
    && curl -s https://api.github.com/repos/qpdf/qpdf/releases/latest \
       | grep "browser_download_url.*-bin-linux-x86_64.zip" \
       | cut -d : -f 2,3 \
       | tr -d \" \
       | wget --no-check-certificate -i - -O /usr/local/qpdf.zip \
    && unzip /usr/local/qpdf.zip -d /usr/local/ \
    && rm /usr/local/qpdf.zip \
    && rm -rf /etc/localtime \
    && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
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

RUN pip install --upgrade pip \
    && pip install --no-cache-dir 'urllib3<2.0' \
    && pip install --no-cache-dir pyHanko \
    && pip install --no-cache-dir 'pyHanko[pkcs11,image-support,opentype,xmp,opentype,image-support]'

COPY ./build/release/edocseal /app/
COPY ./signer.py /app/

RUN chmod +x /app/signer.py

WORKDIR /app

ENTRYPOINT ["/app/edocseal", "server", "-c", "/app/config/config.yaml"]
