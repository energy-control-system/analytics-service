-- +goose Up
create table if not exists finished_tasks
(
    -- Task fields
    task_id                              Int64,
    comment                              Nullable(String),
    plan_visit_at                        Nullable(DateTime('UTC')),
    started_at                           DateTime('UTC'),
    finished_at                          DateTime('UTC'),

    -- Inspection nested fields
    inspection_id                        Int64,
    inspection_type                      Enum8(
        'unknown' = 0,
        'limitation' = 1,
        'resumption' = 2,
        'verification' = 3,
        'unauthorized_connection' = 4
        ),
    inspection_resolution                Enum8(
        'unknown' = 0,
        'limited' = 1,
        'stopped' = 2,
        'resumed' = 3
        ),
    inspection_limit_reason              Nullable(String),
    inspection_method                    String,
    inspection_method_by                 Enum8(
        'unknown' = 0,
        'consumer' = 1,
        'inspector' = 2
        ),
    inspection_reason_type               Enum8(
        'unknown' = 0,
        'not_introduced' = 1,
        'consumer_limited' = 2,
        'inspector_limited' = 3,
        'resumed' = 4
        ),
    inspection_reason_description        Nullable(String),
    inspection_is_restriction_checked    Bool,
    inspection_is_violation_detected     Bool,
    inspection_is_expense_available      Bool,
    inspection_violation_description     Nullable(String),
    inspection_is_unauthorized_consumers Bool,
    inspection_unauthorized_description  Nullable(String),
    inspection_unauthorized_explanation  Nullable(String),
    inspection_inspect_at                DateTime('UTC'),
    inspection_energy_action_at          DateTime('UTC'),

    -- Brigade fields
    brigade_id                           Int64,
    brigade_inspectors                   Array(Tuple(
        id Int64,
        surname String,
        name String,
        patronymic String,
        phone_number String,
        email String,
        assigned_at DateTime('UTC')
        )),

    -- Object fields
    object_id                            Int64,
    object_address                       String,
    object_have_automaton                Bool,

    -- Subscriber nested fields
    subscriber_id                        Int64,
    subscriber_account_number            String,
    subscriber_surname                   String,
    subscriber_name                      String,
    subscriber_patronymic                String,
    subscriber_phone_number              String,
    subscriber_email                     String,
    subscriber_inn                       String,
    subscriber_birth_date                Date,
    subscriber_status                    Enum8(
        'unknown' = 0,
        'active' = 1,
        'violator' = 2,
        'archived' = 3
        )
)
    engine = MergeTree()
        order by finished_at
        partition by toYYYYMM(finished_at)
        ttl finished_at + interval 2 year delete
        settings index_granularity = 8192, merge_with_ttl_timeout = 86400;

create table if not exists reports
(
    id           Int64,
    type         Enum8(
        'unknown' = 0,
        'basic' = 1
        ),

    -- File nested fields
    file_id      Int64,
    file_name    String,
    file_size    Int64,
    file_bucket  Enum8(
        'images' = 1,
        'documents' = 2
        ),
    file_url     String,

    period_start DateTime('UTC'),
    period_end   DateTime('UTC'),
    created_at   DateTime('UTC')
)
    engine = MergeTree()
        order by (id, created_at)
        partition by toYYYYMM(created_at)
        ttl created_at + interval 2 year delete
        settings index_granularity = 8192, merge_with_ttl_timeout = 86400;

alter table reports
    add index idx_report_type type type set(0) granularity 4;
alter table reports
    add index idx_period_start period_start type minmax granularity 4;

-- +goose Down
drop table if exists reports;
drop table if exists finished_tasks;
