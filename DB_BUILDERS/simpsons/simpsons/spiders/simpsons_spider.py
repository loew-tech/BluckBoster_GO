from random import choice
from string import ascii_letters, digits, punctuation
from typing import Any

import scrapy
from collections import namedtuple
from ..items import SimpsonsItem as Item

Character = namedtuple('Character', ['first_name', 'last_name'])


class SimpsonSpider(scrapy.Spider):
    name = 'simpsons'
    start_urls = [
        'https://en.wikipedia.org/wiki/List_of_recurring_The_Simpsons_characters'
    ]
    
    def __init__(self, name: str | None = None, **kwargs: Any):
        super().__init__(name, **kwargs)
        self._gen_id = SimpsonSpider.get_gen_uuid()

    def parse(self, response, **kwargs):
        data = response.css('div[class="mw-heading mw-heading3"] '
                            'h3::text').getall()
        for character in data:
            first_name, *last_name = character.split()

            item = Item()
            item['id'] = self._gen_id()
            item['first_name'] = first_name
            item['last_name'] = ' '.join(last_name)
            item['user_name'] = f'{first_name}_{"_".join(last_name)}'.lower() if last_name else first_name
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
