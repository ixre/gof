/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: aes.ts
 * author: jarrysix (jarrysix@gmail.com)
 * date: 2024-10-28 21:11:21
 * description: AES加密解密
 * history:
 */

import CryptoJS from "crypto-js";

// 示例用法
let key = CryptoJS.enc.Utf8.parse(""); // 32 字节密钥
let iv = CryptoJS.enc.Utf8.parse(""); // 16 字节偏移量

//解密方法(hex)
function Decrypt(word) {
  const encryptedHexStr = CryptoJS.enc.Hex.parse(word); //指定HEx
  const srcs = CryptoJS.enc.Base64.stringify(encryptedHexStr);
  const decrypt = CryptoJS.AES.decrypt(srcs, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  });
  return decrypt.toString(CryptoJS.enc.Utf8);
}

//加密方法(hex)
function Encrypt(word) {
  const srcs = CryptoJS.enc.Utf8.parse(word);
  const encrypted = CryptoJS.AES.encrypt(srcs, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  });
  return encrypted.ciphertext.toString().toUpperCase();
}

export default {
  configure: (options: { key: string; iv: string }) => {
    if (options.key.length != 32) {
      throw new Error("key must be 32 characters");
    }
    if (options.iv.length != 16) {
      throw new Error("iv must be 16 characters");
    }
    key = CryptoJS.enc.Utf8.parse(options.key);
    iv = CryptoJS.enc.Utf8.parse(options.iv);
  },
  decrypt: Decrypt,
  encrypt: Encrypt,
};
