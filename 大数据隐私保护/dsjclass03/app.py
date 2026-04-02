import random
from phe import paillier

if __name__ == '__main__':
    # 使用 phe 库的 Paillier
    print("\n=== 使用 phe 库的 Paillier ===")
    public_key, private_key = paillier.generate_paillier_keypair()

    # 基本加解密测试
    message = 15
    encrypted_phe = public_key.encrypt(message)
    encrypted_phe = public_key.encrypt(message)
    print(f"原始消息: {message}")
    print(f"加密后: {encrypted_phe.ciphertext()}")
    decrypted_phe = private_key.decrypt(encrypted_phe)
    print(f"解密后: {decrypted_phe}")

    # 验证加法同态性
    print("\n=== 验证加法同态性 ===")
    m1, m2 = 10, 20
    c1 = public_key.encrypt(m1)
    c2 = public_key.encrypt(m2)
    c_sum = c1 + c2  # 加密结果直接相加
    m_sum = private_key.decrypt(c_sum)
    print(f"m1 = {m1}, m2 = {m2}")
    print(f"加密后的密文 c1: {c1.ciphertext()}")
    print(f"加密后的密文 c2: {c2.ciphertext()}")
    print(f"密文相加 c_sum: {c_sum.ciphertext()}")
    print(f"解密后的和: {m_sum}")
    print(f"明文加法结果: {m1 + m2}")
    print(f"加法同态性验证: {'通过' if m_sum == (m1 + m2) else '失败'}")

    # 验证标量乘同态性
    print("\n=== 验证标量乘同态性 ===")
    m = 10
    k = 3
    c = public_key.encrypt(m)
    c_prod = c * k  # 加密的消息直接与标量相乘
    m_prod = private_key.decrypt(c_prod)
    print(f"m = {m}, k = {k}")
    print(f"加密后的密文 c: {c.ciphertext()}")
    print(f"密文与标量相乘后的结果 c_prod: {c_prod.ciphertext()}")
    print(f"解密后的结果 m_prod: {m_prod}")
    print(f"明文乘法结果 k * m: {k * m}")
    print(f"标量乘同态性验证: {'通过' if m_prod == (k * m) else '失败'}")

