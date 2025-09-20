-- +goose Up
CREATE TABLE posts (
    id uuid UNIQUE NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    title text,
    description text,
    url text,
    published_at timestamp,
    feed_id uuid NOT NULL,
    CONSTRAINT pk_posts PRIMARY KEY (id),
    CONSTRAINT fk_posts_feed_id FOREIGN KEY (feed_id) REFERENCES feeds.id ON DELETE CASCADE,
    CONSTRAINT ck_posts_title_description_or_null CHECK (description IS NOT NULL OR title IS NOT NULL)
);

-- +goose Down
DROP TABLE posts;

