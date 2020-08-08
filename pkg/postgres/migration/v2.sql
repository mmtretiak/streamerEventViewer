CREATE TABLE streamers (
   id varchar(256),
   name varchar(100) NOT NULL CHECK (name <> '') UNIQUE,
   PRIMARY KEY (id)
);

CREATE TABLE users_to_streamers(
  id varchar(256),
  user_id varchar(256) NOT NULL,
  streamer_id varchar(256) NOT NULL,
  CONSTRAINT fk_streamers_users_user_id
      FOREIGN KEY(user_id)
          REFERENCES users(id),

  CONSTRAINT fk_streamers_users_streamer_id
      FOREIGN KEY(streamer_id)
          REFERENCES streamers(id),
  PRIMARY KEY (id)
);
