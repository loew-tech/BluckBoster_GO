# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

import scrapy


class SimpsonsItem(scrapy.Item):
    id = scrapy.Field()
    first_name = scrapy.Field()
    last_name = scrapy.Field()
    username = scrapy.Field()
    member_type = scrapy.Field()
