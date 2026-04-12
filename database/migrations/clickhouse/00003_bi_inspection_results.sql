-- +goose Up
drop view if exists v_bi_inspection_results;
create view if not exists v_bi_inspection_results as
select
    day,
    inspection_type_ru,
    inspection_result_ru,
    subscriber_status_ru,
    tasks_count,
    round(tasks_count / sum(tasks_count) over (partition by day), 6) as day_tasks_share_ratio
from
(
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
        multiIf(
            subscriber_status = 'active', 'Активен',
            subscriber_status = 'violator', 'Нарушитель',
            subscriber_status = 'archived', 'Архивный',
            'Неизвестно'
        ) as subscriber_status_ru,
        count() as tasks_count
    from finished_tasks
    group by day, inspection_type_ru, inspection_result_ru, subscriber_status_ru
);

-- +goose Down
drop view if exists v_bi_inspection_results;
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
