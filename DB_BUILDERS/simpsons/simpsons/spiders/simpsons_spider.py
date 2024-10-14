import json
from random import choice

import scrapy
from collections import namedtuple
from ..items import SimpsonsItem as Item

Character = namedtuple('Character', ['first_name', 'last_name'])


class SimpsonSpider(scrapy.Spider):
    name = 'simpsons'
    start_urls = [
        'https://en.wikipedia.org/wiki/List_of_recurring_The_Simpsons_characters'
    ]
    _member_types = ['basic', 'advanced', 'premium']

    def parse(self, response, **kwargs):
        data = response.css('div[class="mw-heading mw-heading3"] '
                            'h3::text').getall()
        simpsons = []
        for character in data:
            first_name, *last_name = character.split()

            item = Item()
            item['first_name'] = first_name
            item['last_name'] = ' '.join(last_name)
            item['username'] = f'{first_name}_{"_".join(last_name)}'.lower() \
                if last_name else first_name
            item['member_type'] = choice(SimpsonSpider._member_types)
            simpsons.append(dict(item))
            yield item
        
        with open('../simpsons.json', 'w', encoding='utf-8') as simpsons_json:
            json.dump(simpsons, simpsons_json, ensure_ascii=False, indent=4, sort_keys=True)
