from pathlib import Path

def remove_links_prefix(docs_root: Path):
    print("Remove docs prefix from index.md")
    index = docs_root.joinpath('index.md')

    content = index.read_text()

    content = content.replace("](./docs/", "](")
    content = content.replace("](docs/", "](")

    index.write_text(content)

docs_root = Path('.').joinpath('docs')

remove_links_prefix(docs_root)