import glob
import json
from collections import defaultdict

import pandas as pd

dfs = defaultdict(pd.DataFrame)
for i, filepath in enumerate(glob.glob('/tmp/res/*.json')):
    with open(filepath, 'r') as f:
        j = json.load(f)
    for lang, counters in j["counters_by_language"].items():
        row = {}
        for key, value in counters["total"].items():
            row[f"total_{key}"] = [value]
        for key, value in counters["average"].items():
            row[f"average_{key}"] = [value]
        df = pd.DataFrame.from_dict(data=row)
        dfs[lang] = dfs[lang].append(df)

language = 'java'
metrics = [
    "total_lines_of_code",
    "average_lines_of_code",
    "average_keywords_complexity",
    "average_indentations_complexity",
    "average_indentations_diff_complexity",
]
for metric in metrics:
    pd.Series(dfs[language][metric].values).plot.hist()
