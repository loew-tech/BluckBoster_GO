import json

from time import sleep
from typing import Dict, Generator

from google import genai


with open('.env') as env_:
    API_KEY = env_.read().strip()
client = genai.Client(api_key=API_KEY) 
MODEL = "gemini-2.0-flash"


def get_movies() -> Generator:
    with open('movies.json', encoding='utf-8') as json_:
        data = json.load(json_)

    for d in data:
        yield d['title'], d['year'], d['id']


def get_movie_metrics(title: str, year: int | str) -> Dict[str, int]:
    response = client.models.generate_content(
        model=MODEL,
        contents=f"""
        grade the {year} movie {title} on the following criteria from 0 to 100 in json format: action, comedy, suspense, drama, horror, romance, fantasy, story telling, cinematography, writing, directing, and acting
        """
    )
    if response.text is None:
        return {}
    return json.loads(response.text[response.text.find('{'):
                                    response.text.rfind('}')+1])


def write_metrics_json():
    movies, metrics = get_movies(), {}
    for i, (title, year, id_) in enumerate(movies):
        print(i, id_, title, year)
        if (mets := get_movie_metrics(title, year)):
            metrics[id_] = mets
            with open(r'metrics.json', 'w', encoding='utf-8') as json_:
                json.dump(metrics, json_, ensure_ascii=False, indent=4)
        sleep(10)


if __name__ == '__main__':
    print('hello clustering')
    write_metrics_json()
    print('All Done :)')
