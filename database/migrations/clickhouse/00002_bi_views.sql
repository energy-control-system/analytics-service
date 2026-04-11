-- +goose Up
create view if not exists v_bi_tasks_daily as
select
    toDate(finished_at) as day,
    count() as tasks_count,
    countIf(inspection_type = 'limitation') as limitation_count,
    countIf(inspection_type = 'resumption') as resumption_count,
    countIf(inspection_type = 'verification') as verification_count,
    countIf(inspection_type = 'unauthorized_connection') as unauthorized_connection_count,
    countIf(inspection_is_violation_detected) as violations_detected_count,
    countIf(inspection_is_unauthorized_consumers) as unauthorized_consumers_count,
    round(avg(dateDiff('minute', started_at, finished_at)), 2) as avg_duration_minutes
from finished_tasks
group by day;

create view if not exists v_bi_brigade_performance as
select
    toDate(finished_at) as day,
    brigade_id,
    count() as tasks_count,
    round(avg(dateDiff('minute', started_at, finished_at)), 2) as avg_duration_minutes,
    countIf(inspection_type = 'limitation' and inspection_resolution = 'limited') as successful_limitations_count,
    countIf(inspection_type = 'resumption' and inspection_resolution = 'resumed') as successful_resumptions_count,
    countIf(inspection_is_violation_detected) as violations_detected_count
from finished_tasks
group by day, brigade_id;

create view if not exists v_bi_inspection_results as
select
    toDate(finished_at) as day,
    multiIf(
        inspection_type = 'limitation', 'Ограничение',
        inspection_type = 'resumption', 'Возобновление',
        inspection_type = 'verification', 'Контроль ограничения',
        inspection_type = 'unauthorized_connection', 'Несанкционированное подключение',
        'Неизвестно'
    ) as inspection_type_ru,
    multiIf(
        inspection_type = 'limitation' and inspection_resolution = 'limited', 'Ограничение введено',
        inspection_type = 'limitation', 'Недопуск',
        inspection_type = 'resumption' and inspection_resolution = 'resumed', 'Возобновление выполнено',
        inspection_type = 'resumption', 'Недопуск',
        inspection_is_violation_detected, 'Нарушение выявлено',
        'Нарушение не выявлено'
    ) as inspection_result_ru,
    count() as tasks_count
from finished_tasks
group by day, inspection_type_ru, inspection_result_ru;

create view if not exists v_bi_subscriber_object_profile as
select
    subscriber_id,
    subscriber_account_number,
    multiIf(
        subscriber_status = 'active', 'Активен',
        subscriber_status = 'violator', 'Нарушитель',
        subscriber_status = 'archived', 'Архивный',
        'Неизвестно'
    ) as subscriber_status_ru,
    object_id,
    object_address,
    object_have_automaton,
    if(object_have_automaton, 'Есть автомат', 'Нет автомата') as automaton_state_ru,
    max(toDate(finished_at)) as last_task_day,
    count() as total_tasks_count,
    countIf(inspection_is_violation_detected) as violations_detected_count,
    countIf(inspection_is_unauthorized_consumers) as unauthorized_consumers_count
from finished_tasks
group by
    subscriber_id,
    subscriber_account_number,
    subscriber_status_ru,
    object_id,
    object_address,
    object_have_automaton,
    automaton_state_ru;

-- +goose Down
drop view if exists v_bi_subscriber_object_profile;
drop view if exists v_bi_inspection_results;
drop view if exists v_bi_brigade_performance;
drop view if exists v_bi_tasks_daily;
