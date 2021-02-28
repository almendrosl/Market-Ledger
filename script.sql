drop table if exists customer cascade;
create table customer
(
    name           varchar(255) not null,
    id             bigint generated always as identity
        constraint customer_pkey
            primary key
);

drop table if exists issuer cascade;
create table issuer
(
    customer_id bigint not null
        constraint issuer_pkey
            primary key
        constraint issuer_customer_fkey
            references customer
            on update cascade on delete cascade
);

drop table if exists investor cascade;
create table investor
(
    customer_id bigint not null
        constraint investor_pkey
            primary key
        constraint investor_customer_fkey
            references customer
);

drop table if exists invoice cascade;
create table invoice
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

drop table if exists sell_order cascade;
create table sell_order
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

drop table if exists bid cascade;
create table bid
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

drop table if exists transaction cascade;
create table transaction
(
    id               bigint generated always as identity (maxvalue 2147483647)
        constraint transaction_pkey
            primary key,
    date             timestamp with time zone not null,
    transaction_type varchar(255)             not null,
    details          varchar(255),
    debit            double precision,
    customer_id      bigint                   not null
        constraint transaction_customer_fkey
            references customer,
    sell_order_id    bigint
        constraint transaction_sell_order_fkey
            references sell_order,
    credit           double precision
);

INSERT INTO public.customer (name) VALUES ('investor_1');
INSERT INTO public.customer (name) VALUES ('issuer_2');
INSERT INTO public.investor (customer_id) VALUES (1);
INSERT INTO public.issuer (customer_id) VALUES (2);
INSERT INTO public.invoice (number, description, face_value, issuer_id) VALUES ('123123-332-11', '500 fabric black', 5000, 2);
INSERT INTO public.sell_order (invoice_id, seller_wants, sell_order_state) VALUES (1, 4500, 'Unlocked');
INSERT INTO public.transaction (date, transaction_type, details, debit, customer_id, sell_order_id, credit) VALUES ('2021-02-28 01:12:41.376000', 'Capital', 'Investor 1 add €1000 ', 0, 1, null, 1000);
INSERT INTO public.transaction (date, transaction_type, details, debit, customer_id, sell_order_id, credit) VALUES ('2021-02-28 01:12:41.376000', 'Cash', 'Investor 1 add €1000 ', 1000, 1, null, 0);