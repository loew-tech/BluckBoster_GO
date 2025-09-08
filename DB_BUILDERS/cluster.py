import json

import numpy as np
import pandas as pd
from sklearn.cluster import KMeans


def get_data(file_path = 'metrics.json') -> pd.DataFrame:
    with open(file_path, encoding='utf-8') as json_:
        data = json.load(json_)
    for d in data.values():
        if 'story telling' in d:
            d['story_telling'] = d.pop('story telling')
    return pd.DataFrame.from_records(metrics for metrics in data.values())


def get_movies_cluster(data: pd.DataFrame, n_clusters: int = 20) -> np.ndarray:
    kmeans = KMeans(n_clusters=n_clusters, random_state=0)
    kmeans.fit(data)
    return kmeans


if __name__ == '__main__':
    print('hello clustering')
    df = get_data()
    clusters = get_movies_cluster(df)
    print(len(clusters.cluster_centers_), clusters.cluster_centers_)
    print('All Done :)')
