

// 初始化密钥和偏移量
var key = encrypt.Key.fromUtf8('your 32 byte key');
var iv = encrypt.IV.fromUtf8('your 16 byte iv');

// 解密方法
String aesDecrypt(String encrypted) {
  final encryptedBytes = encrypt.Encrypted(hexToBytes(encrypted));
  final encrypter = encrypt.Encrypter(
      encrypt.AES(key, mode: encrypt.AESMode.cbc, padding: "PKCS7"));
  final decrypted = encrypter.decrypt(encryptedBytes, iv: iv);
  return decrypted;
}

Uint8List hexToBytes(String encrypted) {
  final length = encrypted.length;
  final bytes = <int>[];
  for (var i = 0; i < length; i += 2) {
    final byte = int.parse(encrypted.substring(i, i + 2), radix: 16);
    bytes.add(byte);
  }
  return Uint8List.fromList(bytes);
}

// 加密方法
String aesEncrypt(String plainText) {
  // import 'package:encrypt/encrypt.dart' as encrypt;
  final encrypter = encrypt.Encrypter(
      encrypt.AES(key, mode: encrypt.AESMode.cbc, padding: "PKCS7"));
  final encrypted = encrypter.encrypt(plainText, iv: iv);
  return bytesToHex(encrypted.bytes).toUpperCase(); // 使用自定义的 bytesToHex 函数
}

// 将字节数组转换为十六进制字符串
String bytesToHex(List<int> bytes) {
  return bytes.map((byte) => byte.toRadixString(16).padLeft(2, '0')).join();
}

// 配置密钥和偏移量
void aesConfigure(String newKey, String newIv) {
  if (newKey.length != 32) {
    throw Exception("key length must be 32");
  }
  if (newIv.length != 16) {
    throw Exception("iv length must be 16");
  }
  key = encrypt.Key.fromUtf8(newKey);
  iv = encrypt.IV.fromUtf8(newIv);
}