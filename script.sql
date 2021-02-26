create table if not exists customer
(
    name           varchar(255) not null,
    id             bigint generated always as identity
    constraint customer_pkey
    primary key,
    transaction_id bigint
    );

create table if not exists issuer
(
    customer_id bigint not null
    constraint issuer_pkey
    primary key
    constraint issuer_customer_fkey
    references customer
    on update cascade on delete cascade
);

create table if not exists investor
(
    customer_id bigint not null
    constraint investor_pkey
    primary key
    constraint investor_customer_fkey
    references customer
);

create table if not exists invoice
(
    id          bigint generated always as identity
    constraint invoice_pkey
    primary key,
    number      varchar(255)     not null,
    description varchar(255),
    face_value  double precision not null,
    issuer_id   bigint           not null
    constraint invoice_issuer_fkey
    references issuer
    );

create table if not exists sell_order
(
    id               bigint generated always as identity (maxvalue 2147483647)
    constraint sell_order_pkey
    primary key,
    invoice_id       bigint           not null
    constraint sell_order_invoice_fkey
    references invoice,
    seller_wants     double precision not null,
    sell_order_state varchar(255)
    );

create table if not exists bid
(
    id            bigint generated always as identity (maxvalue 2147483647)
    constraint bid_pkey
    primary key,
    size          double precision not null,
    amount        double precision not null,
    investor_id   bigint           not null
    constraint bid_investor_fkey
    references investor,
    sell_order_id bigint           not null
    constraint bid_sell_order_fkey
    references sell_order
    );

create table if not exists transaction
(
    id                   bigint generated always as identity (maxvalue 2147483647)
    constraint transaction_pkey
    primary key,
    date                 timestamp with time zone not null,
                                       transaction_type     varchar(255)             not null,
    details              varchar(255),
    transaction_d_c_type varchar(255)             not null,
    value                double precision,
    customer_id          bigint                   not null
    constraint transaction_customer_fkey
    references customer,
    sell_order           bigint
    constraint transaction_sell_order_fkey
    references sell_order
    );

