# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html
import boto3


dynamodb = boto3.resource('dynamodb')
table = dynamodb.Table('BluckBuster_movies')

class SimpsonsPipeline(object):

    def process_item(self, item, spider):
        table.put_item(item=dict(item))
        return item
