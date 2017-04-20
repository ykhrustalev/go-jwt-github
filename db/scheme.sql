CREATE TABLE github_user_details (
  id        SERIAL PRIMARY KEY,
  github_id INT  NOT NULL UNIQUE,
  info      JSON NOT NULL,
  emails    JSON NOT NULL
);


CREATE TABLE github_tokens (
  id             SERIAL PRIMARY KEY,
  github_user_id INTEGER REFERENCES github_user_details (id),
  value          JSON NOT NULL
);


CREATE TABLE jwt_tokens (
  id             SERIAL PRIMARY KEY,
  github_user_id INTEGER REFERENCES github_user_details (id),
  value          TEXT NOT NULL
);

