import scrapy
import os

class AnswerSpider(scrapy.Spider):
    name = "answer"
    custom_settings = {
    "ROBOTSTXT_OBEY": False
	}


    allowed_domains = ["mp.weixin.qq.com"]
    start_urls = ["https://mp.weixin.qq.com/s/68TqrkSWdt921UkiIFErUg"]

    def parse(self, response):
        imgs = response.xpath('//img[@class="rich_pages js_insertlocalimg"]')

        save_folder = "weixin_images"
        if not os.path.exists(save_folder):
            os.makedirs(save_folder)

        for img in imgs:
            src = img.xpath('@src').get()
            if src:
                absolute_url = response.urljoin(src)
                yield scrapy.Request(url=absolute_url, callback=self.save_image)

    def save_image(self, response):
        filename = response.url.split("/")[-1]
        filepath = os.path.join("weixin_images", filename)

        with open(filepath, "wb") as f:
            f.write(response.body)
        
        self.log(f"图片已保存: {filepath}")

