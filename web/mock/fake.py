import json
import random

from datetime import datetime, timedelta
from faker import Faker

languages = {
    'english': ['en-gb', 'en-us'],
    'spanish': ['es-es', 'es-us'],
    'french': ['fr-fr', 'fr-ca'],
    'german': ['de-de', 'de-at'],
    'italian': ['it-it', 'it-ch'],
    'portugese': ['pt-pt', 'pt-br'],
    'russian': ['ru-ru'],
    'japanese': ['ja-jp'],
    'korean': ['ko-kr'],
    'chinese': ['zh-cn', 'zh-tw']
}

def generate_vocab():
    generate_vocab_stats()
    generate_vocab_timeline()

def generate_vocab_stats():
    fake = Faker()
    stats = {}
    stats['num_documents'] = fake.pyint(min_value=10000, max_value=20000)
    stats['num_words'] = fake.pyint(min_value=1000000, max_value=2000000)

    remaining_documents = stats['num_documents']
    for lang, locales in languages.items():
        lang_stats = {}
        if lang == list(languages.keys())[-1]:
            lang_stats['num_documents'] = remaining_documents
        else:
            lang_stats['num_documents'] = fake.pyint(min_value=int(remaining_documents*.1), max_value=int(remaining_documents*.5))
        lang_stats['num_words'] = fake.pyint(min_value=10000, max_value=20000)
        lang_stats['num_unique_words'] = fake.pyint(min_value=1000, max_value=2000)
        lang_stats['words_per_document'] = lang_stats['num_words'] / lang_stats['num_documents']

        locale_stats = {}
        if len(locales) == 1:
            locale_stats[locales[0]] = {}
            locale_stats[locales[0]]['num_documents'] = lang_stats['num_documents']
        elif len(locales) == 2:
            locale_stats[locales[0]] = {}
            locale_stats[locales[0]]['num_documents'] = fake.pyint(min_value=1, max_value=lang_stats['num_documents']-1)
            locale_stats[locales[1]] = {}
            locale_stats[locales[1]]['num_documents'] = lang_stats['num_documents'] - locale_stats[locales[0]]['num_documents']
        else:
            raise Exception('Unexpected number of locales')
        lang_stats['locales'] = locale_stats

        stats[lang] = lang_stats
        remaining_documents -= lang_stats['num_documents']

    with open('vocab/stats.json', 'w+') as f:
        json.dump(stats, f, indent=4)

def generate_vocab_timeline():
    fake = Faker()

    start = datetime(2020, 1, 1)
    days = [start + timedelta(days=x) for x in range(0, 365)]
    day_data = []
    for d in days:
        stats = {}
        stats['date'] = d.strftime('%Y-%m-%d')
        langs = {}
        for lang, locales in languages.items():
            if random.randint(0, 1) == 0:
                continue
            lang_stats = {}
            for loc in locales:
                if len(lang_stats) > 0 and random.randint(0, 1) == 0:
                    continue
                loc_stats = {}
                loc_stats['num_documents'] = fake.pyint(min_value=1, max_value=100)
                loc_stats['num_words'] = fake.pyint(min_value=300, max_value=2000)
                loc_stats['num_unique_words'] = fake.pyint(min_value=100, max_value=500)
                loc_stats['sentiment'] = fake.pyfloat(min_value=-1, max_value=1)
                num_most_common = fake.pyint(min_value=1, max_value=10)
                loc_stats['most_common'] = [fake.word() for x in range(0, num_most_common)]
                num_rare = fake.pyint(min_value=1, max_value=10)
                loc_stats['rare_words'] = [fake.word() for x in range(0, num_rare)]
                lang_stats[loc] = loc_stats
            langs[lang] = lang_stats
        stats['languages'] = langs
        day_data.append(stats)
    with open('vocab/timeline.json', 'w+') as f:
        json.dump(day_data, f, indent=4)

if __name__ == "__main__":
    vocab = generate_vocab()