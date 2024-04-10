import json
import sys
from typing import Tuple, List, Any, TypeVar, Callable, Type, cast

from pyhanko import stamp
from pyhanko.pdf_utils import images
from pyhanko.pdf_utils.incremental_writer import IncrementalPdfFileWriter
from pyhanko.sign import signers, fields

T = TypeVar('T')


def list_to_tuple(items: List[int]) -> Tuple[int, int, int, int]:
    assert len(items) == 4
    return items[0], items[1], items[2], items[3]


def from_str(x: Any) -> str:
    assert isinstance(x, str)
    return x


def from_list(f: Callable[[Any], T], x: Any) -> List[T]:
    assert isinstance(x, list)
    return [f(y) for y in x]


def from_int(x: Any) -> int:
    assert isinstance(x, int) and not isinstance(x, bool)
    return x


def to_class(c: Type[T], x: Any) -> dict:
    assert isinstance(x, c)
    return cast(Any, x).to_dict()


class SignatureField:
    signature: str
    key: str
    cert: str
    rect: List[int]

    def __init__(self, signature: str, key: str, cert: str, rect: List[int]) -> None:
        """
        签约人

        :param signature:
            签名图片
        :param key:
            私钥
        :param cert:
            证书
        :param rect:
            签字区域
        """
        self.signature = signature
        self.key = key
        self.cert = cert
        self.rect = rect

    @staticmethod
    def from_dict(obj: Any) -> 'SignatureField':
        assert isinstance(obj, dict)
        signature = from_str(obj.get('signature'))
        key = from_str(obj.get('key'))
        cert = from_str(obj.get('cert'))
        rect = from_list(from_int, obj.get('rect'))
        return SignatureField(signature, key, cert, rect)

    def to_dict(self) -> dict:
        result: dict = {
            'signature': from_str(self.signature),
            'key': from_str(self.key),
            'cert': from_str(self.cert),
            'rect': from_list(from_int, self.rect),
        }
        return result


class SignatureConfig:
    in_file: str
    out_file: str
    signatures: List[SignatureField]

    def __init__(self, in_file: str, out_file: str, signatures: List[SignatureField]) -> None:
        """
        签约配置

        :param in_file:
            待签约文档
        :param out_file:
            签约保存文档
        :param signatures:
            签约人列表（目前只支持两人）
        """
        self.in_file = in_file
        self.out_file = out_file
        self.signatures = signatures

    @staticmethod
    def from_dict(obj: Any) -> 'SignatureConfig':
        assert isinstance(obj, dict)
        in_file = from_str(obj.get('in_file'))
        out_file = from_str(obj.get('out_file'))
        signatures = from_list(SignatureField.from_dict, obj.get('signatures'))
        return SignatureConfig(in_file, out_file, signatures)

    def to_dict(self) -> dict:
        result: dict = {
            'in_file': from_str(self.in_file),
            'out_file': from_str(self.out_file),
            'signatures': from_list(lambda x: to_class(SignatureField, x), self.signatures)
        }
        return result


def signature_config_from_dict(s: Any) -> SignatureConfig:
    return SignatureConfig.from_dict(s)


def signature_config_to_dict(x: SignatureConfig) -> Any:
    return to_class(SignatureConfig, x)


def sign_double():
    s1 = cfg.signatures[0]
    s2 = cfg.signatures[1]

    signer1 = signers.SimpleSigner.load(s1.key, s1.cert)
    signer2 = signers.SimpleSigner.load(s2.key, s2.cert)

    with open(cfg.in_file, 'rb') as inf:
        w = IncrementalPdfFileWriter(inf, strict=False)
        fields.append_signature_field(
            w, sig_field_spec=fields.SigFieldSpec(
                'Signature1', box=list_to_tuple(s1.rect), on_page=-1
            )
        )
        fields.append_signature_field(
            w, sig_field_spec=fields.SigFieldSpec(
                'Signature2', box=list_to_tuple(s2.rect), on_page=-1
            )
        )

        pdf_signer_1 = signers.PdfSigner(
            signers.PdfSignatureMetadata(
                field_name='Signature1',
                certify=True,
                docmdp_permissions=fields.MDPPerm.FILL_FORMS,
            ),
            signer=signer1,
            stamp_style=stamp.TextStampStyle(
                stamp_text='',
                background=images.PdfImage(s1.signature),
                border_width=0,
                background_opacity=0.7,
            ),
        )
        out = pdf_signer_1.sign_pdf(w)

        w = IncrementalPdfFileWriter(out, strict=False)
        pdf_signer_2 = signers.PdfSigner(
            signers.PdfSignatureMetadata(field_name='Signature2'),
            signer=signer2,
            stamp_style=stamp.TextStampStyle(
                stamp_text='',
                background=images.PdfImage(s2.signature),
                border_width=0,
                background_opacity=0.7,
            ),
        )
        out = pdf_signer_2.sign_pdf(w)

        with open(cfg.out_file, 'wb') as outf:
            buf = out.getbuffer()
            outf.write(buf)
            buf.release()
            outf.close()


if __name__ == '__main__':
    """
    调用签名， 使用示例:
    python3 singer.py "{\"in_file\":\"test/document.pdf\",\"out_file\":\"test/document-signed.pdf\",\"signatures\":[{\"signature\":\"test/seal.png\",\"key\":\"test/key1.pem\",\"cert\":\"test/cert1.pem\",\"rect\":[70,320,206,456]},{\"signature\":\"test/signature.png\",\"key\":\"test/key2.pem\",\"cert\":\"test/cert2.pem\",\"rect\":[242,343,317,419]}]}"
    """
    args = sys.argv[1:]
    if len(args) == 1:
        cfg = signature_config_from_dict(json.loads(args[0]))
        if cfg is None or len(cfg.signatures) != 2:
            print('json配置读取失败')
            sys.exit(1)
        sign_double()
    else:
        print('配置获取失败')
        sys.exit(1)
