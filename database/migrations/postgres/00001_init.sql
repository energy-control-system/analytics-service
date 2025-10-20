-- +goose Up
create table if not exists report_types
(
    id   int primary key generated always as identity,
    name text not null
);

insert into report_types (name)
values ('Basic');

create table if not exists reports
(
    id           int primary key generated always as identity,
    type         int references report_types (id) on delete restrict,
    period_start date        not null,
    period_end   date        not null,
    created_at   timestamptz not null default now()
);

create table if not exists attachments
(
    id         int primary key generated always as identity,
    report_id  int         not null references reports (id) on delete cascade,
    file_id    int         not null,
    created_at timestamptz not null default now()
);

create index if not exists idx_reports_type on reports (type);
create index if not exists idx_reports_period on reports (period_start, period_end);
create index if not exists idx_attachments_report on attachments (report_id);

-- +goose Down
drop table if exists attachments;
drop table if exists reports;
drop table if exists report_types;
