import scrapy
from collections import namedtuple
from ..items import SimpsonsItem as Item

Character = namedtuple('Character', ['first_name', 'last_name'])


class SimpsonSpider(scrapy.Spider):

    name = 'simpsons'
    start_urls = [
        'https://en.wikipedia.org/wiki/List_of_recurring_The_Simpsons_characters'
    ]

    def parse(self, response, **kwargs):
        data = response.css('div[class="mw-heading mw-heading3"] '
                            'h3::text').getall()
        for character in data:
            first_name, *last_name = character.split()
            # print(f'{first_name=} last_name={" ".join(last_name)}')

            item = Item()
            item['first_name'] = first_name
            item['last_name'] = ' '.join(last_name)
            yield item
