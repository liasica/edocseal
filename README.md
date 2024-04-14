### 待办

- [ ] 使用队列处理签名任务，例如 [asynq](https://github.com/hibiken/asynq) 或 [nsq](https://github.com/nsqio/nsq)
- [ ] 合同24小时内有效，超时自动删除
- [ ] 删除失效的文件

### 生成模板

```shell
edocseal template templates/原始模板.pdf -t runtime/模板表单.pdf
```

### 编译 Base Docker

```shell
docker build -t liasica/edocseal ./docker/Dockerfile
```

### 签名

#### 依赖

```shell
pip install 'urllib3<2.0'
pip install pyHanko
pip install 'pyHanko[pkcs11,image-support,opentype,xmp,opentype,image-support]'
```

```shell
# 136
pyhanko sign addfields --field -1/70,320,206,456/AUR-SIGN --field -1/242,343,317,419/RIDER-SIGN input.pdf input-s.pdf
pyhanko sign addsig --field AUR-SIGN pemder --key privkey.pem --cert cert.pem --no-pass input-s.pdf output.pdf
pyhanko sign addsig --field AUR-SIGN pemder --key cakey.pem --cert ca.pem --no-pass input-s.pdf output.pdf
pyhanko sign addsig --field AUR-SIGN pemder --key key2.pem --cert cert2.pem --no-pass input-s.pdf output.pdf
pyhanko sign addsig --field RIDER-SIGN pemder --key privkey.pem --cert cert.pem --no-pass output.pdf output.pdf

pyhanko sign addsig --field -1/70,320,206,456/AUR-SIGN --field -1/242,343,317,419/RIDER-SIGN pemder --key cakey.pem --cert ca.pem --key privkey.pem --cert cert.pem --no-pass input.pdf output.pdf
```

```shell
# qpdf
qpdf form.pdf --json > x.json # .acroform.fields
qpdf form.pdf --update-from-json=x.json x.pdf
```

### 尺寸

- 公章尺寸: 125 × 125
- 签名尺寸: 75 × 75

### 签名流程

1. 获取用户手写签名保存为{sn}.png
2. 读取默认 pyhanko.yml 并将用户手写签名图片路径导入并保存为{sn}.yml
3. 使用 pyhanko 指定配置分别签名

### 参考文档

- [pyHanko](https://github.com/MatthiasValvekens/pyHanko)
- [PDF Explained （译作《PDF 解析》）](https://github.com/zxyle/PDF-Explained/blob/master/chapter1.md)
- [使用Go语言生成自签CA证书](https://foreverzmyer.hashnode.dev/go-cert)
- [使用Golang X509签发证书及构建CA架构](https://blog.yeziruo.com/archives/148.html)
