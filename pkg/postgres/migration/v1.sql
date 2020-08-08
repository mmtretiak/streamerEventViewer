CREATE TABLE users (
   id varchar(256),
   name varchar(100) NOT NULL CHECK (name <> ''),
   email varchar(100) NOT NULL UNIQUE,
   thumbnail_url varchar(256),
   PRIMARY KEY (id)
);

CREATE TABLE user_secrets (
  id varchar(256) PRIMARY KEY,
  user_id varchar(256) NOT NULL UNIQUE,
  scopes text[],
  auth_token varchar(256) NOT NULL,
  CONSTRAINT fk_users
      FOREIGN KEY(user_id)
          REFERENCES users(id)
);
