from collections import defaultdict
import json
from typing import Dict, List

import boto3
import numpy as np
import pandas as pd
from sklearn.cluster import KMeans


dynamodb = boto3.resource('dynamodb')
table = dynamodb.Table('BluckBoster_movies')


def get_data(file_path = 'metrics.json') -> pd.DataFrame:
    with open(file_path, encoding='utf-8') as json_:
        data = json.load(json_)
    for d in data.values():
        if 'story telling' in d:
            d['story_telling'] = d.pop('story telling')
    return data


def get_movies_clusters(data: List[Dict[str, int]], keys: List[str], n_clusters: int = 20) -> np.ndarray:
    dataframe_ = pd.DataFrame.from_records([mets[k] for k in keys] for mets in data)
    kmeans = KMeans(n_clusters=n_clusters, random_state=0)
    kmeans.fit(dataframe_)
    return kmeans


def add_centroids_to_dynamo(metrics: Dict[str, Dict[str, int]], keys: List[str], kmeans: KMeans) -> None:  
    def add_centroid(movie_id: str, centroid: int) -> None:
        key = {'id': movie_id}
        expr_attrs_vals = {':c': centroid}
        update_expr = 'SET centroid = :c'
        table.update_item(
            Key=key,
            ReturnValues='NONE',
            UpdateExpression=update_expr,
            ExpressionAttributeValues=expr_attrs_vals,
        )

    for movie, mets in metrics.items():
        centroid = kmeans.predict(np.array([mets[k] for k in keys]).reshape(1, -1))[0]
        print(f'prediction: {movie}: {centroid}')
        add_centroid(movie, int(centroid))


if __name__ == '__main__':
    print('hello clustering')
    metric_data = get_data()
    metric_keys = list(metric_data[list(metric_data.keys())[0]].keys())
    clusters = get_movies_clusters(metric_data.values(), metric_keys)
    add_centroids_to_dynamo(metric_data, metric_keys, clusters)
    print('All Done :)')
