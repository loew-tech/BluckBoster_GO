from random import choice
from string import ascii_letters, digits, punctuation
from typing import Any

import scrapy

from ..items import MoviesItem as Item

class MoviesSpider(scrapy.Spider):
    name = 'movies'
    start_urls = [
        'https://editorial.rottentomatoes.com/guide/best-movies-of-all-time/',
        'https://editorial.rottentomatoes.com/guide/best-movies-of-all-time/2/'
    ]

    def __init__(self, name: str | None = None, **kwargs: Any):
        super().__init__(name, **kwargs)
        self._gen_id = MoviesSpider.get_gen_uuid()

    def parse(self, response, **kwargs):
        for data in response.css('div[class="row countdown-item"]'):
            item = Item()
            item['id'] = self._gen_id()
            header = data.css('div[class="row countdown-item-title-bar"]')
            item['title'] = header[0].css('div[class="article_movie_title"] div '
                                  'h2 a::text').get()
            item['year'] = header[0].css('div[class="article_movie_title"] div '
                                  'h2 span[class="subtle start-year"]::text').get()[1:-1]
            item['rating'] = header[0].css('div[class="article_movie_title"] div h2 '
                                   'span[class="tMeterScore"]::text').get()
            content = data.css('div[class="row countdown-item-details"]')
            item['review'] = content[0].css('div div[class="info '
                                    'critics-consensus"]::text').get()
            item['synopsis'] = content[0].css('div div[class="info synopsis"]::text').\
                get()
            item['cast'] = content[0].css('div div[class="info cast"] a::text').\
                getall()
            item['director'] = content[0].css('div div[class="info '
                                      'director"] a::text').get()
            yield item


    @staticmethod
    def get_gen_uuid():
        chars, used = ascii_letters + digits + punctuation, set()
        
        def gen_uuid():
            id = None
            while id in used:
                id = ''.join(choice(chars) for _ in range(7))
            used.add(id)
            return id
        return gen_uuid