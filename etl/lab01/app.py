# douban_wordcloud.py

import pandas as pd
import jieba
from collections import Counter
from wordcloud import WordCloud
import matplotlib.pyplot as plt
from pyecharts.charts import WordCloud as PyWordCloud
from pyecharts import options as opts

# 1. 读取图书标题
df = pd.read_csv("douban_books.csv")
titles = df["书名"].dropna().tolist()
text = " ".join(titles)

# 2. 加载停用词
with open("stopwords.txt", "r", encoding="utf-8") as f:
    stopwords = set(f.read().splitlines())

# 3. 分词并过滤
words = jieba.lcut(text)
filtered_words = [word for word in words if word not in stopwords and len(word) > 1]

# 4. 词频统计
word_freq = Counter(filtered_words)

# 5. matplotlib词云展示
wc = WordCloud(font_path="simhei.ttf", background_color="white", width=800, height=600)
wc.generate_from_frequencies(word_freq)

plt.figure(figsize=(10, 8))
plt.imshow(wc, interpolation="bilinear")
plt.axis("off")
plt.title("matplotlib", fontsize=16)
plt.show()

# 6. pyecharts词云展示
data = list(word_freq.items())

cloud = (
    PyWordCloud()
    .add("", data, word_size_range=[20, 100])
    .set_global_opts(title_opts=opts.TitleOpts(title="豆瓣图书标题词云（pyecharts）"))
)

cloud.render("douban_wordcloud.html")
print("pyecharts 词云已生成：douban_wordcloud.html")
