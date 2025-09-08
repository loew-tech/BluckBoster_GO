from collections import defaultdict
import json
from typing import Dict, List

import numpy as np
import pandas as pd
from sklearn.cluster import KMeans


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


def add_centroids_to_dynamo() -> None:
    pass


def get_centroid(kmeans: KMeans, metrics: Dict[str, Dict[str, int]], keys: List[str]) -> None:
    clusters = defaultdict(list)
    for movie, mets in metrics.items():
        centroid = kmeans.predict(np.array([mets[k] for k in keys]).reshape(1, -1))[0]
        print(f'prediction: {movie}: {centroid}')
        clusters[centroid].append(movie)


def add_centroid_to_dynamo() -> None:
    pass


if __name__ == '__main__':
    print('hello clustering')
    metric_data = get_data()
    metric_keys = list(metric_data[list(metric_data.keys())[0]].keys())
    # print(f'{metric_keys=}')
    clusters = get_movies_clusters(metric_data.values(), metric_keys)
    # print(len(clusters.cluster_centers_), clusters.cluster_centers_)
    get_centroid(clusters, metric_data, metric_keys)
    print('All Done :)')
