import re

def clean_resume_text(text: str) -> str:
    # Remove HTML tags
    text = re.sub(r'<[^>]+>', '', text)

    # Remove URLs
    text = re.sub(r'http\S+|www\S+', '', text)

    # Strip leading/trailing whitespace
    text = text.strip()

    return text
