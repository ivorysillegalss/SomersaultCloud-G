use mini_gpt;

create table bot
(
    bot_id            int auto_increment
        primary key,
    adjustment_bot_id int                  null,
    is_delete         tinyint(1) default 0 not null,
    is_official       tinyint(1) default 0 not null
);

create table bot_config
(
    bot_id      int      null,
    init_prompt text     null,
    model       char(30) null,
    tag_id      int      null comment '一级标签的id 对应tag表中的子集标签'
);

create table bot_config_tag
(
    tag_id        int auto_increment comment '标签id'
        primary key,
    tag_prompt    text not null comment '该标签提示词',
    parent_tag_id int  not null comment '父级标签id'
);

create table bot_info
(
    bot_id      int          not null,
    description varchar(100) null,
    name        char(30)     not null,
    avatar      text         not null
);

create table chat
(
    chat_id          int auto_increment
        primary key,
    user_id          int        null,
    bot_id           int        not null,
    title            char(15)   null,
    last_update_time int        null,
    is_delete        tinyint(1) null
);

create table chat_ask
(
    record_id int          not null
        primary key,
    message   varchar(100) not null,
    time      int          null comment '创建时间',
    chat_id   int          not null,
    constraint chat_ask_pk
        unique (record_id)
);

create table chat_generation
(
    record_id int  not null
        primary key,
    chat_id   int  not null,
    time      int  null comment '创建时间',
    message   text not null,
    constraint chat_generation_pk
        unique (record_id)
);

create table chat_re
(
    chat_id          int auto_increment
        primary key,
    user_id          int          not null,
    bot_id           int          not null,
    title            varchar(255) null,
    last_update_time bigint       not null,
    Is_delete        tinyint(1)   not null,
    data             blob         null,
    constraint chat_re_pk
        unique (chat_id)
);

create table record_info
(
    record_id       int auto_increment
        primary key,
    chat_id         int       not null,
    create_time     int       null comment '创建的时间戳',
    reference_id    int       null comment '引用的对话的id',
    reference_token char(100) null comment '引用的字符串'
);

create table user_chat
(
    chat_id int auto_increment
        primary key,
    user_id int not null
);

create table user_info
(
    user_id  int auto_increment
        primary key,
    username int null,
    password int null
);

create table user_value
(
    user_value_id int auto_increment
        primary key,
    gender        tinyint(1)  null,
    age           int         null,
    sign          varchar(10) null comment '个性签名',
    pic           text        null comment '头像',
    email         varchar(30) null comment '邮箱
',
    user_id       int         null comment '对应的用户',
    role          int         null
)
    charset = utf8mb3;

