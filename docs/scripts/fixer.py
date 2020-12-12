import os
from shutil import copyfile
from pathlib import Path

def copy_files(root_path: Path, docs_path: Path):
    copyfile(root_path.joinpath('CHANGELOG.md'), docs_path.joinpath('CHANGELOG.md'))

    copyfile(root_path.joinpath('LICENSE'), docs_path.joinpath('LICENSE.md'))


root_path = Path(__file__).parent.absolute().joinpath('..', '..')
docs_path = root_path.joinpath('docs')

copy_files(root_path, docs_path)