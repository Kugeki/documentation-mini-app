create table article
(
    id          serial
        constraint article_pk
            primary key,
    name        text not null,
    description text not null
);

alter table article
    owner to university;

create table example
(
    id                 serial
        constraint example_pk
            primary key,
    name               text    not null,
    description        text    not null,
    code               text    not null,
    output             text    not null,
    highlight_language varchar
);

alter table example
    owner to university;

create table documentation
(
    id                         serial
        constraint documentation_pk
            primary key,
    name                       varchar not null,
    default_highlight_language varchar
);

alter table documentation
    owner to university;

create table documentation_articles
(
    documentation_id integer not null
        constraint documentation_articles_documentation_id_fk
            references documentation,
    article_id       integer not null
        constraint documentation_articles_article_id_fk
            references article,
    constraint documentation_articles_pk
        primary key (documentation_id, article_id)
);

alter table documentation_articles
    owner to university;

create table article_examples
(
    article_id integer           not null
        constraint article_examples_article_id_fk
            references article,
    example_id integer           not null
        constraint article_examples_example_id_fk
            references example,
    priority   integer default 0 not null,
    constraint article_examples_pk
        primary key (article_id, example_id)
);

alter table article_examples
    owner to university;
