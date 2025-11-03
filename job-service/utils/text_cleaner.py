import re
import string


def clean_text(text : str) -> str:
    # Normalizing the text to lower case letters
    text = text.lower()

    # Removing HTML Tags
    text = re.sub(r'<[^>]+>', '', text)

    # Removing Links 
    text = re.sub(r'http\S+|www\S+', '', text)

    # Removing punctuations
    text = text.translate(str.maketrans('', '', string.punctuation))

    # Removing non-alphanumeric characters
    text = re.sub(r'[^a-zA-Z0-9\s]', '', text)

    # Removing extra spaces
    text = re.sub(r'\s+', ' ', text)

    return text

