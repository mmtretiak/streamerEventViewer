CREATE TABLE clips(
    id varchar(256),
    user_id varchar(256),
    streamer_id varchar(256),
    external_id varchar(256) NOT NULL UNIQUE ,
    edit_url varchar(256) NOT NULL,
    view_count integer DEFAULT 0,
    PRIMARY KEY (id),
    CONSTRAINT fk_clips_user_id
                  FOREIGN KEY (user_id)
                  REFERENCES users(id),
    CONSTRAINT fk_clips_streamer_id
        FOREIGN KEY (streamer_id)
            REFERENCES streamers(id)
);
