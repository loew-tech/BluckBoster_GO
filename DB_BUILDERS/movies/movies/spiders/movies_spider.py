import scrapy

from ..items import MoviesItem as Item


class MoviesSpider(scrapy.Spider):
    name = 'movies'
    start_urls = [
        'https://editorial.rottentomatoes.com/guide/best-movies-of-all-time/',
        'https://editorial.rottentomatoes.com/guide/best-movies-of-all-time/2/'
    ]

    def parse(self, response, **kwargs):
        for data in response.css('div[class="row countdown-item"]'):
            item = Item()
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
