const puppeteer = require('puppeteer');
const fs = require('fs');

// Boss直聘使用的AES密钥（Base64编码）
const AES_SECRET_KEY = 'clRwXUJBK1VKK0k0IWFbbQ==';

async function getFp() {
  const browser = await puppeteer.launch({
    headless: 'new',
    // 如果遇到证书问题，可以添加以下参数
    args: ['--ignore-certificate-errors']
  });

  try {
    const page = await browser.newPage();

    // 访问Boss直聘页面加载指纹脚本
    await page.goto('https://www.zhipin.com/web/user/', {
      waitUntil: 'domcontentloaded',
      timeout: 30000
    });

    // 注入CryptoJS库
    await page.addScriptTag({
      url: 'https://static.zhipin.com/library/js/utils/crypto-js.min.js'
    });

    // 等待指纹脚本加载完成
    await page.waitForFunction(
      () => window.zpFingerPrint && typeof window.zpFingerPrint.create === 'function',
      { timeout: 10000 }
    );

    // 修改这部分代码来输出原始指纹和加密后的指纹
    const result = await page.evaluate((key) => {
      return new Promise(resolve => {
        window.zpFingerPrint.create(rawFp => {
          try {
            // AES-CBC加密
            const text = CryptoJS.enc.Utf8.parse(rawFp);
            const keyUtf8 = CryptoJS.enc.Utf8.parse(atob(key));
            const iv = CryptoJS.lib.WordArray.random(16);
            const encrypted = CryptoJS.AES.encrypt(text, keyUtf8, {
              iv: iv,
              mode: CryptoJS.mode.CBC,
              padding: CryptoJS.pad.Pkcs7
            });

            const full = CryptoJS.lib.WordArray.create()
              .concat(iv)
              .concat(encrypted.ciphertext);
            const b64 = CryptoJS.enc.Base64.stringify(full);
            const encodedFp = encodeURIComponent(b64);

            resolve({
              rawFingerprint: rawFp,
              encryptedFingerprint: b64,
              encodedFingerprint: encodedFp
            });
          } catch (err) {
            console.error('加密失败:', err);
            resolve({
              rawFingerprint: '',
              encryptedFingerprint: '',
              encodedFingerprint: ''
            });
          }
        });
      });
    }, AES_SECRET_KEY);

    // 格式化输出结果
    console.log('\n=== 指纹信息 ===');
    console.log('原始指纹:', result.rawFingerprint);
    console.log('\n加密指纹:', result.encryptedFingerprint);
    console.log('\n最终fp参数:', result.encodedFingerprint);
    console.log('===============\n');

    // 保存详细信息到文件
    const outputContent = JSON.stringify(result, null, 2);
    fs.writeFileSync('fp.txt', result.encodedFingerprint);

    return result.encodedFingerprint;

  } catch (err) {
    console.error('获取fp失败:', err);
    return '';
  } finally {
    await browser.close();
  }
}

// 执行并获取fp
getFp().catch(console.error);
