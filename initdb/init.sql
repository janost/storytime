create table stories (
    id serial not null,
    "cast" varchar(250),
    location varchar(250),
    plot varchar(250),
    story_text text,
    model varchar(250),
    primary key(id)
);
