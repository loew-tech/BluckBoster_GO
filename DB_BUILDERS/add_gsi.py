import json
from string import ascii_letters

import boto3

with open('.env') as env_:
    API_KEY = env_.read().strip()

dynamodb = boto3.resource('dynamodb')
table = dynamodb.Table('BluckBoster_movies')


def get_movies():
    with open('movies.json', encoding='utf-8') as json_:
        data = json.load(json_)

    for d in data:
        yield d['id'], d['title'][0] if d['title'][0] in ascii_letters else '#'


def add_gsi(id_, chr_: str) -> None:
    if not id_:
        return

    key = {'id': id_}
    expr_attrs_vals = {':p': chr_}
    update_expr = 'SET paginate_key = :p'
    
    table.update_item(
        Key=key,
        ReturnValues='NONE',
        UpdateExpression=update_expr,
        ExpressionAttributeValues=expr_attrs_vals,
    )
    


if __name__ == '__main__':
    print('Hello GSI')
    for cnt, (i, c) in enumerate(get_movies()):
        print(f'{cnt}. id_={i} chr_={c}')
        add_gsi(i, c)
    print('ALL DONE')
