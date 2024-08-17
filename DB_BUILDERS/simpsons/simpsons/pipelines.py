# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html
import pyodbc
from random import randint
from datetime import datetime as dt
from datetime import timedelta
from random import choices
from string import ascii_lowercase as lc


class SimpsonsPipeline(object):

    def process_item(self, item, spider):
        return item
