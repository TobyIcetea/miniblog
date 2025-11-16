# Copyright 2025 TobyIcetea <x2406862525@163.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/onexstack/miniblog. The professional
# version of this repository is https://github.com/onexstack/onex.

# 定义 Header
HEADER='{"alg":"HS256","typ":"JWT"}'

# 定义 Payload
PAYLOAD='{"sub":"1234567890","name":"John Doe","iat":1516239022}'

# 定义 Secret（用于签名）
SECRET="Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5"

# 1. Base64 编码 Header
HEADER_BASE64=$(echo -n "$HEADER" | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')

# 2. Base64 编码 Payload
PAYLOAD_BASE64=$(echo -n "$PAYLOAD" | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')

# 3. 拼接 Header 和 Payload 为签名数据
SIGNING_INPUT="$HEADER_BASE64.$PAYLOAD_BASE64"

# 4. 使用 HMAC SHA256 算法生成签名
SIGNATURE=$(echo -n "$SIGNING_INPUT" | openssl dgst -sha256 -hmac "$SECRET" -binary | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')

# 5. 拼接最终的 JWT Token
JWT="${SIGNING_INPUT}.${SIGNATURE}"

# 输出 JWT Token
echo "Generated JWT Token"
echo "${JWT}"
