# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html


# useful for handling different item types with a single interface
import boto3

from itemadapter import ItemAdapter

dynamodb = boto3.resource('dynamodb')
table = dynamodb.Table('BluckBuster_movies')


class MoviesPipeline:

    def process_item(self, item, spider):
        table.put_item(Item=dict(item))
        return item
