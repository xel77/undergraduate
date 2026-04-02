import os
import random
from cryptography.fernet import Fernet
from typing import Dict, List, Set, Tuple
import jieba
import jieba.analyse
import hashlib
import logging
import struct
from collections import defaultdict

jieba.setLogLevel(logging.INFO)

class SSE:
    def __init__(self):
        self.K1 = None  # 用于伪随机置换的密钥
        self.K2 = None  # 用于伪随机函数的密钥
        self.K3 = None  # 用于伪随机置换的密钥
        self.array_A = {}  # 存储加密节点的数组
        self.table_T = {}  # 查找表

    def setup(self):
        """生成密钥"""
        self.K1 = Fernet.generate_key()
        self.K2 = Fernet.generate_key()
        self.K3 = Fernet.generate_key()
        return self.K1, self.K2, self.K3

    def pseudo_random_perm(self, key: bytes, value: str) -> int:
        """伪随机置换函数"""
        combined = key + value.encode()
        return int.from_bytes(hashlib.sha256(combined).digest()[:4], 'big') % 10000

    def pseudo_random_func(self, key: bytes, value: str) -> bytes:
        """伪随机函数"""
        combined = key + value.encode()
        return hashlib.sha256(combined).digest()

    def build_index(self, documents: Dict[str, Set[str]]):
        """构建索引"""
        ctr = 1
        for keyword, doc_ids in documents.items():
            doc_ids = list(doc_ids)
            Ki_prev = Fernet.generate_key()

            gamma = self.pseudo_random_perm(self.K3, keyword)
            eta = self.pseudo_random_func(self.K2, keyword)

            first_addr = self.pseudo_random_perm(self.K1, str(ctr))
            self.table_T[gamma] = (first_addr, Ki_prev, eta)

            for j, doc_id in enumerate(doc_ids):
                Ki_curr = Fernet.generate_key() if j < len(doc_ids) - 1 else None
                next_addr = self.pseudo_random_perm(self.K1, str(ctr + 1)) if j < len(doc_ids) - 1 else None

                node = {
                    'doc_id': doc_id,
                    'next_key': Ki_curr,
                    'next_addr': next_addr
                }
                f = Fernet(Ki_prev)
                encrypted_node = f.encrypt(str(node).encode())
                self.array_A[self.pseudo_random_perm(self.K1, str(ctr))] = encrypted_node

                Ki_prev = Ki_curr
                ctr += 1

    def gen_trapdoor(self, keyword: str) -> Tuple[int, bytes]:
        """生成陷门"""
        gamma = self.pseudo_random_perm(self.K3, keyword)
        eta = self.pseudo_random_func(self.K2, keyword)
        return (gamma, eta)

    def search(self, trapdoor: Tuple[int, bytes]) -> List[str]:
        """搜索文档"""
        gamma, eta = trapdoor
        result = []

        if gamma not in self.table_T:
            return result

        addr, key, stored_eta = self.table_T[gamma]

        # 字节串比较
        if stored_eta != eta:
            return result

        while addr is not None:
            try:
                f = Fernet(key)
                encrypted_node = self.array_A[addr]
                decrypted_node = f.decrypt(encrypted_node).decode()
                node = eval(decrypted_node)
                result.append(node['doc_id'])
                addr = node['next_addr']
                key = node['next_key']
            except Exception as e:
                break

        return result

class DocumentManager:
    def __init__(self, documents_folder: str):
        self.documents_folder = documents_folder
        self.documents = defaultdict(set)

    def extract_keywords(self, content: str, topK=20) -> List[str]:
        """从文本中提取关键词"""
        keywords = set()
        try:
            # 使用 TF-IDF 和 TextRank 方法
            keywords.update(jieba.analyse.extract_tags(content, topK=topK))
            keywords.update(jieba.analyse.textrank(content, topK=topK))
        except Exception as e:
            print(f"关键词提取过程中发生错误：{str(e)}")
        return list(keywords)

    def build_documents(self) -> Dict[str, Set[str]]:
        """从文件中构建文档集合"""
        if not os.path.exists(self.documents_folder):
            print(f"文件夹 {self.documents_folder} 不存在")
            return {}

        for filename in os.listdir(self.documents_folder):
            if not filename.endswith('.txt'):
                continue

            try:
                with open(os.path.join(self.documents_folder, filename), 'r', encoding='utf-8') as f:
                    content = f.read().strip()
                    if not content:
                        print(f"文件 {filename} 为空")
                        continue

                    keywords = self.extract_keywords(content)
                    if not keywords:
                        print(f"从文件 {filename} 中没有提取到关键词")
                        continue

                    for keyword in keywords:
                        self.documents[keyword].add(filename)

            except Exception as e:
                print(f"处理文件 {filename} 时发生错误：{str(e)}")

        return self.documents

def main():
    documents_folder = r"test"  # 文件夹路径
    document_manager = DocumentManager(documents_folder)

    # 从文件构建文档集合
    documents = document_manager.build_documents()

    if not documents:
        print("没有成功构建文档索引")
        return

    # 初始化并运行SSE
    sse = SSE()
    sse.setup()

    # 构建索引
    sse.build_index(documents)

    while True:
        keyword = input("请输入要搜索的关键词 (输入 'q' 退出): ")
        if keyword.lower() == 'q':
            break

        # 生成陷门并搜索
        trapdoor = sse.gen_trapdoor(keyword)
        result = sse.search(trapdoor)

        if result:
            print(f"\n找到以下文件包含关键词 '{keyword}':")
            for i, file in enumerate(result, 1):
                print(f"{i}. {file}")
        else:
            print(f"\n未找到包含关键词 '{keyword}' 的文件")

if __name__ == "__main__":
    main()
