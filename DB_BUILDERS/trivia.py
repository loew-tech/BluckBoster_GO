import json
from typing import Generator

from google import genai
from google.genai.types import GenerateContentResponse


with open('.env') as env_:
    API_KEY = env_.read().strip()
client = genai.Client(api_key=API_KEY) 
MODEL = "gemini-2.0-flash"


def get_movie_names() -> Generator:
    with open('movies.json', encoding='utf-8') as json_:
        data = json.load(json_)

    for d in data:
        yield (d['title'], d['year'])


def get_trivia(movie: str, year: int | str) -> str:
    def parse_response(resp: GenerateContentResponse) -> str:
        trivia = []
        text = resp.text.split('**Question ')
        for t in text[1:]:
            print(f'\t{t=}')
            _, q, *_, ans = t.strip().split('\n\n')
            ans = ans.split('\n')[0]
            trivia.append(f'{q}:{ans}')
        return '&:&'.join(trivia)


    response = client.models.generate_content(
        model=MODEL, 
        contents=f'give me 3 trivia questions and answers on the {year} movie {movie}'
    )
    return parse_response(response)


if __name__ == '__main__':
    print('Hello trivia')
    for entry in get_movie_names():
        print(entry)
        m, y = entry
        triv = get_trivia(m, y)
        print(f'\n{triv=}')
        break
    print('All Done :)')
