import json
import re
from time import sleep
from typing import Generator

import boto3
from google import genai
from google.genai.types import GenerateContentResponse


with open('.env') as env_:
    API_KEY = env_.read().strip()
client = genai.Client(api_key=API_KEY) 
MODEL = "gemini-2.0-flash"

dynamodb = boto3.resource('dynamodb')
table = dynamodb.Table('BluckBoster_movies')


def get_movies() -> Generator:
    with open('movies.json', encoding='utf-8') as json_:
        data = json.load(json_)

    for d in data:
        yield (d['title'], d['year'], d['id'])


def get_trivia(movie: str, year: int | str) -> str:
    def parse_response(resp: GenerateContentResponse) -> str:
        trivia = []
        pattern = r'\\*\\*Question \d+\\*\\*'
        trivia_qs = re.split(pattern, resp.text)[1:]
        for t in trivia_qs:
            q, ans = t.replace('*', '').strip().replace('\n', '').split('Answer:')
            trivia.append(f'{q}:{ans}')
        return '&:&'.join(trivia)


    response = client.models.generate_content(
        model=MODEL, 
        contents=f'give me 3 trivia questions and answers on the {year} movie {movie}'
    )
    return parse_response(response)


def update_movie(id_, trivia: str) -> None:
    if not trivia:
        return

    key = {'id': id_}
    expr_attrs_vals = {':t': trivia}
    update_expr = 'SET trivia = :t'
    
    table.update_item(
        Key=key,
        ReturnValues='NONE',
        UpdateExpression=update_expr,
        ExpressionAttributeValues=expr_attrs_vals,
    )
    


if __name__ == '__main__':
    print('Hello trivia')
    for entry in get_movies():
        m, y, movie_id = entry
        triv = get_trivia(m, y)
        update_movie(movie_id, triv)
        sleep(3)
    print('All Done :)')
