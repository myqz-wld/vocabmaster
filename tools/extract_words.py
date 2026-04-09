#!/usr/bin/env python3
"""
Extract high-quality words from ECDICT SQLite and JLPT JSON,
generate english.json and japanese.json for vocabmaster.
"""

import json
import sqlite3
import re
import sys

ECDICT_DB = "/tmp/ecdict-sqlite/stardict.db"
JLPT_JSON = "/tmp/jlpt.json"
OUTPUT_DIR = sys.argv[1] if len(sys.argv) > 1 else "."


def clean_phonetic(p):
    """Clean phonetic string - add slashes if missing"""
    if not p:
        return ""
    p = p.strip()
    if p and not p.startswith("/") and not p.startswith("["):
        p = "/" + p + "/"
    return p


def clean_translation(t):
    """Extract clean Chinese translation"""
    if not t:
        return ""
    # Remove [医] [法] etc. prefixes
    t = re.sub(r'\[.*?\]\s*', '', t)
    # Take first line only
    lines = [l.strip() for l in t.split('\n') if l.strip()]
    if not lines:
        return ""
    # Remove POS prefix like "n. " "v. "
    result = lines[0]
    result = re.sub(r'^[a-z]+\.\s*', '', result)
    # Clean up
    result = result.strip()
    # Truncate if too long
    if len(result) > 50:
        parts = result.split('；')
        if len(parts) > 1:
            result = '；'.join(parts[:3])
        else:
            parts = result.split(',')
            result = ','.join(parts[:3])
    return result


def extract_pos(pos_str):
    """Extract primary part of speech"""
    if not pos_str:
        return ""
    # pos format: "n:80/v:20" or "n:100"
    parts = pos_str.split('/')
    if parts:
        main = parts[0].split(':')[0].strip()
        pos_map = {
            'n': 'noun', 'v': 'verb', 'adj': 'adjective',
            'adv': 'adverb', 'prep': 'preposition',
            'conj': 'conjunction', 'pron': 'pronoun',
            'interj': 'interjection', 'det': 'determiner'
        }
        return pos_map.get(main, main)
    return ""


def determine_difficulty(row):
    """Determine difficulty level based on Collins, Oxford, BNC, frequency"""
    collins = row['collins'] or 0
    oxford = row['oxford'] or 0
    bnc = row['bnc'] or 99999
    frq = row['frq'] or 99999
    tag = row['tag'] or ''

    # Beginner: Oxford 3000 OR Collins 4-5 OR very high frequency
    if oxford == 1 or collins >= 4 or (bnc <= 2000 and frq <= 2000):
        return 1
    # Intermediate: Collins 3 OR medium frequency OR CET4/CET6
    if collins == 3 or (bnc <= 8000 and frq <= 8000):
        return 2
    if 'cet4' in tag or 'cet6' in tag:
        return 2
    # Advanced: must have at least one quality signal
    if collins in (1, 2):
        return 3
    if 'gre' in tag or 'ielts' in tag or 'toefl' in tag:
        return 3
    if 'cet6' in tag:
        return 3
    if bnc <= 20000 or frq <= 20000:
        return 3
    # No quality signal at all — skip (return -1)
    return -1


def is_base_form(row):
    """Check if word is likely a base form, not a conjugation/plural"""
    exchange = row['exchange'] or ''
    # If the word appears as someone else's derived form, it's not a base form
    # exchange field contains things like "s:apples" or "p:ran/d:running"
    # A word that has "0:base" means it IS a derived form of "base"
    if exchange.startswith('0:') or '/0:' in exchange:
        return False
    # Words ending in common inflection suffixes that have short base forms
    word = row['word']
    if re.match(r'.+(?:ing|ed|tion|ness|ment|able|ible|ful|less|ous|ive|ly|er|est|ism|ist)$', word):
        # Check if there's a shorter base form in exchange
        if exchange and '0:' not in exchange and len(word) > 6:
            return True  # Has exchange data but no base pointer — likely IS a base form
        if not exchange:
            return True  # No exchange data, accept it
    return True


def make_word_id(word, lang):
    """Create a word ID"""
    clean = re.sub(r'[^a-zA-Z0-9]', '_', word.lower())
    clean = re.sub(r'_+', '_', clean).strip('_')
    return f"{lang}_{clean}"


def extract_english(db_path, max_advanced=5000):
    """Extract English words from ECDICT"""
    conn = sqlite3.connect(db_path)
    conn.row_factory = sqlite3.Row

    # Filter out function words (articles, pronouns, prepositions, conjunctions, etc.)
    stop_words = {
        'the', 'a', 'an', 'be', 'is', 'am', 'are', 'was', 'were', 'been', 'being',
        'of', 'and', 'to', 'in', 'on', 'at', 'by', 'for', 'with', 'from',
        'it', 'he', 'she', 'we', 'they', 'me', 'him', 'her', 'us', 'them',
        'my', 'your', 'his', 'its', 'our', 'their', 'mine', 'yours', 'ours', 'theirs',
        'this', 'that', 'these', 'those', 'who', 'whom', 'which', 'what', 'whose',
        'i', 'you', 'do', 'does', 'did', 'done', 'doing',
        'have', 'has', 'had', 'having',
        'will', 'would', 'shall', 'should', 'may', 'might', 'can', 'could', 'must',
        'not', 'no', 'nor', 'but', 'or', 'if', 'then', 'so', 'as', 'than',
        'up', 'out', 'off', 'over', 'into', 'onto', 'upon', 'about', 'after',
        'before', 'between', 'through', 'during', 'without', 'within',
        'also', 'just', 'only', 'very', 'too', 'quite', 'rather', 'still',
        'here', 'there', 'where', 'when', 'how', 'why', 'all', 'each', 'every',
        'both', 'few', 'more', 'most', 'some', 'any', 'many', 'much', 'own',
        'other', 'another', 'such', 'even', 'well', 'back', 'now', 'way',
        'get', 'got', 'getting', 'go', 'going', 'gone', 'went',
        'come', 'came', 'coming', 'take', 'took', 'taken', 'taking',
        'make', 'made', 'making', 'give', 'gave', 'given', 'giving',
        'say', 'said', 'saying', 'see', 'saw', 'seen', 'seeing',
        'know', 'knew', 'known', 'knowing', 'think', 'thought', 'thinking',
        'let', 'put', 'set', 'keep', 'kept', 'keeping',
        # interjections / filler
        'oh', 'ooh', 'ah', 'uh', 'hmm', 'wow', 'hey', 'oops', 'yeah', 'yep',
        'nope', 'nah', 'huh', 'shh', 'psst', 'phew', 'gee', 'gosh', 'duh',
    }

    # Query ALL words with phonetic AND Chinese translation
    query = """
    SELECT word, phonetic, definition, translation, pos, collins, oxford, tag, bnc, frq, exchange
    FROM stardict
    WHERE phonetic != '' AND phonetic IS NOT NULL
      AND translation != '' AND translation IS NOT NULL
      AND length(word) > 2
      AND word NOT LIKE '%-%'
      AND word NOT LIKE '%.%'
      AND word NOT LIKE '%''%'
      AND word NOT LIKE '% %'
      AND word GLOB '[a-z]*'
    ORDER BY
      CASE WHEN oxford = 1 THEN 0 ELSE 1 END,
      CASE WHEN collins > 0 THEN 6 - collins ELSE 6 END,
      COALESCE(bnc, 99999)
    """

    rows = conn.execute(query).fetchall()
    conn.close()

    words_by_level = {1: [], 2: [], 3: []}
    seen = set()

    for row in rows:
        word = row['word']
        if word in seen or word in stop_words:
            continue
        seen.add(word)

        level = determine_difficulty(row)
        if level == -1:
            continue  # No quality signal, skip
        # No cap on beginner/intermediate, cap advanced at max_advanced
        if level == 3 and len(words_by_level[3]) >= max_advanced:
            continue
        # Filter out derived forms (plurals, conjugations)
        if not is_base_form(row):
            continue

        phonetic = clean_phonetic(row['phonetic'])
        chinese = clean_translation(row['translation'])
        pos = extract_pos(row['pos'])

        if not chinese or not phonetic:
            continue

        word_data = {
            "id": make_word_id(word, "en"),
            "language": "en",
            "text": word,
            "pronunciation": phonetic,
            "chinese_def": chinese,
            "difficulty": level,
            "part_of_speech": pos,
            "examples": [],
            "linked_word_ids": [],
            "tags": []
        }

        # Add tags based on ECDICT tag field
        tag = row['tag'] or ''
        if 'zk' in tag:
            word_data['tags'].append('中考')
        if 'gk' in tag:
            word_data['tags'].append('高考')
        if 'cet4' in tag:
            word_data['tags'].append('四级')
        if 'cet6' in tag:
            word_data['tags'].append('六级')
        if 'gre' in tag:
            word_data['tags'].append('GRE')
        if 'ielts' in tag:
            word_data['tags'].append('雅思')
        if 'toefl' in tag:
            word_data['tags'].append('托福')

        words_by_level[level].append(word_data)

    # Combine all levels
    all_words = words_by_level[1] + words_by_level[2] + words_by_level[3]

    print(f"English words extracted:")
    print(f"  Beginner: {len(words_by_level[1])}")
    print(f"  Intermediate: {len(words_by_level[2])}")
    print(f"  Advanced: {len(words_by_level[3])}")
    print(f"  Total: {len(all_words)}")

    return {
        "version": "1.0",
        "language": "en",
        "words": all_words
    }


def extract_japanese(jlpt_path, count_per_level=200):
    """Extract Japanese words from JLPT vocabulary"""
    with open(jlpt_path, 'r', encoding='utf-8') as f:
        data = json.load(f)

    # Map JLPT levels to difficulty
    # N5, N4 -> Beginner (1)
    # N3 -> Intermediate (2)
    # N2, N1 -> Advanced (3)
    level_map = {5: 1, 4: 1, 3: 2, 2: 3, 1: 3}

    words_by_level = {1: [], 2: [], 3: []}

    for word, entries in data.items():
        for entry in entries:
            jlpt_level = entry.get('level', 3)
            reading = entry.get('reading', '')
            difficulty = level_map.get(jlpt_level, 2)

            if len(words_by_level[difficulty]) >= count_per_level:
                continue

            # Create word ID from reading (romanized)
            word_id = make_word_id(reading if reading else word, "ja")
            # Avoid duplicate IDs
            base_id = word_id
            counter = 2
            existing_ids = {w['id'] for ws in words_by_level.values() for w in ws}
            while word_id in existing_ids:
                word_id = f"{base_id}_{counter}"
                counter += 1

            word_data = {
                "id": word_id,
                "language": "ja",
                "text": word,
                "pronunciation": reading,
                "chinese_def": "",  # Will be filled by LLM
                "difficulty": difficulty,
                "part_of_speech": "",
                "examples": [],
                "linked_word_ids": [],
                "tags": [f"JLPT-N{jlpt_level}"]
            }

            words_by_level[difficulty].append(word_data)

    all_words = words_by_level[1] + words_by_level[2] + words_by_level[3]

    print(f"\nJapanese words extracted:")
    print(f"  Beginner (N5+N4): {len(words_by_level[1])}")
    print(f"  Intermediate (N3): {len(words_by_level[2])}")
    print(f"  Advanced (N2+N1): {len(words_by_level[3])}")
    print(f"  Total: {len(all_words)}")

    return {
        "version": "1.0",
        "language": "ja",
        "words": all_words
    }


def main():
    # Extract English from ECDICT
    # English: no cap on beginner/intermediate, cap advanced at 5000
    english_data = extract_english(ECDICT_DB, max_advanced=5000)

    # Japanese: take ALL words
    japanese_data = extract_japanese(JLPT_JSON, count_per_level=99999)

    # Write output
    en_path = f"{OUTPUT_DIR}/english.json"
    ja_path = f"{OUTPUT_DIR}/japanese.json"

    with open(en_path, 'w', encoding='utf-8') as f:
        json.dump(english_data, f, ensure_ascii=False, indent=2)
    print(f"\nWritten to {en_path}")

    with open(ja_path, 'w', encoding='utf-8') as f:
        json.dump(japanese_data, f, ensure_ascii=False, indent=2)
    print(f"Written to {ja_path}")


if __name__ == '__main__':
    main()
