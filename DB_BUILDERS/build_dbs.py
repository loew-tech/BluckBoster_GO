import json
import os

import boto3


if __name__ == '__main__':
    os.system('cd ./movies & scrapy crawl movies')
    os.system('cd ./simpsons & scrapy crawl simpsons')

    dynamodb = boto3.resource('dynamodb')

    with open('movies.json', encoding='utf-8') as movies_json:
        movies, table = json.load(movies_json), dynamodb.Table('BluckBoster_movies')
        with table.batch_writer() as batch:
            for movie in movies:
                batch.put_item(Item=movie)

    with open('simpsons.json', encoding='utf-8') as character_json:
        characters, table = json.load(character_json), dynamodb.Table('BluckBoster_members')
        with table.batch_writer() as batch:
            for character in characters:
                batch.put_item(Item=character)

    print('\nALL DONE :)')
