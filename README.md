
### 常用指令
```shell
# 136
pyhanko sign addfields --field -1/70,320,206,456/AUR-SIGN input.pdf input-s.pdf
pyhanko sign addsig --field AUR-SIGN pemder --key privkey.pem --cert cert.pem --no-pass input-s.pdf output.pdf
```

### 参考文档
- [pyHanko](https://github.com/MatthiasValvekens/pyHanko)
- [PDF Explained （译作《PDF 解析》）](https://github.com/zxyle/PDF-Explained/blob/master/chapter1.md)
- [使用Go语言生成自签CA证书](https://foreverzmyer.hashnode.dev/go-cert)
- [使用Golang X509签发证书及构建CA架构](https://blog.yeziruo.com/archives/148.html)