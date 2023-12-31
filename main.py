import os
import subprocess
import re
import json


# TODO: error handling for:
# - .vscode and settings.json not existing
# - npx not installed
# settings.json 'cSpell.words' not a key
def run_cspell_and_parse_output():
    # Run cspell command
    # npx --yes cspell "**/*" --exclude="**/target/**" --exclude="**/out/**"
    command = ["npx", "--yes", "cspell", "lint", "**/*", "--exclude=./main"]
    result = subprocess.run(
        command,
        capture_output=True,
        text=True,
        check=False
    )

    # Cspell writes normal output to stdout
    output = result.stdout

    # Parse output
    pattern = r'Unknown word \((.*?)\)'
    return re.findall(pattern, output)


unknown_words = run_cspell_and_parse_output()

if os.path.isfile('.vscode/settings.json'):
    with open('.vscode/settings.json') as f:
        data = json.load(f)

    data['cSpell.words'] = list(set(data['cSpell.words']) | set(unknown_words))

    with open('.vscode/settings.json', 'w') as f:
        json.dump(data, f, indent=2)
else:
    with open('.vscode/settings.json', 'w') as f:
        data = {
            "cSpell.words": list(set(unknown_words))
        }
        json.dump(data, f, indent=2)
