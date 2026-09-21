import json
import re
import os

os.makedirs('data', exist_ok=True)

try:
    with open('data/file_list.json', 'r', encoding='utf-8') as f:
        raw_files = json.load(f)
except Exception:
    try:
        with open('file_list.json', 'r', encoding='utf-8') as f:
            raw_files = json.load(f)
    except Exception:
        raw_files = []

master_table = {}
text_data = json.dumps(raw_files)
tables = re.findall(r'm_[a-z0-9_]+', text_data)

for t in set(tables):
    master_table[t] = []

base_essentials = ["m_arena_class", "m_unit", "m_shop_contents", "m_shop_group", "m_stage"]
for b in base_essentials:
    if b not in master_table:
        master_table[b] = []

with open('data/master_table.json', 'w', encoding='utf-8') as f:
    json.dump(master_table, f, indent=2)

print(f"¡Listo! data/master_table.json generado con {len(master_table)} tablas.")