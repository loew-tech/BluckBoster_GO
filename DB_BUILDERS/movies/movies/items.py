# Define here the models for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

import scrapy


class MoviesItem(scrapy.Item):
    title = scrapy.Field()
    year = scrapy.Field()
    rating = scrapy.Field()
    review = scrapy.Field()
    synopsis = scrapy.Field()
    cast = scrapy.Field()
    director = scrapy.Field()
