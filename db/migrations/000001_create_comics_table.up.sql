CREATE TABLE IF NOT EXISTS comics(
    comic_id INTEGER PRIMARY KEY,
    title TEXT,
    date DATETIME,
    img_url TEXT,
    news TEXT,
    safe_title TEXT,
    transcription TEXT,
    alt_transcription TEXT,
    link TEXT
);